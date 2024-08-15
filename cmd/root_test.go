package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
)

// Mock
var osExit = os.Exit

func TestExecute(t *testing.T) {
	// Override os.Exit to intercept calls during tests.
	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}

	defer func() { osExit = os.Exit }() // Restore original

	// Simulate successful execution
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		exitCode = 0
	}
	Execute()
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, but got %d", exitCode)
	}

	// Simulate execution with an error
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		exitCode = 1
	}
	Execute()
	if exitCode != 1 {
		t.Fatalf("Expected exit code 1, but got %d", exitCode)
	}
}
