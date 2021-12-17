package cmd

import (
	"fmt"
	"os"

	"github.com/LukasKnuth/EzBackup/k8s"
	"github.com/LukasKnuth/EzBackup/operations"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup [PVC]",
	Short: "Finds any Pods accessing the given PVC, scales them down, runs the backup and scales them back up.",
	Long: `Depending on the data that is written to the Persistent Volume Claim, it can be required to stop the Pod
writing to it first. An example is a database, which keeps data in-memory. A relyable way to flush this
data to disc is to shut the application down.

For other use-cases like static asset hosting, this might not be required.`,
// todo document restic specific environment variables. Need to use other dependency to automate this???
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pvcName := args[0]
		fmt.Printf("Command: backup\nPersistent Volume Claim: %s\nNamespace: %s\nTimeout: %s\n\n", pvcName, Namespace, Timeout)

		options, err := k8s.InCluster(Namespace)
		if err != nil {
			fmt.Println("Couldn't establish connection to Kubernetes API server", err)
			os.Exit(1)
		}

		owners, err := operations.ScaleDown(pvcName, options, Force, Timeout)
		if err != nil {
			fmt.Println("Error while scaling down resources: ", err)
			os.Exit(1)
		}

		err = operations.Backup(pvcName)
		if err != nil {
			fmt.Println("Error while running backup: ", err)
			fmt.Println("Proceeding to scale-up to restore cluster state...")
		}

		err = operations.ScaleUp(owners, options)
		if err != nil {
			fmt.Println("Error while scaling up resources: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	// Define any additional flags and configuration settings.
}
