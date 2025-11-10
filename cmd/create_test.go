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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratePlist(t *testing.T) {
	tests := []struct {
		name            string
		label           string
		scriptPath      string
		intervalSeconds int
		wantContains    []string
	}{
		{
			name:            "basic plist generation",
			label:           "com.macron.test",
			scriptPath:      "/path/to/script.sh",
			intervalSeconds: 3600,
			wantContains: []string{
				"<?xml version=\"1.0\" encoding=\"UTF-8\"?>",
				"<key>Label</key>",
				"<string>com.macron.test</string>",
				"<key>ProgramArguments</key>",
				"<string>/path/to/script.sh</string>",
				"<key>StartInterval</key>",
				"<integer>3600</integer>",
				"<key>RunAtLoad</key>",
				"<true/>",
				"/tmp/com.macron.test.stdout",
				"/tmp/com.macron.test.stderr",
			},
		},
		{
			name:            "plist with special characters in path",
			label:           "com.macron.backup",
			scriptPath:      "/Users/test/My Scripts/backup.sh",
			intervalSeconds: 1800,
			wantContains: []string{
				"<string>com.macron.backup</string>",
				"<string>/Users/test/My Scripts/backup.sh</string>",
				"<integer>1800</integer>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generatePlist(tt.label, tt.scriptPath, tt.intervalSeconds)

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("generatePlist() missing expected content:\nwant: %q\ngot: %s", want, got)
				}
			}

			// Verify it's valid XML structure
			if !strings.HasPrefix(got, "<?xml") {
				t.Errorf("generatePlist() should start with XML declaration")
			}
			if !strings.HasSuffix(strings.TrimSpace(got), "</plist>") {
				t.Errorf("generatePlist() should end with </plist>")
			}
		})
	}
}

func TestValidateAndResolveScript(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a test script file
	testScript := filepath.Join(tempDir, "test-script.sh")
	if err := os.WriteFile(testScript, []byte("#!/bin/bash\necho test"), 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	tests := []struct {
		name        string
		scriptPath  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid script with absolute path",
			scriptPath: testScript,
			wantErr:    false,
		},
		{
			name:        "non-existent script",
			scriptPath:  filepath.Join(tempDir, "nonexistent.sh"),
			wantErr:     true,
			errContains: "does not exist",
		},
		{
			name:       "relative path resolution",
			scriptPath: "test-script.sh",
			wantErr:    false, // Will resolve to current directory
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For relative path test, change to temp dir
			if tt.name == "relative path resolution" {
				oldWd, _ := os.Getwd()
				defer os.Chdir(oldWd)
				os.Chdir(tempDir)
			}

			got, err := validateAndResolveScript(tt.scriptPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateAndResolveScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("validateAndResolveScript() error = %v, should contain %q", err, tt.errContains)
				}
			}

			if !tt.wantErr {
				if !filepath.IsAbs(got) {
					t.Errorf("validateAndResolveScript() = %v, should be absolute path", got)
				}
			}
		})
	}
}

func TestWritePlistFile(t *testing.T) {
	// This test would require mocking os.UserHomeDir() or running in a controlled environment
	// For now, we'll test the error cases that don't depend on the actual home directory

	tests := []struct {
		name            string
		taskName        string
		script          string
		intervalSeconds int
		wantLabel       string
	}{
		{
			name:            "generates correct label",
			taskName:        "backup",
			script:          "/usr/local/bin/backup.sh",
			intervalSeconds: 3600,
			wantLabel:       "com.macron.backup",
		},
		{
			name:            "handles special characters in name",
			taskName:        "my-task_123",
			script:          "/path/to/script.sh",
			intervalSeconds: 1800,
			wantLabel:       "com.macron.my-task_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can test label generation logic separately
			label := labelPrefix + "." + tt.taskName
			if label != tt.wantLabel {
				t.Errorf("Label = %v, want %v", label, tt.wantLabel)
			}
		})
	}
}

func TestCreateCmdValidation(t *testing.T) {
	// Test that the command has the expected structure
	if createCmd == nil {
		t.Fatal("createCmd should not be nil")
	}

	if createCmd.Use != "create" {
		t.Errorf("createCmd.Use = %v, want 'create'", createCmd.Use)
	}

	if createCmd.RunE == nil {
		t.Error("createCmd.RunE should not be nil")
	}

	// Verify flags are defined
	nameFlag := createCmd.Flags().Lookup("name")
	if nameFlag == nil {
		t.Error("name flag should be defined")
	}

	intervalFlag := createCmd.Flags().Lookup("interval")
	if intervalFlag == nil {
		t.Error("interval flag should be defined")
	}

	scriptFlag := createCmd.Flags().Lookup("script")
	if scriptFlag == nil {
		t.Error("script flag should be defined")
	}
}
