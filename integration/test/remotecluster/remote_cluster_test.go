// +build k8srequired

package remotecluster

import (
	"context"
	"io/ioutil"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/apptest"
	"github.com/giantswarm/apptest/integration/env"
	"github.com/giantswarm/apptest/integration/setup"
)

var (
	config setup.Config
)

func init() {
	var err error

	{
		config, err = setup.NewConfig()
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestRemoteCluster(t *testing.T) {
	var err error

	ctx := context.Background()

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-cluster",
		},
	}

	config.AppTest.CtrlClient().Create(ctx, ns)
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}

	bytes, err := ioutil.ReadFile(env.KubeConfigPath())
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}

	kubeConfig := string(bytes)

	appOperatorVersion := "2.8.0"
	appOperatorValues := `Installation:
  V1:
    Kubernetes:
      API:
        Address: my-cluster.svc.cluster.local`

	apps := []apptest.App{
		{
			AppCRNamespace:     "giantswarm",
			AppOperatorVersion: "0.0.0", // Install app-operator using app-operator-unique.
			CatalogName:        "control-plane-catalog",
			Name:               "app-operator",
			Namespace:          "giantswarm",
			ValuesYAML:         appOperatorValues,
			Version:            "2.8.0",
			WaitForDeploy:      true,
		},
		{
			AppCRNamespace:     "my-cluster",       // CR is created in the remote cluster namespace.
			AppOperatorVersion: appOperatorVersion, // app-operator version to use.
			CatalogName:        "control-plane-test-catalog",
			KubeConfig:         kubeConfig, // Specify a kubeconfig for the remote cluster.
			Name:               "apptest-app",
			Namespace:          "kube-system",
			SHA:                env.CircleSHA(), // The commit to be tested.
			WaitForDeploy:      true,
		},
	}

	err = config.AppTest.InstallApps(ctx, apps)
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}

	// Cleanup deletes all app resources. This is important if the test is run
	// outside of temporary kind clusters.
	err = config.AppTest.CleanUp(ctx, apps)
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}
}
