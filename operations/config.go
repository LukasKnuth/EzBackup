package operations

import (
	"fmt"

	"github.com/LukasKnuth/EzBackup/k8s"
	"github.com/LukasKnuth/EzBackup/util"
)

func AutoConfigure(flags *util.Flags) (*k8s.RequestOptions, error) {
	if flags.Kubeconfig != "" {
		fmt.Println("Config: out-of-cluster mode, using Kubeconfig from %s...", flags.Kubeconfig)
		return k8s.FromKubeconfig(flags.Kubeconfig, flags.Namespace)
	} else {
		fmt.Println("Config: in-cluster mode, auto-detecting configuration...")
		return k8s.InCluster(flags.Namespace)
	}
}