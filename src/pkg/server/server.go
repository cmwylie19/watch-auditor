package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cmwylie19/watch-auditor/src/pkg/handlers"
	"github.com/cmwylie19/watch-auditor/src/pkg/logging"
	"github.com/cmwylie19/watch-auditor/src/pkg/scheduler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
)

type Server struct {
	Port      int           `json:"port"`
	Every     time.Duration `json:"every"`
	Handlers  *handlers.Handlers
	Namespace string `json:"namespace"`
	Logger    logging.LoggerInterface
	Client    kubernetes.Interface
}

func (s *Server) Start() error {
	http.HandleFunc("/healthz", s.Handlers.Healthz)
	http.Handle("/metrics", promhttp.Handler())
	scheduler := scheduler.NewScheduler(s.Every, s.Namespace, s.Logger, s.Client)
	go scheduler.Start()
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}
