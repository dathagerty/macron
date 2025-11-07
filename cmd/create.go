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

var (
	name     string
	interval string
	script   string
)

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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new launchd cron task",
	Long: `Create a new launchd cron task with NAME to run SCRIPT over an INTERVAL.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		interval, _ := cmd.Flags().GetString("interval")
		script, _ := cmd.Flags().GetString("script")

		// Validate interval can be parsed as time.Duration
		duration, err := time.ParseDuration(interval)
		if err != nil {
			return fmt.Errorf("invalid interval '%s': must be a valid duration (e.g., 1h, 30m, 1h30m): %w", interval, err)
		}

		// Convert script to absolute path
		absScript, err := filepath.Abs(script)
		if err != nil {
			return fmt.Errorf("error resolving script path: %w", err)
		}

		// Validate script file exists
		if _, err := os.Stat(absScript); os.IsNotExist(err) {
			return fmt.Errorf("script file does not exist: %s", absScript)
		} else if err != nil {
			return fmt.Errorf("error checking script file: %w", err)
		}

		// Get home directory
		home, err := homedir.Dir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}

		// Create LaunchAgents directory path
		launchAgentsDir := filepath.Join(home, "Library", "LaunchAgents")

		// Ensure LaunchAgents directory exists
		if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
			return fmt.Errorf("error creating LaunchAgents directory: %w", err)
		}

		// Generate label and filename
		label := fmt.Sprintf("com.macron.%s", name)
		plistFilename := fmt.Sprintf("%s.plist", label)
		plistPath := filepath.Join(launchAgentsDir, plistFilename)

		// Convert duration to seconds
		intervalSeconds := int(duration.Seconds())

		// Generate plist content
		plistContent := generatePlist(label, absScript, intervalSeconds)

		// Write plist file
		if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
			return fmt.Errorf("error writing plist file: %w", err)
		}

		fmt.Printf("Successfully created launchd task:\n")
		fmt.Printf("  Label: %s\n", label)
		fmt.Printf("  Interval: %s (%d seconds)\n", interval, intervalSeconds)
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
	createCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the launchd task (required)")
	createCmd.Flags().StringVarP(&interval, "interval", "i", "", "Interval for the task execution (required)")
	createCmd.Flags().StringVarP(&script, "script", "s", "", "Path to the script to execute (required)")

	// Mark flags as required
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("interval")
	createCmd.MarkFlagRequired("script")
}
