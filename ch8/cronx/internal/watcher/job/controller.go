package job

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/time/rate"

	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	batchinformers "k8s.io/client-go/informers/batch/v1"
	batchlisters "k8s.io/client-go/listers/batch/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"cronx/internal/pkg/store"
	"cronx/pkg/logx"
)

// Controller watches Kubernetes Job events (Add/Update/Delete),
// processes them using a workqueue, and synchronizes the Job status
// back into the Cronx database.
//
// This implements the standard controller pattern from client-go:
// - informer watches resources
// - handler enqueues resource keys
// - workers dequeue keys and reconcile state with DB
type Controller struct {
	// Kubernetes Job informer components.
	jobsInformer cache.SharedIndexInformer
	jobsSynced   cache.InformerSynced
	jobsLister   batchlisters.JobLister

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.TypedRateLimitingInterface[cache.ObjectName]

	// DB job store for syncing status.
	jobStore store.JobStore
}

// NewController constructs a Kubernetes Job controller with event handlers.
func NewController(jobInformer batchinformers.JobInformer, jobStore store.JobStore) *Controller {
	ctx := context.Background()

	// Use exponential backoff and token bucket to protect the system from overload.
	ratelimiter := workqueue.NewTypedMaxOfRateLimiter(
		workqueue.NewTypedItemExponentialFailureRateLimiter[cache.ObjectName](5*time.Millisecond, 1000*time.Second),
		&workqueue.TypedBucketRateLimiter[cache.ObjectName]{Limiter: rate.NewLimiter(rate.Limit(50), 300)},
	)

	ctrl := &Controller{
		jobsLister: jobInformer.Lister(),
		jobsSynced: jobInformer.Informer().HasSynced,
		workqueue:  workqueue.NewTypedRateLimitingQueue(ratelimiter),
		jobStore:   jobStore,
	}

	// Register event handlers for Jobs.
	if _, err := jobInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			logx.Info(ctx, "AddFunc", "obj", obj.(*batchv1.Job).Name)
			ctrl.enqueueObj(obj)
		},
		UpdateFunc: func(oldObj, newObj any) {
			newJob := newObj.(*batchv1.Job)
			oldJob := oldObj.(*batchv1.Job)
			// Skip if RV is identical — not a real update.
			if newJob.ResourceVersion == oldJob.ResourceVersion {
				logx.Info(ctx, "old job equals new job", "oldJob", oldJob.Name, "newJob", newJob.Name)
				// Periodic resync will send update events for all known Jobs.
				// Two different versions of the same Jobs will always have different RVs.
				return
			}
			logx.Info(ctx, "UpdateFunc", "oldJob", oldJob.Name, "newJob", newJob.Name)
			ctrl.enqueueObj(newObj)
		},
		DeleteFunc: func(obj any) {
			logx.Info(ctx, "DeleteFunc", "obj", obj.(*batchv1.Job).Name)
			ctrl.enqueueObj(obj)
		}, // optional: we can also handle deletes
	}); err != nil {
		logx.Error(ctx, err.Error())
	}

	return ctrl
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(ctx context.Context, workers int) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Wait for the caches to be synced before starting workers
	if ok := cache.WaitForCacheSync(ctx.Done(), c.jobsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	// Launch two workers to process Job resources
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runWorker, time.Second)
	}

	logx.Info(ctx, "Started workers")
	<-ctx.Done()
	logx.Info(ctx, "Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker(ctx context.Context) {
	ctx = logx.WithTraceID(ctx, uuid.NewString())
	for c.processNextWorkItem(ctx) {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem(ctx context.Context) bool {
	objRef, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}

	// We call Done at the end of this func so the workqueue knows we have
	// finished processing this item. We also must remember to call Forget
	// if we do not want this work item being re-queued. For example, we do
	// not call Forget if a transient error occurs, instead the item is
	// put back on the workqueue and attempted again after a back-off
	// period.
	defer c.workqueue.Done(objRef)

	// Run the syncHandler, passing it the structured reference to the object to be synced.
	err := c.syncHandler(ctx, objRef)
	if err == nil {
		// If no error occurs then we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(objRef)
		logx.Info(ctx, "Successfully synced", "objectName", objRef)
		return true
	}
	// there was a failure so be sure to report it.  This method allows for
	// pluggable error handling which can be used for things like
	// cluster-monitoring.
	utilruntime.HandleErrorWithContext(ctx, err, "Error syncing; requeuing for later retry", "objectReference", objRef)
	// since we failed, we should requeue the item to work on later.  This
	// method will add a backoff to avoid hotlooping on particular items
	// (they're probably still not going to work right away) and overall
	// controller protection (everything I've done is broken, this controller
	// needs to calm down or it can starve other useful work) cases.
	c.workqueue.AddRateLimited(objRef)
	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Job resource
// with the current status of the resource.
func (c *Controller) syncHandler(ctx context.Context, objectRef cache.ObjectName) error {
	logx.Info(ctx, "sync job status to db", "objectRef", objectRef)

	// Get the Job resource with this namespace/name.
	job, err := c.jobsLister.Jobs(objectRef.Namespace).Get(objectRef.Name)
	if err != nil {
		// The Job resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleErrorWithContext(ctx, err, "Job referenced by item in work queue no longer exists", "objectReference", objectRef)
			return nil
		}

		return err
	}

	// Convert k8s Job status to db job status.
	jobStatus := toJobStatus(job)

	// Get job ID from label.
	jobIDStr := job.Labels[LabelJobID]
	jobID, _ := strconv.ParseInt(jobIDStr, 10, 64)

	// Load DB Job.
	modJob, err := c.jobStore.Get(ctx, jobID)
	if err != nil {
		return err
	}

	// Skip if already up-to-date.
	if modJob.Status == jobStatus {
		logx.Info(ctx, "Nothing to update", "namespace", modJob.Namespace, "name", modJob.Name, "status", modJob.Status)
		return nil // no-op
	}

	// CAS update to prevent race conditions.
	if err := c.jobStore.UpdateStatus(ctx, jobID, jobStatus, modJob.Status); err != nil {
		return fmt.Errorf("failed update job status: %w", err)
	}

	return nil
}

// enqueueObj converts an informer object into a queue reference key.
func (c *Controller) enqueueObj(obj any) {
	if objectRef, err := cache.ObjectToName(obj); err != nil {
		utilruntime.HandleError(err)
		return
	} else {
		c.workqueue.Add(objectRef)
	}
}
