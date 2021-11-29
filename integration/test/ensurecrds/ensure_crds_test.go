//go:build k8srequired
// +build k8srequired

package ensurecrds

import (
	"context"
	"testing"

	"github.com/giantswarm/app/v6/pkg/crd"
	monitoringv1alpha1 "github.com/giantswarm/silence-operator/api/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"

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

func TestEnsureCRDs(t *testing.T) {
	var err error

	ctx := context.Background()

	runtimeScheme := runtime.NewScheme()

	// Add the extra CRDs you need to the scheme.
	appSchemeBuilder := runtime.SchemeBuilder{
		monitoringv1alpha1.AddToScheme,
	}
	err = appSchemeBuilder.AddToScheme(runtimeScheme)
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}

	c := apptest.Config{
		Logger: config.Logger,
		Scheme: runtimeScheme,

		KubeConfigPath: env.KubeConfigPath(),
	}

	appTest, err := apptest.New(c)
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}

	var crdGetter *crd.CRDGetter
	{
		cc := crd.Config{
			Logger: config.Logger,
		}

		crdGetter, err = crd.NewCRDGetter(cc)
		if err != nil {
			t.Fatalf("expected nil got %#q", err)
		}
	}

	customResourceDefinition, err := crdGetter.LoadCRD(ctx, "monitoring.giantswarm.io", "Silence")
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}

	crds := []*apiextensionsv1.CustomResourceDefinition{
		customResourceDefinition,
	}

	// Ensure the CRD exists in the cluster.
	err = appTest.EnsureCRDs(ctx, crds)
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}

	silences := &monitoringv1alpha1.SilenceList{}

	err = appTest.CtrlClient().List(ctx, silences)
	if err != nil {
		t.Fatalf("expected nil got %#q", err)
	}

	if len(silences.Items) != 0 {
		t.Fatalf("expected 0 silences got %d", len(silences.Items))
	}
}
