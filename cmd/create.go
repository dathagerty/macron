/*
Copyright © 2025 David Hagerty <david@dathagerty.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const (
	labelPrefix      = "com.macron"
	launchAgentsDirPerms = 0755
	plistFilePerms       = 0644
)

var (
	createName     string
	createInterval string
	createScript   string
)

// validateAndResolveScript validates the script exists and returns its absolute path
func validateAndResolveScript(scriptPath string) (string, error) {
	absScript, err := filepath.Abs(scriptPath)
	if err != nil {
		return "", fmt.Errorf("error resolving script path: %w", err)
	}

	if _, err := os.Stat(absScript); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("script file does not exist: %s", absScript)
		}
		return "", fmt.Errorf("error checking script file: %w", err)
	}

	return absScript, nil
}

// generatePlist creates the plist XML content for a launchd task
func generatePlist(label, scriptPath string, intervalSeconds int) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>%s</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>StartInterval</key>
	<integer>%d</integer>
	<key>RunAtLoad</key>
	<true/>
	<key>StandardOutPath</key>
	<string>/tmp/%s.stdout</string>
	<key>StandardErrorPath</key>
	<string>/tmp/%s.stderr</string>
</dict>
</plist>`, label, scriptPath, intervalSeconds, label, label)
}

// writePlistFile writes the plist content to the LaunchAgents directory
func writePlistFile(name, absScript string, intervalSeconds int) (string, string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", "", fmt.Errorf("error getting home directory: %w", err)
	}

	launchAgentsDir := filepath.Join(home, "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, launchAgentsDirPerms); err != nil {
		return "", "", fmt.Errorf("error creating LaunchAgents directory: %w", err)
	}

	label := fmt.Sprintf("%s.%s", labelPrefix, name)
	plistPath := filepath.Join(launchAgentsDir, label+".plist")

	plistContent := generatePlist(label, absScript, intervalSeconds)
	if err := os.WriteFile(plistPath, []byte(plistContent), plistFilePerms); err != nil {
		return "", "", fmt.Errorf("error writing plist file: %w", err)
	}

	return label, plistPath, nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new launchd cron task",
	Long: `Create a new launchd cron task with NAME to run SCRIPT over an INTERVAL.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate interval can be parsed as time.Duration
		duration, err := time.ParseDuration(createInterval)
		if err != nil {
			return fmt.Errorf("invalid interval '%s': must be a valid duration (e.g., 1h, 30m, 1h30m): %w", createInterval, err)
		}

		// Validate script file exists and resolve to absolute path
		absScript, err := validateAndResolveScript(createScript)
		if err != nil {
			return err
		}

		// Convert duration to seconds
		intervalSeconds := int(duration.Seconds())

		// Write plist file to LaunchAgents directory
		label, plistPath, err := writePlistFile(createName, absScript, intervalSeconds)
		if err != nil {
			return err
		}

		// Output success message
		fmt.Printf("Successfully created launchd task:\n")
		fmt.Printf("  Label: %s\n", label)
		fmt.Printf("  Interval: %s (%d seconds)\n", createInterval, intervalSeconds)
		fmt.Printf("  Script: %s\n", absScript)
		fmt.Printf("  Plist: %s\n", plistPath)
		fmt.Printf("\nTo load the task, run:\n")
		fmt.Printf("  launchctl load %s\n", plistPath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Define flags for the create command
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Name of the launchd task (required)")
	createCmd.Flags().StringVarP(&createInterval, "interval", "i", "", "Interval for the task execution (required)")
	createCmd.Flags().StringVarP(&createScript, "script", "s", "", "Path to the script to execute (required)")

	// Mark flags as required
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("interval")
	createCmd.MarkFlagRequired("script")
}
