package operations

import (
	"fmt"
	"time"

	"github.com/LukasKnuth/EzBackup/k8s"
	"github.com/LukasKnuth/EzBackup/util"
)

func ScaleDown(pvcName string, options *k8s.RequestOptions, force bool, timeout time.Duration) (util.SpecificSet, error) {
	fmt.Println("Looking for Pods mounting the PVC...")
	filtered, err := k8s.FindMountingPods(pvcName, options)
	if err != nil {
		return nil, err
	}

	tree := util.MakeSet(len(filtered)) // could be smaller, but this is a good guess...

	for _, pod := range filtered {
		fmt.Printf("Mounted by %s: %s\n", "Pod", pod.Name)	
		podDependencies, err := k8s.DependencyTree(&pod, options)
		if err != nil {
			if force {
				fmt.Printf("WARN: %s -- Ignoring and continuing because --force was specified\n", err)
			} else {
				return nil, err
			}
		} else {
			tree.PutAll(podDependencies)
		}
	}

	// Output only.
	fmt.Printf("Should scale %d resources\n", len(tree))
	for _, res := range tree {
		fmt.Printf("  %s: %s\n", res.Kind(), res.Name())
	}

	if len(tree) > 0 {
		term_signal := k8s.AwaitTermination(filtered, options, timeout)
		for _, res := range tree {
			err := res.Surrender(options, force)
			if err != nil {
				return nil, err // todo do we need to stop the await here?
			}
		}
		if <-term_signal == k8s.Timeout {
			fmt.Println("Timed out waiting for Pods to go down!")
			// todo do we restore everything here?
		}
	}
	fmt.Println("All dependencies shut down, continuing...")
	return tree, nil
}

func ScaleUp(owners util.SpecificSet, options *k8s.RequestOptions) (err error) {
	fmt.Printf("Scaling %d resources back up\n", len(owners))
	for _, res := range owners {
		fmt.Printf("  %s: %s\n", res.Kind(), res.Name())
		err = res.Restore(options)
		if err != nil {
			return
		}
	}
	return nil
}