# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



## [Unreleased]

## [1.3.0] - 2023-11-10

### Changed

- Add a switch for PSP CR installation.

## [1.2.1] - 2022-07-28

### Fixed

- Catalog CR is fixed to include "Repositories" section

## [1.2.0] - 2022-03-03

### Added

- Add constants for catalog names.

## [1.1.1] - 2022-03-01

### Fixed

- Add retry logic when accessing catalog repos to get app version.

## [1.1.0] - 2022-02-21

### Changed

- Downgrade k8s modules to `< 0.21.0` version and controller-runtime to `< 0.7.0` version.

## [1.0.1] - 2022-01-24

### Fixed

- Create `Catalog` CRs instead of deprecated `AppCatalog` CRs.

## [1.0.0] - 2021-11-29

### Changed

- Drop `apiextensions` dependency.

## [0.12.0] - 2021-08-24

### Added

- Add `RESTConfig()` for use in integration tests.

## [0.11.0] - 2021-06-16

### Added

- Create `Catalog` CRs for the integration test.

## [0.10.3] - 2021-03-09

### Fixed

- Fix to `appcatalog` library for getting the latest entry in the app catalog.
- Update user values configmap if it already exists.

## [0.10.2] - 2021-02-08

### Fixed

- Extend support for setting the app CR name.

## [0.10.1] - 2021-02-08

### Fixed

- Revert `sigs.k8s.io/controller-runtime` to v0.6.4.

## [0.10.0] - 2021-02-03

### Added

- Add support for setting the app CR name.
- Set `app.kubernetes.io/name` label for app CRs.

## [0.9.0] - 2020-12-15

### Added

- Adding App Upgrade test.

## [0.8.2] - 2020-12-14

### Fixed

- Fix namespace handling when waiting for deployed apps.

## [0.8.1] - 2020-12-14

### Fixed

- Fix namespace handling when waiting for deployed apps.

## [0.8.0] - 2020-12-10

### Added

- Add support for setting kubeconfig secret in app CRs for remote clusters.
- Add clean up function to remove resources created while installing apps.

### Changed

- User config values are created on App namespace.

## [0.7.1] - 2020-11-26

### Fixed

- Comparing `SHA` parameter with either app version or version.

## [0.7.0] - 2020-11-17

### Fixed

- Install specified app version instead of latest when passing the `SHA` parameter.

### Changed

- Updated `appcatalog` library.

## [0.6.0] - 2020-11-12

### Changed

- Don't fail when ensuring a CRD that's already present.

## [0.5.0] - 2020-11-06

### Fixed

- Add new methods to the interface.

### Added

- Expose method `EnsureCRDs` to register CRDs in the k8s API.
- A custom `Scheme` can be passed to configure the controller-runtime client.
- Add getter method that returns the controller-runtime client.
- Generate catalog URLs for known catalogs.

## [0.4.1] - 2020-10-30

### Added

- Support both explicit kubeconfigs and file paths.

### Fixed

- Optimize apps wait interval as app-operator has a status webhook.

## [0.4.0] - 2020-10-29

### Changed

- Remove k8sclient dependency and use controller-runtime client for managing CRs.

## [0.3.0] - 2020-10-08

### Added

- Add support for configuring app CRs with values.

## [0.2.0] - 2020-10-06

### Added

- Add support for setting app version from SHA for test catalogs.

## [0.1.0] - 2020-09-30

### Added

- Add initial version that implements InstallApps for use in apptestctl and
Go integration tests.

[Unreleased]: https://github.com/giantswarm/apptest/compare/v1.3.0...HEAD
[1.3.0]: https://github.com/giantswarm/apptest/compare/v1.2.1...v1.3.0
[1.2.1]: https://github.com/giantswarm/apptest/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/giantswarm/apptest/compare/v1.1.1...v1.2.0
[1.1.1]: https://github.com/giantswarm/apptest/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/giantswarm/giantswarm/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/giantswarm/apptest/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/apptest/compare/v0.12.0...v1.0.0
[0.12.0]: https://github.com/giantswarm/apptest/compare/v0.11.0...v0.12.0
[0.11.0]: https://github.com/giantswarm/apptest/compare/v0.11.0...v0.11.0
[0.11.0]: https://github.com/giantswarm/apptest/compare/v0.10.3...v0.11.0
[0.10.3]: https://github.com/giantswarm/apptest/compare/v0.10.2...v0.10.3
[0.10.2]: https://github.com/giantswarm/apptest/compare/v0.10.1...v0.10.2
[0.10.1]: https://github.com/giantswarm/apptest/compare/v0.10.0...v0.10.1
[0.10.0]: https://github.com/giantswarm/apptest/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/giantswarm/apptest/compare/v0.8.2...v0.9.0
[0.8.2]: https://github.com/giantswarm/apptest/compare/v0.8.1...v0.8.2
[0.8.1]: https://github.com/giantswarm/apptest/compare/v0.8.0...v0.8.1
[0.8.0]: https://github.com/giantswarm/apptest/compare/v0.7.1...v0.8.0
[0.7.1]: https://github.com/giantswarm/apptest/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/giantswarm/apptest/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/giantswarm/apptest/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/giantswarm/apptest/compare/v0.4.1...v0.5.0
[0.4.1]: https://github.com/giantswarm/apptest/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/giantswarm/apptest/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/giantswarm/apptest/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/giantswarm/apptest/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/apptest/releases/tag/v0.1.0
