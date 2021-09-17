package k8s

import (
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ScaleInfo struct {
	replicaCount int32
}

type Scalable interface {
	ScaleDown(options *RequestOptions, force bool) (ScaleInfo, error)
	ScaleUp(options *RequestOptions, info *ScaleInfo) error
}

// Pod
func (p *pod) ScaleDown(options *RequestOptions, force bool) (info ScaleInfo, err error) {
	info = ScaleInfo{}
	if force {
		err = options.Clientset.CoreV1().Pods(options.Namespace).Delete(options.Context, p.Name(), metav1.DeleteOptions{})
	} else {
		err = errors.New("The Pod resource does not support non-destructive scaling.")
	}
	return
}

func (p *pod) ScaleUp(options *RequestOptions, info *ScaleInfo) error {
	fmt.Printf("The Pod resource can't be scaled up. Re-create Pod: %s manually", p.Kind())
	return nil
}

// Deployment
func (d *deployment) ScaleDown(options *RequestOptions, force bool) (info ScaleInfo, err error) {
	info = ScaleInfo{}
	scale, err := options.Clientset.AppsV1().Deployments(options.Namespace).GetScale(options.Context, d.Name(), metav1.GetOptions{})
	if err != nil {
		return
	}
	fmt.Printf("Scaling Pod: %s from %d to 0 replicas", d.Name(), scale.Spec.Replicas)
	info.replicaCount = scale.Spec.Replicas
	scale.Spec.Replicas = 0
	_, err = options.Clientset.AppsV1().Deployments(options.Namespace).UpdateScale(options.Context, d.Name(), scale, metav1.UpdateOptions{})
	return
}

func (d *deployment) ScaleUp(options *RequestOptions, info *ScaleInfo) error {
	scale, err := options.Clientset.AppsV1().Deployments(options.Namespace).GetScale(options.Context, d.Name(), metav1.GetOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Scaling Pod: %s from 0 to %d replicas", d.Name(), info.replicaCount)
	scale.Spec.Replicas = info.replicaCount
	_, err = options.Clientset.AppsV1().Deployments(options.Namespace).UpdateScale(options.Context, d.Name(), scale, metav1.UpdateOptions{})
	return err
}

// ReplicaSet
func (rs *replicaSet) ScaleUp(options *RequestOptions, force bool) (info ScaleInfo, err error) {
	return ScaleInfo{}, nil
}

func (rs *replicaSet) ScaleDown(options *RequestOptions, info *ScaleInfo) error {
	return nil
}

// job
func (j *job) ScaleUp(options *RequestOptions, force bool) (info ScaleInfo, err error) {
	return ScaleInfo{}, nil
}

func (j *job) ScaleDown(options *RequestOptions, info *ScaleInfo) error {
	return nil
}

// cronJob
func (cj *cronJob) ScaleUp(options *RequestOptions, force bool) (info ScaleInfo, err error) {
	return ScaleInfo{}, nil
}

func (cj *cronJob) ScaleDown(options *RequestOptions, info *ScaleInfo) error {
	return nil
}

// daemonSet
func (ds *daemonSet) ScaleUp(options *RequestOptions, force bool) (info ScaleInfo, err error) {
	return ScaleInfo{}, nil
}

func (ds *daemonSet) ScaleDown(options *RequestOptions, info *ScaleInfo) error {
	return nil
}

// statefulSet
