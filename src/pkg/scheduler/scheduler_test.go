package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/cmwylie19/watch-auditor/src/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestScheduler_CreatePod(t *testing.T) {
	logger := &logging.MockLogger{}
	clientset := fake.NewSimpleClientset()
	scheduler := &Scheduler{
		Every:     10 * time.Second,
		client:    clientset,
		Namespace: "default",
		Logger:    logger,
	}

	name := "testpod"
	scheduler.CreatePod(name)

	createdPod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), "auto-testpod", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Expected pod to be created, but got error: %v", err)
	}

	if createdPod.Name != "auto-testpod" {
		t.Errorf("Expected pod name to be auto-testpod, got %s", createdPod.Name)
	}
}

func TestScheduler_CheckPod(t *testing.T) {
	logger := &logging.MockLogger{}
	clientset := fake.NewSimpleClientset()
	watchFailures := promauto.NewCounter(prometheus.CounterOpts{
		Name: "test_watch_failures",
		Help: "Number of watch failures for testing",
	})
	scheduler := &Scheduler{
		Every:               10 * time.Second,
		client:              clientset,
		Namespace:           "default",
		Logger:              logger,
		watchFailuresMetric: watchFailures,
	}

	name := "testpod"
	scheduler.CreatePod(name)

	// Simulate deleting the pod to test the check
	clientset.CoreV1().Pods("default").Delete(context.TODO(), "auto-"+name, metav1.DeleteOptions{})

	// Now call CheckPod which should detect the pod was deleted
	scheduler.CheckPod(name)

	expectedLog := "INFO: Watch Controller successfully deleted pod: auto-testpod"
	if logger.Messages[len(logger.Messages)-1] != expectedLog {
		t.Errorf("Expected log message '%s', got '%s'", expectedLog, logger.Messages[len(logger.Messages)-1])
	}
}
