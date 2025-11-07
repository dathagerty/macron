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
	Run: func(cmd *cobra.Command, args []string) {
		// Get flag values
		name, _ := cmd.Flags().GetString("name")
		interval, _ := cmd.Flags().GetString("interval")
		script, _ := cmd.Flags().GetString("script")

		fmt.Printf("Creating launchd cron task:\n")
		fmt.Printf("  Name: %s\n", name)
		fmt.Printf("  Interval: %s\n", interval)
		fmt.Printf("  Script: %s\n", script)
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
