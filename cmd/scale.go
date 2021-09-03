package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// scaleCmd represents the scale command
var scaleCmd = &cobra.Command{
	Use:   "scale [PVC]",
	Short: "Finds any Pods accessing the given PVC and sacles them down through their Deployments/ReplicationSets",
	Long: `Depending on the data that is written to the Persistent Volume Claim, it can be required to stop the Pod
writing to it first. An example is a database, which keeps data in-memory. A relyable way to flush this
data to disc is to shut the application down.

For other use-cases like static asset hosting, this might not be required.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Attempting to scale down anything that uses %s\n", args[0])

		config, err := clientcmd.BuildConfigFromFlags("", "/Users/lukasknuth/k3sup/kubeconfig")
		if err != nil {
			panic(err.Error())
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

		// todo add namespace flag!
		// todo use another "context"!?

		// ######
		pods, err := clientset.CoreV1().Pods("infrastructure").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		filtered := filterMountingPods(pods.Items, args[0])
		tree := make([]metav1.ObjectMeta, 0)
		for _, pod := range filtered {
			fmt.Printf("Mounted by  %s: %s\n", "Pod", pod.Name)
			options := RequestOptions{Namespace: "infrastructure", Context: context.TODO(), Clientset: clientset}
			tree, err = dependencyTree(pod.ObjectMeta, tree, 1, &options)
			if err != nil {
				panic(err.Error())
			}
		}
		fmt.Printf("Should scale %d resources\n", len(tree))
		for _, res := range tree {
			fmt.Printf("  %s\n", res.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(scaleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scaleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//scaleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func filterMountingPods(pods []corev1.Pod, pvcName string) []corev1.Pod {
	var filtered []corev1.Pod
	for _, pod := range pods {
		for _, vol := range pod.Spec.Volumes { // todo also filter by state + filter yourself (by hostname)
			if vol.VolumeSource.PersistentVolumeClaim != nil && vol.VolumeSource.PersistentVolumeClaim.ClaimName == pvcName {
				filtered = append(filtered, pod)
			}
		}
	}
	return filtered
}

type RequestOptions struct {
	Clientset *kubernetes.Clientset
	Context context.Context
	Namespace string
}

func dependencyTree(resource metav1.ObjectMeta, toScale []metav1.ObjectMeta, level int, options *RequestOptions) ([]metav1.ObjectMeta, error) {
	if len(resource.OwnerReferences) > 0 {
		for _, owner := range resource.OwnerReferences {
			ownerRes, err := fetchResource(&owner, options)
			if err != nil {
				return nil, err
			} else if ownerRes != nil {
				fmt.Printf("%s-> %s (%s)\n", strings.Repeat(" ", level), ownerRes.Name, "todo") // todo use wrappers instead?
				toScale, err = dependencyTree(*ownerRes, toScale, level + 1, options)
				if err != nil {
					return nil, err
				}
			} else {
				fmt.Printf("%s-> Unsupported owner %s of type %s\n", strings.Repeat(" ", level), owner.Name, owner.Kind)
			}
		}
		return toScale, nil
	} else {
		return append(toScale, resource), nil
	}
}

func fetchResource(ref *metav1.OwnerReference, options *RequestOptions) (*metav1.ObjectMeta, error) {
	switch ref.Kind {
	case "ReplicaSet":
		rs, err := options.Clientset.AppsV1().ReplicaSets(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return &rs.ObjectMeta, nil
		}
	case "Deployment":
		d, err := options.Clientset.AppsV1().Deployments(options.Namespace).Get(options.Context, ref.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		} else {
			return &d.ObjectMeta, nil
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