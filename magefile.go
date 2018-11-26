// +build mage

package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/postfinance/mage/git"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Test

func init() {
	v, err := strconv.ParseFloat(strings.TrimLeft(runtime.Version(), "go"), 3)
	if err == nil && v > 1.10 {
		fmt.Fprintln(os.Stderr, "go version > 1.10 - unset GOPATH")
		os.Unsetenv("GOPATH")
	}
	os.Setenv("GO111MODULE", "on")
	os.Setenv("CGO_ENABLED", "0") // only stativ builds
	os.Setenv("GOOS", "linux")
}

// Vendor go packages (does wipe first)
func Vendor() error {
	mg.Deps(Clean)
	return sh.Run(mg.GoCmd(), "mod", "vendor")
}

// Clean the workspace (public, vendor)
func Clean() error {
	for _, d := range []string{"dist", "public", "vendor"} {
		if err := sh.Rm(d); err != nil {
			return err
		}
	}
	return nil
}

// TOC for README.md
func TOC() error {
	return sh.Run("doctoc", "--gitlab", "README.md")
}

// Lint run linter
func Lint() error {
	return sh.Run(mg.GoCmd(), "tool", "vet", "cmd", "pkg")
}

// Test run all tests
func Test() error {
	mg.Deps(Lint)
	return sh.Run(mg.GoCmd(), "test", withRace(), "./pkg/...", "./cmd/...")
}

// Coverage run all tests with coverage
func Coverage() error {
	var dir = "public" // target dir for the reports
	var report = fmt.Sprintf("%s/coverage.json", dir)
	var htmlReport = fmt.Sprintf("%s/gocoverage.html", dir)

	mg.Deps(Lint)
	createReport := func() error {
		stdout, err := create(report)
		if err != nil {
			return err
		}
		defer stdout.Close()
		_, err = sh.Exec(nil, stdout, os.Stderr, "gocov", "test", withRace(), "./pkg/...")
		return err
	}
	if err := createReport(); err != nil {
		return err
	}
	convertReport := func() error {
		stdout, err := create(htmlReport)
		if err != nil {
			return err
		}
		defer stdout.Close()
		_, err = sh.Exec(nil, stdout, os.Stderr, "gocov-html", report)
		return err
	}
	return convertReport()
}

var (
	buildSource = "main.go"
)

func getldflags() string {
	return "-d -s -w"
}

// BuildAuth build the auth binary
func BuildAuth() error {
	ldflags := getldflags()
	return sh.Run(mg.GoCmd(), "build", "-ldflags", ldflags, "-o", "dist/authenticator", "cmd/authenticator/main.go")
}

// BuildSync build the sync binary
func BuildSync() error {
	ldflags := getldflags()
	return sh.Run(mg.GoCmd(), "build", "-ldflags", ldflags, "-o", "dist/synchronizer", "cmd/synchronizer/main.go")
}

// BuildRenew build the renew binary
func BuildRenew() error {
	ldflags := getldflags()
	return sh.Run(mg.GoCmd(), "build", "-ldflags", ldflags, "-o", "dist/token-renewer", "cmd/token-renewer/main.go")
}

// BuildAuthImage build vault-kubernetes-authenticator docker image
func BuildAuthImage() error {
	mg.Deps(BuildAuth)
	g, err := git.New(".", git.WithSemverTemplate())
	if err != nil {
		mg.Fatal(1, err)
	}

	image := "vault-kubernetes-authenticator"
	latestImage := fmt.Sprintf("%s:latest", image)
	versionImage := fmt.Sprintf("%s:%s", image, g)

	err = sh.Run("docker", "build", "--build-arg", "BINARY=dist/authenticator", "-t", versionImage, "-f", "packaging/docker/authenticator/Dockerfile", ".")
	if err != nil {
		return err
	}
	err = sh.Run("docker", "tag", versionImage, latestImage)
	if err != nil {
		return err
	}
	return tagAndPush(versionImage, latestImage)
}

// BuildSyncImage build vault-kubernetes-synchronizer docker image
func BuildSyncImage() error {
	mg.Deps(BuildSync)
	g, err := git.New(".", git.WithSemverTemplate())
	if err != nil {
		mg.Fatal(1, err)
	}

	image := "vault-kubernetes-synchronizer"
	latestImage := fmt.Sprintf("%s:latest", image)
	versionImage := fmt.Sprintf("%s:%s", image, g)

	err = sh.Run("docker", "build", "--build-arg", "BINARY=dist/synchronizer", "-t", versionImage, "-f", "packaging/docker/synchronizer/Dockerfile", ".")
	if err != nil {
		return err
	}
	err = sh.Run("docker", "tag", versionImage, latestImage)
	if err != nil {
		return err
	}
	return tagAndPush(versionImage, latestImage)
}

// BuildRenewImage build vault-kubernetes-renew docker image
func BuildRenewImage() error {
	mg.Deps(BuildRenew)
	g, err := git.New(".", git.WithSemverTemplate())
	if err != nil {
		mg.Fatal(1, err)
	}

	image := "vault-kubernetes-token-renewer"
	latestImage := fmt.Sprintf("%s:latest", image)
	versionImage := fmt.Sprintf("%s:%s", image, g)

	err = sh.Run("docker", "build", "--build-arg", "BINARY=dist/token-renewer", "-t", versionImage, "-f", "packaging/docker/token-renewer/Dockerfile", ".")
	if err != nil {
		return err
	}
	err = sh.Run("docker", "tag", versionImage, latestImage)
	if err != nil {
		return err
	}
	return tagAndPush(versionImage, latestImage)
}

// BuildAllImages execute all image build targets
func BuildAllImages() error {
	if err := BuildAuthImage(); err != nil {
		return err
	}
	if err := BuildSyncImage(); err != nil {
		return err
	}
	if err := BuildRenewImage(); err != nil {
		return err
	}
	return nil
}

// tagAndPush tag and push images if DOCKER_TARGET environment variable is set
func tagAndPush(images ...string) error {
	target := os.Getenv("DOCKER_TARGET")
	if len(target) == 0 {
		return nil
	}
	for _, name := range images {
		var err error
		t := path.Join(target, name)
		err = sh.Run("docker", "tag", name, t)
		if err != nil {
			return err
		}
		err = sh.Run("docker", "push", t)
		if err != nil {
			return err
		}
	}
	return nil
}

// withRace sets the race flage if appropriate
func withRace() string {
	if os.Getenv("CGO_ENABLED") == "0" {
		fmt.Fprintf(os.Stderr, "race detector disabled: -race requires cgo; enable cgo by setting CGO_ENABLED=1\n")
		return ""
	}
	return "-race"
}

// create file and all directories if necessary
func create(filename string) (*os.File, error) {
	dir := path.Dir(filename)
	// ensure directories
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	return os.Create(filename)
}
