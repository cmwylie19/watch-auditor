package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/cmwylie19/watch-auditor/src/pkg/handlers"
	"github.com/cmwylie19/watch-auditor/src/pkg/logging"
	"k8s.io/client-go/kubernetes/fake"
)

func TestServer_Start(t *testing.T) {
	logger := logging.NewMockLogger()
	mockHandlers := &handlers.Handlers{}
	clientset := fake.NewSimpleClientset()

	server := &Server{
		Port:      8080,
		Every:     10 * time.Second,
		Handlers:  mockHandlers,
		Namespace: "default",
		Logger:    logger,
		Client:    clientset,
	}

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			t.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	req, err := http.NewRequest("GET", "http://localhost:8080/healthz", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
	}

	req, err = http.NewRequest("GET", "http://localhost:8080/metrics", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
	}
}
