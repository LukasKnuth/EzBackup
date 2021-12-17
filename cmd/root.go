package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/LukasKnuth/EzBackup/util"
)

var Flags util.Flags

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "EzBackup",
	Short: "Easy backups of Kubernetes Persistent Volume Claim contents.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Cobra also local flags "Flags()" for only this command
	// Or global flags "PersistentFlags()" for all commands
	// todo add a flag to change in-cluster config to out-cluster and specify the kubefile file.

	rootCmd.PersistentFlags().StringVarP(&Flags.Namespace, "namespace", "n", "default", "The Kubernetes namespace to work inside of")

	forceUsage := "When forcing, the command will continue if resources mounting "+
	"the PVC can't be scaled or are of an unsupported Kind."
	rootCmd.PersistentFlags().BoolVar(&Flags.Force, "force", false, forceUsage)

	timeoutUsage := "Specifies a duration to wait for active Pods with write mount to shut down."
	rootCmd.PersistentFlags().DurationVar(&Flags.Timeout, "timeout", 2 * time.Minute, timeoutUsage)
}
