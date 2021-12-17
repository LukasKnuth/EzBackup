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
	
	kubeconfigUsage := "If specified, uses the given Kubeconfig file instead of the in-cluster auto configuration "+
	"to talk to the Kubernetes API Server.\nThis is helpful for testing from outside the cluster, for example."
	rootCmd.PersistentFlags().StringVar(&Flags.Kubeconfig, "kubeconfig", "", kubeconfigUsage)

	rootCmd.PersistentFlags().StringVarP(&Flags.Namespace, "namespace", "n", "default", "The Kubernetes namespace to work inside of")

	forceUsage := "When forcing, the command will continue if resources mounting "+
	"the PVC can't be scaled or are of an unsupported Kind."
	rootCmd.PersistentFlags().BoolVar(&Flags.Force, "force", false, forceUsage)

	timeoutUsage := "Specifies a duration to wait for active Pods with write mount to shut down."
	rootCmd.PersistentFlags().DurationVar(&Flags.Timeout, "timeout", 2 * time.Minute, timeoutUsage)
}
