package k8s

import (
	"fmt"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AwaitResult int32
const (
	Completed AwaitResult = iota
	Timeout AwaitResult = iota
)

func AwaitTermination(pods []corev1.Pod, options *RequestOptions, timeout time.Duration) (<-chan AwaitResult) {
	lookup := makeLookup(pods)
	informer := makeInformer(options)
	stopper := make(chan struct{})
	// Buffer one. If timeout stops informer and DeleteFunc wants to send result, it doesn't block but is never used.
	// This prevents the goroutine from lingering, since it blocks on sending forever.
	internalChan := make(chan AwaitResult, 1)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			podName := obj.(metav1.Object).GetName()
			_, ok := lookup[podName]
			if ok {
				fmt.Printf("Pod: %s was deleted\n", podName)
			}
			delete(lookup, podName)
			if len(lookup) == 0 {
				internalChan <- Completed
			}
		},
	})

	resultChan := make(chan AwaitResult)
	go func() {
		select {
		case res := <-internalChan:
			// We're done before the timeout. Return!
			close(stopper)
			resultChan <- res
		case <-time.After(timeout):
			// timeout!
			close(stopper)
			resultChan <- Timeout
		}
	}()

	// Blocks until "stopper" channel is closed
	go informer.Run(stopper)
	return resultChan
}

func makeInformer(options *RequestOptions) cache.SharedIndexInformer {
	opts := informers.WithNamespace(options.Namespace)
	factory := informers.NewSharedInformerFactoryWithOptions(options.Clientset, time.Second, opts)
	return factory.Core().V1().Pods().Informer()
}

func makeLookup(pods []corev1.Pod) map[string]bool {
	lookup := make(map[string]bool)
	for _, pod := range pods {
		lookup[pod.Name] = true
	}
	return lookup
}