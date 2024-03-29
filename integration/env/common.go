//go:build k8srequired
// +build k8srequired

package env

import (
	"fmt"
	"os"
)

const (
	// EnvVarCircleSHA is the process environment variable representing the
	// CIRCLE_SHA1 env var.
	EnvVarCircleSHA = "CIRCLE_SHA1"
	// EnvVarE2EKubeconfig is the process environment variable representing the
	// E2E_KUBECONFIG env var.
	EnvVarE2EKubeconfig = "E2E_KUBECONFIG"
)

var (
	circleSHA      string
	kubeconfigPath string
)

func init() {
	circleSHA = os.Getenv(EnvVarCircleSHA)
	if circleSHA == "" {
		panic(fmt.Sprintf("env var '%s' must not be empty", EnvVarCircleSHA))
	}

	kubeconfigPath = os.Getenv(EnvVarE2EKubeconfig)
	if kubeconfigPath == "" {
		panic(fmt.Sprintf("env var '%s' must not be empty", EnvVarE2EKubeconfig))
	}
}

func CircleSHA() string {
	return circleSHA
}

func KubeConfigPath() string {
	return kubeconfigPath
}
