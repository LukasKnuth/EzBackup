package k8s

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
)

type RequestOptions struct {
	Clientset *kubernetes.Clientset
	Context context.Context
	Namespace string
}

func FromKubeconfig(kubeconfig_path string, namespace string) (*RequestOptions, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig_path)
	if err != nil {
		return nil, err
	} else {
		return fromRestConfig(config, namespace)
	}
}

func InCluster(namespace string) (*RequestOptions, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	} else {
		return fromRestConfig(config, namespace)
	}
}

func fromRestConfig(config *rest.Config, namespace string) (*RequestOptions, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	} else {
		return &RequestOptions{
			Clientset: clientset,
			Context: context.TODO(), // todo what do use here?
			Namespace: namespace,
		}, nil
	}
}