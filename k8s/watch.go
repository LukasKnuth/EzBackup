package k8s

import (
	"fmt"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AwaitTermination(pods []corev1.Pod, options *RequestOptions) (<-chan string) {
	lookup := makeLookup(pods)
	informer := makeInformer(options)
	stopper := make(chan struct{})
	done := make(chan string) // todo when timeout is here, send success info!

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			podName := obj.(metav1.Object).GetName()
			_, ok := lookup[podName]
			if ok {
				fmt.Printf("Pod: %s was deleted\n", podName)
			}
			delete(lookup, podName)
			if len(lookup) == 0 {
				close(stopper)
				done <- "success"
			}
		},
	})

	// todo add timeout option (and parameter) and close stopper when time runs out.
	// Blocks until "stopper" channel is closed
	go informer.Run(stopper)
	return done
}

func makeInformer(options *RequestOptions) cache.SharedIndexInformer {
	factory := informers.NewSharedInformerFactory(options.Clientset, time.Second)
	return factory.Core().V1().Pods().Informer()
}

func makeLookup(pods []corev1.Pod) map[string]bool {
	lookup := make(map[string]bool)
	for _, pod := range pods {
		lookup[pod.Name] = true
	}
	return lookup
}