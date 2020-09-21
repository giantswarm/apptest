package apptest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/giantswarm/apiextensions/v2/pkg/apis/application/v1alpha1"
	"github.com/giantswarm/apiextensions/v2/pkg/label"
	"github.com/giantswarm/backoff"
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"gopkg.in/yaml.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	deployedStatus = "deployed"
	namespace      = "giantswarm"
)

// Config represents the configuration used to setup the apps.
type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

// AppSetup implements the logic for managing the app setup.
type AppSetup struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger
}

// New creates a new configured app setup library.
func New(config Config) (*AppSetup, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}

	a := &AppSetup{
		k8sClient: config.K8sClient,
		logger:    config.Logger,
	}

	return a, nil
}

// InstallApps creates appcatalog and app CRs for use in automated tests
// and ensures they are installed by our app platform.
func (a *AppSetup) InstallApps(ctx context.Context, apps []App) error {
	var err error

	err = a.createAppCatalogs(ctx, apps)
	if err != nil {
		return microerror.Mask(err)
	}

	err = a.createApps(ctx, apps)
	if err != nil {
		return microerror.Mask(err)
	}

	err = a.waitForDeployedApps(ctx, apps)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (a *AppSetup) createAppCatalogs(ctx context.Context, apps []App) error {
	var err error

	for _, app := range apps {
		a.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("creating %#q appcatalog cr", app.CatalogName))

		appCatalogCR := &v1alpha1.AppCatalog{
			ObjectMeta: metav1.ObjectMeta{
				Name: app.CatalogName,
				Labels: map[string]string{
					// Processed by app-operator-unique.
					label.AppOperatorVersion: "0.0.0",
				},
			},
			Spec: v1alpha1.AppCatalogSpec{
				Description: app.CatalogName,
				Title:       app.CatalogName,
				Storage: v1alpha1.AppCatalogSpecStorage{
					Type: "helm",
					URL:  app.CatalogURL,
				},
			},
		}
		_, err = a.k8sClient.G8sClient().ApplicationV1alpha1().AppCatalogs().Create(ctx, appCatalogCR, metav1.CreateOptions{})
		if apierrors.IsAlreadyExists(err) {
			a.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("%#q appcatalog CR already exists", appCatalogCR.Name))
		} else if err != nil {
			return microerror.Mask(err)
		}

		a.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("created %#q appcatalog cr", app.CatalogName))
	}

	return nil
}

func (a *AppSetup) createApps(ctx context.Context, apps []App) error {
	for _, app := range apps {
		a.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("creating %#q app cr", app.Name))

		entry, err := getLatestEntry(ctx, app.CatalogURL, app.Name, app.Version)
		if err != nil {
			return microerror.Mask(err)
		}

		appCR := &v1alpha1.App{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Name,
				Namespace: app.Namespace,
				Labels: map[string]string{
					// Processed by app-operator-unique.
					label.AppOperatorVersion: "0.0.0",
				},
			},
			Spec: v1alpha1.AppSpec{
				Catalog: app.CatalogName,
				KubeConfig: v1alpha1.AppSpecKubeConfig{
					InCluster: true,
				},
				Name:      app.Name,
				Namespace: app.Namespace,
				Version:   entry.Version,
			},
		}
		_, err = a.k8sClient.G8sClient().ApplicationV1alpha1().Apps(namespace).Create(ctx, appCR, metav1.CreateOptions{})
		if apierrors.IsAlreadyExists(err) {
			a.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("%#q app CR already exists", appCR.Name))
			return nil
		} else if err != nil {
			return microerror.Mask(err)
		}

		a.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("created %#q app cr", appCR.Name))
	}

	return nil
}

func (a *AppSetup) waitForDeployedApps(ctx context.Context, apps []App) error {
	for _, app := range apps {
		if app.WaitForDeploy {
			err := a.waitForDeployedApp(ctx, app.Name)
			if err != nil {
				return microerror.Mask(err)
			}
		} else {
			a.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("skipping wait for deploy of %#q app cr", app.Name))
		}
	}

	return nil
}

func (a *AppSetup) waitForDeployedApp(ctx context.Context, appName string) error {
	a.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensuring %#q app CR is %#q", appName, deployedStatus))

	o := func() error {
		app, err := a.k8sClient.G8sClient().ApplicationV1alpha1().Apps(namespace).Get(ctx, appName, metav1.GetOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
		if app.Status.Release.Status != deployedStatus {
			return microerror.Maskf(executionFailedError, "waiting for %#q, current %#q", deployedStatus, app.Status.Release.Status)
		}
		return nil
	}

	n := func(err error, t time.Duration) {
		a.logger.Log("level", "debug", "message", fmt.Sprintf("failed to get app CR status '%s': retrying in %s", deployedStatus, t), "stack", fmt.Sprintf("%v", err))
	}

	b := backoff.NewConstant(20*time.Minute, 60*time.Second)
	err := backoff.RetryNotify(o, b, n)
	if err != nil {
		return microerror.Mask(err)
	}

	a.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("ensured %#q app CR is deployed", appName))

	return nil
}

func getLatestEntry(ctx context.Context, storageURL, app, appVersion string) (entry, error) {
	index, err := getIndex(storageURL)
	if err != nil {
		return entry{}, microerror.Mask(err)
	}

	entries, ok := index.Entries[app]
	if !ok {
		return entry{}, microerror.Maskf(notFoundError, "no app %#q in index.yaml", app)
	}

	var latestCreated *time.Time
	var latestEntry entry
	for _, e := range entries {
		if appVersion != "" && e.AppVersion != appVersion {
			continue
		}

		t, err := parseTime(e.Created)
		if err != nil {
			return entry{}, microerror.Mask(err)
		}

		if latestCreated == nil || t.After(*latestCreated) {
			latestCreated = t
			latestEntry = e
			continue
		}
	}

	if latestEntry.Name != "" {
		return latestEntry, nil
	}

	return entry{}, microerror.Maskf(notFoundError, "no app %#q in index.yaml with given appVersion %#q", app, appVersion)
}

func getIndex(storageURL string) (index, error) {
	indexURL := fmt.Sprintf("%s/index.yaml", storageURL)

	// We use https in catalog URLs so we can disable the linter in this case.
	resp, err := http.Get(indexURL) // #nosec
	if err != nil {
		return index{}, microerror.Mask(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return index{}, microerror.Mask(err)
	}

	var i index
	err = yaml.Unmarshal(body, &i)
	if err != nil {
		return i, microerror.Mask(err)
	}

	return i, nil
}

func parseTime(created string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, created)
	if err != nil {
		return nil, microerror.Maskf(executionFailedError, "wrong timestamp format %#q", created)
	}
	return &t, nil
}
