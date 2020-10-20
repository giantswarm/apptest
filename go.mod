module github.com/giantswarm/apptest

go 1.14

require (
	github.com/giantswarm/apiextensions/v3 v3.2.0
	github.com/giantswarm/appcatalog v0.2.7
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/k8sclient/v4 v4.0.1-0.20201020101418-0115e035b0dc
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/micrologger v0.3.3
	k8s.io/api v0.18.9
	k8s.io/apimachinery v0.18.9
)

replace sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
