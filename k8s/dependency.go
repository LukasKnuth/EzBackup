package k8s

import (
	"strings"
	"fmt"
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

func (p *pod) Name() string {
	return p.Pod.Name
}
func (p *pod) Kind() string {
	return "Pod"
}
func (p *pod) Owners() []metav1.OwnerReference {
	return p.OwnerReferences
}

func (r *deployment) Name() string {
	return r.Deployment.Name
}
func (r *deployment) Kind() string {
	return "Deployment"
}
func (r *deployment) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}

func (r *replicaSet) Name() string {
	return r.ReplicaSet.Name
}
func (r *replicaSet) Kind() string {
	return "ReplicaSet"
}
func (r *replicaSet) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}

func (r *job) Name() string {
	return r.Job.Name
}
func (r *job) Kind() string {
	return "Job"
}
func (r *job) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}

func (r *cronJob) Name() string {
	return r.CronJob.Name
}
func (r *cronJob) Kind() string {
	return "CronJob"
}
func (r *cronJob) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}

func (r *daemonSet) Name() string {
	return r.DaemonSet.Name
}
func (r *daemonSet) Kind() string {
	return "DaemonSet"
}
func (r *daemonSet) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}

func (r *statefulSet) Name() string {
	return r.StatefulSet.Name
}
func (r *statefulSet) Kind() string {
	return "StatefulSet"
}
func (r *statefulSet) Owners() []metav1.OwnerReference {
	return r.OwnerReferences
}


// ----------- FUNCTIONS -------------

func DependencyTree(mountingPod *corev1.Pod, options *RequestOptions) ([]BlockingMountOwner, error) {
	tree := make([]BlockingMountOwner, 0)
	return doDependencyTree(&pod{mountingPod}, tree, 1, options)
}

func doDependencyTree(resource BlockingMountOwner, toScale []BlockingMountOwner, level int, options *RequestOptions) ([]BlockingMountOwner, error) {
	if len(resource.Owners()) > 0 {
		for _, owner := range resource.Owners() {
			ownerRes, err := fetchResource(&owner, options)
			if err != nil {
				return nil, err
			} else if ownerRes != nil {
				fmt.Printf("%sowned by %s: %s\n", strings.Repeat(" ", level), ownerRes.Kind(), ownerRes.Name())
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

func fetchResource(ref *metav1.OwnerReference, options *RequestOptions) (BlockingMountOwner, error) {
	switch ref.Kind {
	case "ReplicaSet":
		rs, err := options.Clientset.AppsV1().ReplicaSets(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return &replicaSet{rs, *rs.Spec.Replicas}, nil
		}
	case "Deployment":
		d, err := options.Clientset.AppsV1().Deployments(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return &deployment{d, *d.Spec.Replicas}, nil
		}
	case "Job":
		j, err := options.Clientset.BatchV1().Jobs(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return &job{j}, nil
		}
	case "CronJob":
		cj, err := options.Clientset.BatchV1().CronJobs(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return &cronJob{cj}, nil
		}
	case "DaemonSet":
		ds, err := options.Clientset.AppsV1().DaemonSets(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return &daemonSet{ds}, nil
		}
	case "StatefulSet":
		ss, err := options.Clientset.AppsV1().StatefulSets(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return &statefulSet{ss}, nil
		}
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported owner %s: %s", ref.Kind, ref.Name))
	}
}