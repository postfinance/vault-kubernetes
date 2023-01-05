## 0.3.0 (2022-10-18)


### Build System

* **deps**: bump github.com/stretchr/testify from 1.7.0 to 1.7.1 ([de97fc76](https://github.com/postfinance/vault-kubernetes/commit/de97fc76))
  > Bumps [github.com/stretchr/testify](https://github.com/stretchr/testify) from 1.7.0 to 1.7.1.
  > - [Release notes](https://github.com/stretchr/testify/releases)
  > - [Commits](https://github.com/stretchr/testify/compare/v1.7.0...v1.7.1)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/stretchr/testify
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/api from 0.23.4 to 0.23.5 ([fb63e9f2](https://github.com/postfinance/vault-kubernetes/commit/fb63e9f2))
  > Bumps [k8s.io/api](https://github.com/kubernetes/api) from 0.23.4 to 0.23.5.
  > - [Release notes](https://github.com/kubernetes/api/releases)
  > - [Commits](https://github.com/kubernetes/api/compare/v0.23.4...v0.23.5)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/api
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/client-go from 0.23.4 to 0.23.5 ([daf9ad13](https://github.com/postfinance/vault-kubernetes/commit/daf9ad13))
  > Bumps [k8s.io/client-go](https://github.com/kubernetes/client-go) from 0.23.4 to 0.23.5.
  > - [Release notes](https://github.com/kubernetes/client-go/releases)
  > - [Changelog](https://github.com/kubernetes/client-go/blob/master/CHANGELOG.md)
  > - [Commits](https://github.com/kubernetes/client-go/compare/v0.23.4...v0.23.5)
  > 
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/client-go
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: k8s.io/api 0.24.3 -> 0.25.3 ([26827933](https://github.com/postfinance/vault-kubernetes/commit/26827933))
* **deps**: k8s.io/client-go 0.24.3 -> 0.25.3 ([22f9a423](https://github.com/postfinance/vault-kubernetes/commit/22f9a423))


### New Features

* **common**: use distroless image from gcr.io ([583d2704](https://github.com/postfinance/vault-kubernetes/commit/583d2704))



## 0.2.6 (2022-03-14)


### Bug Fixes

* **common**: fixes #24 ([10bae6f8](https://github.com/postfinance/vault-kubernetes/commit/10bae6f8))



## 0.2.5 (2022-02-28)

merge @pszmytka-viacom PR to prevent crashes on non key-value secrets

## 0.2.4 (2022-02-28)

## 0.2.4 (2022-02-28)


## 0.2.3 (2021-10-15)



## 0.2.2 (2021-06-30)



## 0.2.1 (2021-06-10)


### Bug Fixes

* **decode**: fixing decode function for base64 secrets ([421d6017](https://github.com/postfinance/vault-kubernetes/commit/421d6017))



## 0.2.0 (2021-05-28)


### Bug Fixes

* **goreleaser**: enable github releases ([d13ffdfd](https://github.com/postfinance/vault-kubernetes/commit/d13ffdfd))


### New Features

* **common**: use goreleaser to build images ([#19](https://github.com/postfinance/vault-kubernetes/issues/19), [2f6a6548](https://github.com/postfinance/vault-kubernetes/commit/2f6a6548))



## 0.1.7 (2021-01-27)


### New Features

* **synchronizer**: labels for secrets added ([14a1737c](https://github.com/postfinance/vault-kubernetes/commit/14a1737c))



## 0.1.3 (2019-07-18)


### New Features

* **common**: allow to customize the annotation put/searched on managed secrets ([64b32383](https://github.com/postfinance/vault-kubernetes/commit/64b32383))
