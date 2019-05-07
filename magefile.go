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
var Default = BuildAllImages

func init() {
	v, err := strconv.ParseFloat(strings.TrimLeft(runtime.Version(), "go"), 3)
	if err == nil && v > 1.10 {
		fmt.Fprintln(os.Stderr, "go version > 1.10 - unset GOPATH")
		os.Unsetenv("GOPATH")
	}
	os.Setenv("GO111MODULE", "on")
	os.Setenv("CGO_ENABLED", "0") // only static builds
	os.Setenv("GOOS", "linux")
}

// TOC for README.md
func TOC() error {
	return sh.Run("doctoc", "--gitlab", "README.md")
}

func getldflags() string {
	return ""
}

// BuildAuthImage build vault-kubernetes-authenticator docker image
func BuildAuthImage() error {
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
