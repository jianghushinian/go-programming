package job

import (
	"strconv"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cronx/internal/pkg/model"
)

const (
	LabelAPP   = "app"
	LabelJobID = "cronx-job-id"
)

// toK8sJob converts job model into a Kubernetes batch/v1 Job object.
// - Sets basic metadata
// - Defines Pod template
// - Injects identifying labels for sync back
func toK8sJob(job *model.Job) *batchv1.Job {
	backoffLimit := int32(1)
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: job.Name,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						LabelAPP:   "cronx-job",
						LabelJobID: strconv.Itoa(int(job.ID)), // Used to map back to DB job
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "job",
							Image:           job.Info.Image,
							Command:         job.Info.Command,
							Args:            job.Info.Args,
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backoffLimit,
		},
	}
	return jobSpec
}

// toJobStatus inspects Kubernetes Job conditions and maps them to DB status.
func toJobStatus(job *batchv1.Job) model.JobStatus {
	switch {
	case isJobCompleted(job):
		return model.JobStatusSucceeded
	case isJobFailed(job):
		return model.JobStatusFailed
	case isJobSuspended(job):
		return model.JobStatusPending
	case isJobFailureTarget(job):
		return model.JobStatusFailed
	case isJobSuccessCriteriaMet(job):
		return model.JobStatusSucceeded
	case isJobActive(job):
		return model.JobStatusRunning
	default:
		return model.JobStatusUnknown
	}
}

func isJobCompleted(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func isJobFailed(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobFailed && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func isJobSuspended(job *batchv1.Job) bool {
	return job.Spec.Suspend != nil && *job.Spec.Suspend
}

func isJobFailureTarget(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobFailureTarget && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func isJobSuccessCriteriaMet(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobSuccessCriteriaMet && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func isJobActive(job *batchv1.Job) bool {
	return job.Status.Active > 0
}
