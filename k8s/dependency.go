package k8s

import (
	"strings"
	"fmt"
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	batchv1 "k8s.io/api/batch/v1"
)

type Dependency interface {
	MetaName() string
	MetaKind() string
	Owners() []metav1.OwnerReference
}

type pod struct { *corev1.Pod }

func (p pod) MetaName() string {
	return p.Name
}
func (p pod) MetaKind() string {
	return "Pod"
}
func (p pod) Owners() []metav1.OwnerReference {
	return p.OwnerReferences
}

type deployment struct { *appsv1.Deployment }

func (r deployment) MetaName() string {
	return r.Name
}
func (r deployment) MetaKind() string {
	return "Deployment"
}
func (r deployment) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}

type replicaSet struct { *appsv1.ReplicaSet }

func (r replicaSet) MetaName() string {
	return r.Name
}
func (r replicaSet) MetaKind() string {
	return "ReplicaSet"
}
func (r replicaSet) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}

type job struct { *batchv1.Job }

func (r job) MetaName() string {
	return r.Name
}
func (r job) MetaKind() string {
	return "Job"
}
func (r job) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}

type cronJob struct { *batchv1.CronJob }

func (r cronJob) MetaName() string {
	return r.Name
}
func (r cronJob) MetaKind() string {
	return "CronJob"
}
func (r cronJob) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}

type daemonSet struct { *appsv1.DaemonSet }

func (r daemonSet) MetaName() string {
	return r.Name
}
func (r daemonSet) MetaKind() string {
	return "DaemonSet"
}
func (r daemonSet) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}

type statefulSet struct { *appsv1.StatefulSet }

func (r statefulSet) MetaName() string {
	return r.Name
}
func (r statefulSet) MetaKind() string {
	return "StatefulSet"
}
func (r statefulSet) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}


// ----------- FUNCTIONS -------------

func DependencyTree(mountingPod *corev1.Pod, options *RequestOptions) ([]Dependency, error) {
	tree := make([]Dependency, 0)
	return doDependencyTree(pod{mountingPod}, tree, 1, options)
}

func doDependencyTree(resource Dependency, toScale []Dependency, level int, options *RequestOptions) ([]Dependency, error) {
	if len(resource.Owners()) > 0 {
		for _, owner := range resource.Owners() {
			ownerRes, err := fetchResource(&owner, options)
			if err != nil {
				return nil, err
			} else if ownerRes != nil {
				fmt.Printf("%sowned by %s: %s\n", strings.Repeat(" ", level), ownerRes.MetaKind(), ownerRes.MetaName())
				toScale, err = doDependencyTree(ownerRes, toScale, level + 1, options)
				if err != nil {
					return nil, err
				}
			}
		}
		return toScale, nil
	} else {
		return append(toScale, resource), nil
	}
}

func fetchResource(ref *metav1.OwnerReference, options *RequestOptions) (Dependency, error) {
	switch ref.Kind {
	case "ReplicaSet":
		rs, err := options.Clientset.AppsV1().ReplicaSets(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return replicaSet{rs}, nil
		}
	case "Deployment":
		d, err := options.Clientset.AppsV1().Deployments(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return deployment{d}, nil
		}
	case "Job":
		j, err := options.Clientset.BatchV1().Jobs(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return job{j}, nil
		}
	case "CronJob":
		cj, err := options.Clientset.BatchV1().CronJobs(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return cronJob{cj}, nil
		}
	case "DaemonSet":
		ds, err := options.Clientset.AppsV1().DaemonSets(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return daemonSet{ds}, nil
		}
	case "StatefulSet":
		ss, err := options.Clientset.AppsV1().StatefulSets(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return statefulSet{ss}, nil
		}
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported owner %s: %s\n", ref.Kind, ref.Name))
	}
}