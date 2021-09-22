package cmd

import (
	"fmt"
	"time"

	"github.com/LukasKnuth/EzBackup/k8s"
	"github.com/LukasKnuth/EzBackup/operations"

	"github.com/spf13/cobra"
)

var Namespace string
var Force bool
var Timeout time.Duration

// backupCmd represents the scale command
var backupCmd = &cobra.Command{
	Use:   "backup [PVC]",
	Short: "Finds any Pods accessing the given PVC, scales them down, runs the backup and scales them back up.",
	Long: `Depending on the data that is written to the Persistent Volume Claim, it can be required to stop the Pod
writing to it first. An example is a database, which keeps data in-memory. A relyable way to flush this
data to disc is to shut the application down.

For other use-cases like static asset hosting, this might not be required.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pvcName := args[0]
		fmt.Printf("Command: backup\nPersistent Volume Claim: %s\nNamespace: %s\nTimeout: %s\n\n", pvcName, Namespace, Timeout)

		options, err := k8s.FromKubeconfig("/Users/lukasknuth/k3sup/kubeconfig", Namespace)
		if err != nil {
			panic(err.Error())
		}

		owners, err := operations.ScaleDown(pvcName, options, Force, Timeout)
		if err != nil {
			fmt.Println("Error while scaling down resources: ", err)
		}

		operations.Backup()

		err = operations.ScaleUp(owners, options)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// todo we probably need to move this into the main file?
	rootCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "default", "The Kubernetes namespace to work inside of")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	forceUsage := "When forcing, the command will continue if resources mounting "+
	"the PVC can't be scaled or are of an unsupported Kind."
	backupCmd.Flags().BoolVar(&Force, "force", false, forceUsage)

	timeoutUsage := "Specifies a duration to wait for active Pods with write mount to shut down."
	backupCmd.Flags().DurationVar(&Timeout, "timeout", 2 * time.Minute, timeoutUsage)
}
