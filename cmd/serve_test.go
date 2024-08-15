package cmd

import (
	"testing"

	"time"

	"github.com/cmwylie19/watch-auditor/src/pkg/logging"
	"github.com/spf13/cobra"
)

func TestServeCommand(t *testing.T) {
	logger := logging.NewMockLogger()

	serveCmd := &cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info("Server command executed")
		},
	}

	// Manually add the flags as done in the init() function
	serveCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	serveCmd.PersistentFlags().DurationVarP(&every, "every", "e", 30*time.Second, "Interval to check in seconds")
	serveCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, error)")
	serveCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "pepr-demo", "Namespace to audit")

	serveCmd.SetArgs([]string{"serve", "--port=8080", "--log-level=info"})
	err := serveCmd.Execute()
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	// Validate the logger captured the expected message
	if len(logger.Messages) == 0 || logger.Messages[0] != "INFO: Server command executed" {
		t.Errorf("Expected 'Server command executed' log message, got %v", logger.Messages)
	}
}
