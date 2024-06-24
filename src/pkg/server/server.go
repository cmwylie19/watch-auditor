package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cmwylie19/watch-auditor/src/pkg/handlers"
	"github.com/cmwylie19/watch-auditor/src/pkg/scheduler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Port             int           `json:"port"`
	Every            time.Duration `json:"every"`
	Handlers         *handlers.Handlers
	Mode             string `json:"mode"`
	FailureThreshold int    `json:"failureThreshold"`
}

func (s *Server) Start() error {
	http.HandleFunc("/healthz", s.Handlers.Healthz)
	http.Handle("/metrics", promhttp.Handler())
	scheduler := scheduler.NewScheduler(s.FailureThreshold, s.Every, s.Mode)
	go scheduler.Start()
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}
