package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/cmwylie19/watch-auditor/src/config/lang"
	"github.com/cmwylie19/watch-auditor/src/pkg/logging"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Scheduler struct {
	Every                time.Duration
	client               *kubernetes.Clientset
	watchFailuresMetric  prometheus.Counter
	watchDeletionsMetric prometheus.Counter
	failureThreshold     int
	failureCount         int
	Mode                 string
	Namespace            string
}

func NewScheduler(failureThreshold int, every time.Duration, mode, namespace string) *Scheduler {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	watchFailures := promauto.NewCounter(prometheus.CounterOpts{
		Name: lang.PromWatchFailuresName,
		Help: lang.PromWatchFailuresHelp,
	})

	watchDeletions := promauto.NewCounter(prometheus.CounterOpts{
		Name: lang.PromWatchDeletionsName,
		Help: lang.PromWatchDeletionsHelp,
	})

	return &Scheduler{
		Every:                every,
		client:               clientset,
		watchFailuresMetric:  watchFailures,
		watchDeletionsMetric: watchDeletions,
		failureThreshold:     failureThreshold,
		failureCount:         0,
		Mode:                 mode,
		Namespace:            namespace,
	}
}
func (s *Scheduler) Start() {

	ticker := time.NewTicker(time.Duration(s.Every))
	defer ticker.Stop()

	for range ticker.C {
		id := uuid.New()

		go func() {
			s.CreatePod(id.String())

			time.Sleep(10 * time.Second)
			s.CheckPod(id.String())

			// time.Sleep(5 * time.Second)
			// s.DeletePod()
		}()
	}
}

func (s *Scheduler) CreatePod(name string) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("auto-%s", name),
			Namespace: s.Namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sh", "-c", "sleep 5 && exit 0"},
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	_, err := s.client.CoreV1().Pods(s.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		logging.Info(lang.SchedulerAuditorFailedCreation)

		return
	} else {
		logging.Info(lang.SchedulerAuditorSuccessCreation)
	}

}

func (s *Scheduler) CheckPod(name string) {
	_, err := s.client.CoreV1().Pods(s.Namespace).Get(context.TODO(), fmt.Sprintf("auto-%s", name), metav1.GetOptions{})
	if err != nil {
		logging.Info(lang.SchedulerWatcherSuccessDeletion)
	} else {
		logging.Info(lang.SchedulerWatcherFailedDeletion)
		s.watchFailuresMetric.Inc()
		s.failureCount++
	}

	if s.failureCount >= s.failureThreshold {
		if s.Mode == "enforcing" {
			s.DeleteWatcherPod("pepr-system")
		}
		s.failureCount = 0
	}
}

func (s *Scheduler) DeletePod() {
	err := s.client.CoreV1().Pods(s.Namespace).Delete(context.TODO(), "auto-kill-pod", metav1.DeleteOptions{})
	if err != nil {
		logging.Info(lang.SchedulerAuditorFailedDeletion)
		return
	} else {
		logging.Info(lang.SchedulerAuditorSuccessDeletion)
	}
}

func (s *Scheduler) DeleteWatcherPod(namespace string) {
	labelSelector := "pepr.dev/controller=watcher"

	if err := s.client.CoreV1().Pods(namespace).DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: labelSelector,
	}); err != nil {
		logging.Info(lang.SchedulerWatcherPodFailedDeletion)
	} else {
		logging.Info(lang.SchedulerWatcherPodSuccessDeletion)
		s.watchDeletionsMetric.Inc()
	}
}
