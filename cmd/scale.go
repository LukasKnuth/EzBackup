package cmd

import (
	"context"
	"fmt"

	"github.com/LukasKnuth/EzBackup/k8s"

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
		fmt.Printf("Attempting to scale down anything that uses PVC: %s\n", args[0])

		config, err := clientcmd.BuildConfigFromFlags("", "/Users/lukasknuth/k3sup/kubeconfig") // todo change to in-cluster config
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
		options := k8s.RequestOptions{Namespace: "infrastructure", Context: context.TODO(), Clientset: clientset}
		tree := make([]k8s.Dependency, 0)

		for _, pod := range filtered {
			fmt.Printf("Mounted by %s: %s\n", "Pod", pod.Name)	
			podDependencies, err := k8s.DependencyTree(&pod, &options)
			if err != nil {
				panic(err.Error())
			} else {
				tree = append(tree, podDependencies...)
			}
		}
		fmt.Printf("Should scale %d resources\n", len(tree))
		for _, res := range tree {
			fmt.Printf("  %s: %s\n", res.MetaKind(), res.MetaName())
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
