package scheduler

import (
	"context"
	"fmt"

	"github.com/cmwylie19/watch-auditor/src/config/lang"
	"github.com/cmwylie19/watch-auditor/src/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Scheduler struct {
	Every                int
	Unit                 string
	client               *kubernetes.Clientset
	watchFailuresMetric  prometheus.Counter
	watchDeletionsMetric prometheus.Counter
	failureThreshold     int
	failureCount         int
}

func NewScheduler(failureThreshold, every int, unit string) *Scheduler {

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

	prometheus.MustRegister(watchFailures)
	prometheus.MustRegister(watchDeletions)

	return &Scheduler{
		Every:                every,
		Unit:                 unit,
		client:               clientset,
		watchFailuresMetric:  watchFailures,
		watchDeletionsMetric: watchDeletions,
		failureThreshold:     failureThreshold,
		failureCount:         0,
	}
}
func (s *Scheduler) Start() {
	for {
		s.CreatePod()
		s.CheckPod()
		s.DeletePod()
		if s.failureCount >= s.failureThreshold {
			s.DeleteWatcherPod("pepr-system")
			s.failureCount = 0
		}
	}
}
func (s *Scheduler) CreatePod() {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "auto-kill-pod",
			Namespace: "neuvector",
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

	_, err := s.client.CoreV1().Pods("neuvector").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		logging.Info(lang.SchedulerFailedCreate + err.Error())
		return
	}
}

func (s *Scheduler) CheckPod() {
	_, err := s.client.CoreV1().Pods("neuvector").Get(context.TODO(), "auto-kill-pod", metav1.GetOptions{})
	if err != nil {
		logging.Info(lang.SchedulerWatcherSuccessDeletion)
	} else {
		logging.Info(lang.SchedulerWatcherFailedDeletion)
		s.watchFailuresMetric.Inc()
	}
}

func (s *Scheduler) DeletePod() {
	err := s.client.CoreV1().Pods("neuvector").Delete(context.TODO(), "auto-kill-pod", metav1.DeleteOptions{})
	if err != nil {
		logging.Info(lang.SchedulerFailedDelete + err.Error())
		return
	}
}

func (s *Scheduler) DeleteWatcherPod(namespace string) {
	labelSelector := "pepr.dev/controller=watcher"
	deleteOptions := metav1.DeleteOptions{}
	listOptions := metav1.ListOptions{LabelSelector: labelSelector}
	if err := s.client.CoreV1().Pods(namespace).DeleteCollection(context.TODO(), deleteOptions, listOptions); err != nil {
		logging.Error(fmt.Sprintf(lang.SchedulerWatcherPodFailedDeletion+namespace, err.Error()))
	} else {
		logging.Info(lang.SchedulerWatcherPodSuccessDeletion + namespace)
		s.watchDeletionsMetric.Inc()
	}
}