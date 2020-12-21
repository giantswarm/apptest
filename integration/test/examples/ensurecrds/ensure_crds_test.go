// +build k8srequired

package ensurecrds

import (
	"context"
	"testing"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/monitoring/v1alpha1"
	"github.com/giantswarm/apiextensions/v3/pkg/crd"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

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

func TestEnsureCRDs(t *testing.T) {
	var err error

	ctx := context.Background()

	crds := []*apiextensionsv1.CustomResourceDefinition{
		crd.LoadV1("monitoring.giantswarm.io", "Silence"),
	}

	err = config.AppTest.EnsureCRDs(ctx, crds)
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}

	silences := &v1alpha1.SilenceList{}

	err = config.AppTest.CtrlClient().List(ctx, silences)
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}

	if len(silences.Items) != 0 {
		t.Fatalf("expected 0 silences got %d", len(silences.Items))
	}
}
