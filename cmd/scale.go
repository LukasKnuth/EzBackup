package cmd

import (
	"fmt"

	"github.com/LukasKnuth/EzBackup/k8s"

	"github.com/spf13/cobra"
)

var Namespace string
var Force bool

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
		fmt.Printf("Command: scale\nPersistent Volume Claim: %s\nNamespace: %s\n\n", args[0], Namespace)

		options, err := k8s.FromKubeconfig("/Users/lukasknuth/k3sup/kubeconfig", Namespace)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Looking for Pods mounting the PVC...")
		filtered, err := k8s.FindMountingPods(args[0], options)
		if err != nil {
			panic(err.Error())
		}

		tree := make([]k8s.BlockingMountOwner, 0)

		for _, pod := range filtered {
			fmt.Printf("Mounted by %s: %s\n", "Pod", pod.Name)	
			podDependencies, err := k8s.DependencyTree(&pod, options)
			if err != nil {
				if Force {
					fmt.Printf("WARN: %s -- Ignoring and continuing because --force was specified\n", err)
				} else {
					panic(err.Error())
				}
			} else {
				tree = append(tree, podDependencies...) // todo can have duplicates! If one deployment creates 2 pods, the deployment is here twice!
			}
		}
		fmt.Printf("Should scale %d resources\n", len(tree))
		for _, res := range tree {
			fmt.Printf("  %s: %s\n", res.Kind(), res.Name())
			res.Surrender(options, Force)
		}
		if len(tree) > 0 {
			k8s.AwaitTermination(filtered, options)
			fmt.Println("All dependencies shut down, continuing...")
		}
		fmt.Println("PRETENDING: backup...") // todo how can we "return" here and continue afterwads?
		fmt.Printf("Scaling %d resources back up\n", len(tree))
		for _, res := range tree {
			fmt.Printf("  %s: %s\n", res.Kind(), res.Name())
			res.Restore(options)
		}
	},
}

func init() {
	rootCmd.AddCommand(scaleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// todo we probably need to move this into the main file?
	rootCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "default", "The Kubernetes namespace to work inside of")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	forceUsage := "When forcing, the command will continue if resources mounting "+
	"the PVC can't be scaled or are of an unsupported Kind."
	scaleCmd.Flags().BoolVar(&Force, "force", false, forceUsage)
}
