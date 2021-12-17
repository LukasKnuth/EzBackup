package util

import (
	"time"
)

type Flags struct {
	Namespace string
	Force bool
	Kubeconfig string
	Timeout time.Duration
}