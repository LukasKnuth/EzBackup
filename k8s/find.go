package k8s

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// todo allow annotating PVC with "readSafe" (or something) to show that the PVC doesn't require scaling down accessing PVCs

func FindMountingPods(pvc_name string, options *RequestOptions) ([]corev1.Pod, error) {
	pods, err := options.Clientset.CoreV1().Pods(options.Namespace).List(options.Context, metav1.ListOptions{})
	if err != nil {
		return nil, err
	} else {
		return filterMountingPods(pods.Items, pvc_name), nil
	}
}

func filterMountingPods(pods []corev1.Pod, pvcName string) []corev1.Pod {
	var filtered []corev1.Pod
	for _, pod := range pods {
		if podMatches(&pod, pvcName) {
			filtered = append(filtered, pod)
		}
	}
	return filtered
}

func podMatches(pod *corev1.Pod, pvcName string) bool {
	for _, vol := range pod.Spec.Volumes {
		if vol.VolumeSource.PersistentVolumeClaim != nil && vol.VolumeSource.PersistentVolumeClaim.ClaimName == pvcName { // is mounting our PVC
			if vol.VolumeSource.PersistentVolumeClaim.ReadOnly == false { // not a read-only mount
				if pod.Status.Phase == "Running" { // Is still running
					return true
				} else {
					fmt.Printf(" -> Skipping Pod: %s because it's not Running\n", pod.Name)
				}
			} else {
				fmt.Printf(" -> Skipping Pod: %s because it mounts ReadOnly\n", pod.Name)
			}
		}
	}
	return false
}
