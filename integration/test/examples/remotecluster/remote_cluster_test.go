// +build k8srequired

package remotecluster

import (
	"context"
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
			Name: "eggs2",
		},
	}

	config.AppTest.CtrlClient().Create(ctx, ns)
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}

	apps := []apptest.App{
		{
			AppCRNamespace:     "eggs2", // CR is created in the remote cluster namespace.
			AppOperatorVersion: "0.0.0", // app-operator version to use.
			CatalogName:        "default",
			// KubeConfig:      "kubeconfig" // TODO Sort out kubeconfig.
			Name:          "test-app",
			Namespace:     "kube-system",
			Version:       env.CircleSHA(),
			WaitForDeploy: true,
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
