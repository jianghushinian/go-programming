package job

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	"cronx/internal/pkg/meta"
	"cronx/internal/pkg/model"
	"cronx/internal/pkg/store"
	"cronx/internal/watcher"
	"cronx/pkg/logx"
)

var _ watcher.Watcher = (*jobWatcher)(nil)

// jobWatcher implements the Watcher interface and periodically syncs
// Cronx job records into Kubernetes Job resources. It also launches
// the Kubernetes Job controller to watch and reconcile Job status
// back into the database.
type jobWatcher struct {
	store      store.IStore         // Abstraction over DB stores
	clientset  kubernetes.Interface // Kubernetes API client
	controller *Controller          // Custom Kubernetes Job controller

	wg sync.WaitGroup
}

// Init initialize watcher dependencies.
func (w *jobWatcher) Init(ctx context.Context, config *watcher.Config) error {
	logx.Info(ctx, "Job watcher init")

	w.store = config.Store
	w.clientset = config.Clientset

	// Create a shared informer factory with 30s resync interval.
	// This informer watches Kubernetes Job resources.
	kubeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(
		w.clientset,
		time.Second*30,
		// Restricts informer to only watch Jobs in the "cronx" namespace.
		kubeinformers.WithNamespace("cronx"),
	)

	// Build the custom controller that translates k8s Job events into DB updates.
	w.controller = NewController(kubeInformerFactory.Batch().V1().Jobs(), w.store.Jobs())

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(ctx.done())
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(ctx.Done())

	// Start controller workers
	go func() {
		if err := w.controller.Run(ctx, 2); err != nil {
			logx.Error(ctx, err.Error(), "Error running controller")
		}
	}()

	logx.Info(ctx, "Job watcher initiated")
	return nil
}

// Spec returns the cron expression for running the watcher.
// Here the watcher runs every 30 seconds.
func (w *jobWatcher) Spec() string {
	return "@every 30s"
}

// Run executes the periodic sync logic:
// - Find all jobs with `Normal` status in DB
// - Create corresponding Kubernetes Jobs
// - Update DB job status to Pending
func (w *jobWatcher) Run() {
	ctx := context.Background()
	ctx = logx.WithTraceID(ctx, uuid.NewString())

	logx.Debug(ctx, "Job sync period is start")

	// Find DB jobs that should be created in Kubernetes.
	_, modJobs, err := w.store.Jobs().List(ctx, meta.WithFilter(map[string]any{
		"status": model.JobStatusNormal,
	}))
	if err != nil {
		logx.Error(ctx, err.Error(), "Failed to list jobs")
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(modJobs))

	// Create Kubernetes Jobs concurrently.
	for _, modJob := range modJobs {
		go func(modJob *model.Job) {
			defer wg.Done()
			ctx := logx.WithJobID(ctx, modJob.ID)

			// Create a Kubernetes Job resource based on DB model.
			job, err := w.clientset.BatchV1().Jobs(modJob.Namespace).Create(ctx, toK8sJob(modJob), metav1.CreateOptions{})
			if err != nil {
				logx.Error(ctx, err.Error(), "job", modJob.Name)
				return
			}

			// Mark status as Pending after creation.
			modJob.Status = model.JobStatusPending
			if err := w.store.Jobs().Update(ctx, modJob); err != nil {
				logx.Error(ctx, err.Error(), "job", modJob.Name)
				return
			}
			logx.Info(ctx, "Successfully created job", "namespace", job.Namespace, "name", job.Name)
		}(modJob)
	}

	wg.Wait()

	logx.Debug(ctx, "Job sync period is complete")
}

func init() {
	// Register this watcher so it is automatically discovered and executed.
	watcher.Register(&jobWatcher{})
}
