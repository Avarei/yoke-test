//go:build mage
// +build mage

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

type Build mg.Namespace
type Push mg.Namespace
type Lint mg.Namespace

const (
	Registry   = "ghcr.io"
	Repository = "avarei/yoke-test"
)

func (Build) Cluster() error {
	fmt.Println("Building...")
	cmd := exec.Command(mg.GoCmd(), "build", "-C", "./cluster", "-o", "../.out/flight-cluster.wasm", "./flight")
	cmd.Env = append(os.Environ(), "GOOS=wasip1", "GOARCH=wasm")
	return cmd.Run()
}

func (Push) Cluster() error {
	fmt.Println("getting appropriate tag:")
	tag, err := getTag("cluster")
	if err != nil {
		return err
	}
	target := fmt.Sprintf("oci://%s/%s/flight-cluster:%s", Registry, Repository, tag)
	cmd := exec.Command("yoke", "push", ".out/flight-cluster.wasm", target)

	fmt.Println("pushing to", target)
	return cmd.Run()
}

func getTag(name string) (string, error) {
	cmd := exec.Command("git", "describe", "--tags", fmt.Sprintf("--match=%s-v*", name), "--always", "--dirty")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}
	longTag := stdout.String()
	tag := strings.TrimPrefix(longTag, fmt.Sprintf("%s-", name))
	return tag, nil
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll(".out")
}
