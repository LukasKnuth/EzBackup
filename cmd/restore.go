package cmd

import (
	"fmt"
	
	"github.com/LukasKnuth/EzBackup/operations"

	"github.com/spf13/cobra"
)

// versionCmd represents the scale command
var restoreCmd = &cobra.Command{
	Use:   "restore [PVC]",
	Short: "Restores a previously created backup.",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		pvcName := args[0]
		fmt.Printf("Command: restore\nPersistent Volume Claim: %s\nNamespace: %s\nTimeout: %s\n\n", pvcName, Flags.Namespace, Flags.Timeout)

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

		err = operations.Restore(pvcName)
		if err != nil {
			fmt.Println("Error while running restore:", err)
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
	rootCmd.AddCommand(restoreCmd)
	// Define any additional flags and configuration settings.
}
