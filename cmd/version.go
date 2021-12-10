package cmd

import (
	"fmt"
	
	"github.com/LukasKnuth/EzBackup/restic"

	"github.com/spf13/cobra"
)

// versionCmd represents the scale command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version information of EzBackup and restic",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("EzBackup: 0.1.0") // todo set version in build?
		err := restic.Version()
		if err != nil {
			fmt.Println("Couldn't get restic version. See output...")
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
