# Changelog

All notable changes to this project will be documented in this file.

## [0.3.8] - 2025-02-21

### 🐛 Bug Fixes

- Fix tests
- Fix linting

### 💼 Other

- *(deps)* Bump golang.org/x/crypto from 0.24.0 to 0.31.0
- *(deps)* Bump github.com/stretchr/testify from 1.9.0 to 1.10.0
- *(deps)* Bump golang.org/x/net from 0.26.0 to 0.33.0 (#156)
- *(deps)* Bump k8s.io/api from 0.31.2 to 0.32.2
- *(deps)* Bump k8s.io/client-go from 0.31.2 to 0.32.2
- *(deps)* Bump google.golang.org/grpc from 1.43.0 to 1.56.3
- Update all to recent versions

### ⚙️ Miscellaneous Tasks

- *(lint)* Fix linting config (#157)
- Tidy
- Rename default branch to main
- Update to vaultkv 0.0.6
- Update vaultk8s to latest release

## 0.3.5 (2024-04-02)

### Bug Fixes

* **common**: unit tests ([c7369121](https://github.com/postfinance/vault-kubernetes/commit/c7369121))

### Build System

* **deps**: bump github.com/stretchr/testify from 1.8.4 to 1.9.0 ([5bceab69](https://github.com/postfinance/vault-kubernetes/commit/5bceab69))
  > Bumps [github.com/stretchr/testify](https://github.com/stretchr/testify) from 1.8.4 to 1.9.0.
  > - [Release notes](https://github.com/stretchr/testify/releases)
  > - [Commits](https://github.com/stretchr/testify/compare/v1.8.4...v1.9.0)
  >
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/stretchr/testify
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump k8s.io/api from 0.29.0 to 0.29.3 ([c01c6096](https://github.com/postfinance/vault-kubernetes/commit/c01c6096))
  > Bumps [k8s.io/api](https://github.com/kubernetes/api) from 0.29.0 to 0.29.3.
  > - [Commits](https://github.com/kubernetes/api/compare/v0.29.0...v0.29.3)
  >
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/api
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/apimachinery from 0.29.0 to 0.29.3 ([cb7868ba](https://github.com/postfinance/vault-kubernetes/commit/cb7868ba))
  > Bumps [k8s.io/apimachinery](https://github.com/kubernetes/apimachinery) from 0.29.0 to 0.29.3.
  > - [Commits](https://github.com/kubernetes/apimachinery/compare/v0.29.0...v0.29.3)
  >
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/apimachinery
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/client-go from 0.29.0 to 0.29.3 ([7a1f6282](https://github.com/postfinance/vault-kubernetes/commit/7a1f6282))
  > Bumps [k8s.io/client-go](https://github.com/kubernetes/client-go) from 0.29.0 to 0.29.3.
  > - [Changelog](https://github.com/kubernetes/client-go/blob/main/CHANGELOG.md)
  > - [Commits](https://github.com/kubernetes/client-go/compare/v0.29.0...v0.29.3)
  >
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/client-go
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...



## 0.3.4 (2024-01-10)


### Build System

* **deps**: bump golang.org/x/net from 0.8.0 to 0.17.0 ([17304fae](https://github.com/postfinance/vault-kubernetes/commit/17304fae))
  > Bumps [golang.org/x/net](https://github.com/golang/net) from 0.8.0 to 0.17.0.
  > - [Commits](https://github.com/golang/net/compare/v0.8.0...v0.17.0)
  >
  > ---
  > updated-dependencies:
  > - dependency-name: golang.org/x/net
  >   dependency-type: indirect
  > ...



## 0.3.3 (2023-07-19)


### Build System

* **deps**: bump github.com/stretchr/testify from 1.8.2 to 1.8.3 ([56790217](https://github.com/postfinance/vault-kubernetes/commit/56790217))
  > Bumps [github.com/stretchr/testify](https://github.com/stretchr/testify) from 1.8.2 to 1.8.3.
  > - [Release notes](https://github.com/stretchr/testify/releases)
  > - [Commits](https://github.com/stretchr/testify/compare/v1.8.2...v1.8.3)
  >
  > ---
  > updated-dependencies:
  > - dependency-name: github.com/stretchr/testify
  >   dependency-type: direct:production
  >   update-type: version-update:semver-patch
  > ...
* **deps**: bump k8s.io/api from 0.26.2 to 0.27.2 ([9db77990](https://github.com/postfinance/vault-kubernetes/commit/9db77990))
  > Bumps [k8s.io/api](https://github.com/kubernetes/api) from 0.26.2 to 0.27.2.
  > - [Commits](https://github.com/kubernetes/api/compare/v0.26.2...v0.27.2)
  >
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/api
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump k8s.io/apimachinery from 0.26.2 to 0.27.2 ([6f82347a](https://github.com/postfinance/vault-kubernetes/commit/6f82347a))
  > Bumps [k8s.io/apimachinery](https://github.com/kubernetes/apimachinery) from 0.26.2 to 0.27.2.
  > - [Commits](https://github.com/kubernetes/apimachinery/compare/v0.26.2...v0.27.2)
  >
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/apimachinery
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...
* **deps**: bump k8s.io/client-go from 0.26.2 to 0.27.2 ([0bef4de3](https://github.com/postfinance/vault-kubernetes/commit/0bef4de3))
  > Bumps [k8s.io/client-go](https://github.com/kubernetes/client-go) from 0.26.2 to 0.27.2.
  > - [Changelog](https://github.com/kubernetes/client-go/blob/main/CHANGELOG.md)
  > - [Commits](https://github.com/kubernetes/client-go/compare/v0.26.2...v0.27.2)
  >
  > ---
  > updated-dependencies:
  > - dependency-name: k8s.io/client-go
  >   dependency-type: direct:production
  >   update-type: version-update:semver-minor
  > ...



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
  > - [Changelog](https://github.com/kubernetes/client-go/blob/main/CHANGELOG.md)
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
