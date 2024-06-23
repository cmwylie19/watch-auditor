/*
Copyright © 2024 Case Wylie <casewylie@gmail.com>

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

	"github.com/cmwylie19/watch-auditor/src/config/lang"
	"github.com/cmwylie19/watch-auditor/src/pkg/logging"
	"github.com/cmwylie19/watch-auditor/src/pkg/server"
	"github.com/spf13/cobra"
)

var (
	port             int
	every            int
	unit             string
	metrics          bool
	logLevel         string
	mode             string
	failureThreshold int
)

func init() {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: lang.CmdServeShort,
		Long:  lang.CmdServeLong,
		Run: func(cmd *cobra.Command, args []string) {

			logging.SetupLogging(logLevel)
			logging.Info(fmt.Sprintf("Server is starting on %d", port))
			if mode != "audit" && mode != "enforcing" {
				logging.Debug("Mode must be either 'audit' or 'enforcing', defaulting to 'enforcing'")
				mode = "enforcing"
			}

			server := server.Server{
				Port:  port,
				Every: every,
				Unit:  unit,
				Mode:  mode,
			}

			if err := server.Start(); err != nil {
				logging.Error(err.Error())
				os.Exit(1)
			}
		},
	}
	serveCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Port to listen on (default: 8080)")
	serveCmd.PersistentFlags().IntVarP(&every, "every", "e", 1, "Interval to check")
	serveCmd.PersistentFlags().StringVarP(&unit, "unit", "u", "minute", "Unit of time to check (minute, hour, day)")
	serveCmd.PersistentFlags().BoolVar(&metrics, "metrics", true, "Enable metrics")
	serveCmd.Flags().IntVarP(&failureThreshold, "failure-threshold", "f", 3, "Failure threshold to roll watch controller pod")
	serveCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, error)")
	serveCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "enforcing", "Mode to run in (audit, enforcing)")
	rootCmd.AddCommand(serveCmd)
}
