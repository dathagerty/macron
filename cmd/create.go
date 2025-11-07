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
	"time"

	"github.com/spf13/cobra"
)

var (
	name     string
	interval string
	script   string
)

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

		// Validate script file exists
		if _, err := os.Stat(script); os.IsNotExist(err) {
			return fmt.Errorf("script file does not exist: %s", script)
		} else if err != nil {
			return fmt.Errorf("error checking script file: %w", err)
		}

		fmt.Printf("Creating launchd cron task:\n")
		fmt.Printf("  Name: %s\n", name)
		fmt.Printf("  Interval: %s (%v)\n", interval, duration)
		fmt.Printf("  Script: %s\n", script)

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
