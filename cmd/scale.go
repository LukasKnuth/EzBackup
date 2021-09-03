package cmd

import (
	"fmt"

	"github.com/LukasKnuth/EzBackup/k8s"

	"github.com/spf13/cobra"
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

		// todo add namespace flag!

		options, err := k8s.FromKubeconfig("/Users/lukasknuth/k3sup/kubeconfig", "infrastructure") // todo change to in-cluster config
		if err != nil {
			panic(err.Error())
		}

		filtered, err := k8s.FindMountingPods(args[0], options)
		if err != nil {
			panic(err.Error())
		}

		tree := make([]k8s.Dependency, 0)

		for _, pod := range filtered {
			fmt.Printf("Mounted by %s: %s\n", "Pod", pod.Name)	
			podDependencies, err := k8s.DependencyTree(&pod, options)
			if err != nil {
				// todo control via flag what to do: Continue/fail
				fmt.Print(err)
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
