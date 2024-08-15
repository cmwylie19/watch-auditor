package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/cmwylie19/watch-auditor/src/config/lang"
	logging "github.com/cmwylie19/watch-auditor/src/pkg/logging"
	"github.com/cmwylie19/watch-auditor/src/pkg/server"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	port      int
	every     time.Duration
	logLevel  string
	namespace string
)

func init() {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: lang.CmdServeShort,
		Long:  lang.CmdServeLong,
		Run: func(cmd *cobra.Command, args []string) {

			logger, err := logging.NewLogger("")
			if err != nil {
				fmt.Printf("Failed to initialize logger: %v\n", err)
				os.Exit(1)
			}
			defer logger.CloseFile()

			switch logLevel {
			case "debug":
				logger.SetLevel(slog.LevelDebug)
			case "info":
				logger.SetLevel(slog.LevelInfo)
			case "warn":
				logger.SetLevel(slog.LevelWarn)
			case "error":
				logger.SetLevel(slog.LevelError)
			default:
				logger.SetLevel(slog.LevelInfo) // Default to INFO level
			}

			config, err := rest.InClusterConfig()
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to create in-cluster config: %v", err))
				os.Exit(1)
			}

			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to create Kubernetes client: %v", err))
				os.Exit(1)
			}

			server := server.Server{
				Port:      port,
				Every:     every,
				Namespace: namespace,
				Logger:    logger,
				Client:    clientset,
			}

			if err := server.Start(); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}

			logger.Info("Server started successfully")
		},
	}

	serveCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	serveCmd.PersistentFlags().DurationVarP(&every, "every", "e", 30*time.Second, "Interval to check in seconds")
	serveCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, error)")
	serveCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "pepr-demo", "Namespace to audit")
	rootCmd.AddCommand(serveCmd)
}
