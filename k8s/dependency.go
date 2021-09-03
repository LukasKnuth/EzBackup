package k8s

import (
	"context"
	"strings"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes"
)

type RequestOptions struct {
	Clientset *kubernetes.Clientset
	Context context.Context
	Namespace string
}

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
				fmt.Printf("%s-> %s: %s\n", strings.Repeat(" ", level), ownerRes.MetaKind(), ownerRes.MetaName())
				toScale, err = doDependencyTree(ownerRes, toScale, level + 1, options)
				if err != nil {
					return nil, err
				}
			} else {
				fmt.Printf("%s-> Unsupported owner %s: %s\n", strings.Repeat(" ", level), owner.Kind, owner.Name)
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
		fallthrough
	case "CronJob":
		fallthrough
	case "DaemonSet":
		fallthrough
	case "StatefulSet":
		fallthrough
	default:
		return nil, nil
	}
}