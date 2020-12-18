// +build k8srequired

package upgradeapp

import (
	"context"
	"testing"

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

func TestUpgradeApp(t *testing.T) {
	var err error
	ctx := context.Background()

	currentApp := apptest.App{
		CatalogName:   "default",
		Name:          "test-app",
		Namespace:     "giantswarm",
		WaitForDeploy: true,
	}

	desiredApp := apptest.App{
		CatalogName:   "default-test",
		Name:          "test-app",
		Namespace:     "giantswarm",
		SHA:           env.CircleSHA(),
		WaitForDeploy: true,
	}

	err = config.AppTest.UpgradeApp(ctx, currentApp, desiredApp)
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}
}
