package server

import (
	"fmt"
	"net/http"

	"github.com/cmwylie19/watch-auditor/src/pkg/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Port             int    `json:"port"`
	Every            int    `json:"every"`
	Unit             string `json:"unit"`
	Handlers         *handlers.Handlers
	Mode             string `json:"mode"`
	FailureThreshold int    `json:"failureThreshold"`
}

func (s *Server) Start() error {
	http.HandleFunc("/healthz", s.Handlers.Healthz)
	http.Handle("/metrics", promhttp.Handler())
	// scheduler := scheduler.NewScheduler(s.failureThreshold,s.Every, s.Unit)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}