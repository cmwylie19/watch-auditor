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
	Every               time.Duration
	client              *kubernetes.Clientset
	watchFailuresMetric prometheus.Counter
	Namespace           string
}

func NewScheduler(every time.Duration, namespace string) *Scheduler {

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

	return &Scheduler{
		Every:               every,
		client:              clientset,
		watchFailuresMetric: watchFailures,
		Namespace:           namespace,
	}
}
func (s *Scheduler) Start() {

	ticker := time.NewTicker(time.Duration(s.Every))
	defer ticker.Stop()

	for range ticker.C {
		id := uuid.New().String()[:5]

		go func() {
			s.CreatePod(id)

			time.Sleep(10 * time.Second)
			s.CheckPod(id)
		}()
	}
}

func (s *Scheduler) CreatePod(name string) {
	var formattedName = fmt.Sprintf("auto-%s", name)
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      formattedName,
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
		logging.Info(fmt.Sprintf(lang.SchedulerAuditorFailedCreation, formattedName, err.Error()))

		return
	} else {
		logging.Info(fmt.Sprintf(lang.SchedulerAuditorSuccessCreation, formattedName))
	}

}

func (s *Scheduler) CheckPod(name string) {
	var formattedName = fmt.Sprintf("auto-%s", name)
	_, err := s.client.CoreV1().Pods(s.Namespace).Get(context.TODO(), formattedName, metav1.GetOptions{})
	if err != nil {
		logging.Info(fmt.Sprintf(lang.SchedulerWatcherSuccessDeletion, formattedName))
	} else {
		logging.Info(fmt.Sprintf(lang.SchedulerWatcherFailedDeletion, formattedName))
		s.watchFailuresMetric.Inc()
	}
}
