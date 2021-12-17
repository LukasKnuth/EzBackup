package cmd

import (
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		pvcName := args[0]
		fmt.Printf("Command: backup\nPersistent Volume Claim: %s\nNamespace: %s\nTimeout: %s\n\n", pvcName, Flags.Namespace, Flags.Timeout)

		options, err := operations.AutoConfigure(&Flags)
		if err != nil {
			fmt.Println("Couldn't establish connection to Kubernetes API server:", err)
			return
		}

		owners, err := operations.ScaleDown(pvcName, options, Flags.Force, Flags.Timeout)
		if err != nil {
			fmt.Println("Error while scaling down resources:", err)
			return
		}

		err = operations.Backup(pvcName)
		if err != nil {
			fmt.Println("Error while running backup:", err)
			fmt.Println("Proceeding to scale-up to restore cluster state...")
		}

		err = operations.ScaleUp(owners, options)
		if err != nil {
			fmt.Println("Error while scaling up resources:", err)
			return
		}

		// All good!
		return nil
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	// Define any additional flags and configuration settings.
}
