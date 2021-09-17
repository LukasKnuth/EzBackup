package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	batchv1 "k8s.io/api/batch/v1"
)

/// A Kubernetes resource which mounts the PVC with writing permissions. This blocks us from creating a backup
/// because the resource could be writing to the PVC while we're backing it up, causing the backup to be
/// incomplete or invalid.
type BlockingMountOwner interface {
	Name() string
	Kind() string
	Owners() []metav1.OwnerReference
	Surrender(options *RequestOptions, force bool) error
	Restore(options *RequestOptions) error
}

type pod struct { *corev1.Pod }
type deployment struct {
	*appsv1.Deployment
	originalReplicas int32
}
type replicaSet struct {
	*appsv1.ReplicaSet
	originalReplicas int32
}
type job struct { *batchv1.Job }
type cronJob struct { *batchv1.CronJob }
type daemonSet struct { *appsv1.DaemonSet }
type statefulSet struct { *appsv1.StatefulSet }