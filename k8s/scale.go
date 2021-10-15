package k8s

import (
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pod
func (p *pod) Surrender(options *RequestOptions, force bool) (err error) {
	if force {
		err = options.Clientset.CoreV1().Pods(options.Namespace).Delete(options.Context, p.Name(), metav1.DeleteOptions{})
	} else {
		err = errors.New(fmt.Sprintf("%s: The Pod resource does not support non-destructive scaling.", p.Name()))
	}
	return
}

func (p *pod) Restore(options *RequestOptions) error {
	fmt.Printf("The Pod resource can't be scaled up. Re-create Pod: %s manually\n", p.Kind())
	return nil
}

// Deployment
func (d *deployment) Surrender(options *RequestOptions, force bool) (err error) {
	scale, err := options.Clientset.AppsV1().Deployments(options.Namespace).GetScale(options.Context, d.Name(), metav1.GetOptions{})
	if err != nil {
		return
	}
	fmt.Printf("%s: Scaling Deployment from %d to 0 replicas\n", d.Name(), scale.Spec.Replicas)
	d.originalReplicas = scale.Spec.Replicas
	scale.Spec.Replicas = 0
	_, err = options.Clientset.AppsV1().Deployments(options.Namespace).UpdateScale(options.Context, d.Name(), scale, metav1.UpdateOptions{})
	return
}

func (d *deployment) Restore(options *RequestOptions) error {
	scale, err := options.Clientset.AppsV1().Deployments(options.Namespace).GetScale(options.Context, d.Name(), metav1.GetOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("%s, Scaling Deployment from 0 to %d replicas\n", d.Name(), d.originalReplicas)
	scale.Spec.Replicas = d.originalReplicas
	_, err = options.Clientset.AppsV1().Deployments(options.Namespace).UpdateScale(options.Context, d.Name(), scale, metav1.UpdateOptions{})
	return err
}

// todo all below!
// ReplicaSet
func (rs *replicaSet) Surrender(options *RequestOptions, force bool) (err error) {
	return nil
}

func (rs *replicaSet) Restore(options *RequestOptions) error {
	return nil
}

// job
func (j *job) Surrender(options *RequestOptions, force bool) (err error) {
	return nil
}

func (j *job) Restore(options *RequestOptions) error {
	return nil
}

// cronJob
func (cj *cronJob) Surrender(options *RequestOptions, force bool) (err error) {
	return nil
}

func (cj *cronJob) Restore(options *RequestOptions) error {
	return nil
}

// daemonSet
func (ds *daemonSet) Surrender(options *RequestOptions, force bool) (err error) {
	return nil
}

func (ds *daemonSet) Restore(options *RequestOptions) error {
	return nil
}

// statefulSet
func (ds *statefulSet) Surrender(options *RequestOptions, force bool) (err error) {
	return nil
}

func (ds *statefulSet) Restore(options *RequestOptions) error {
	return nil
}