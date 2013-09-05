package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func Exec(bundlePath, pwd, cmd string, args ...string) {
	// Figure out the package name, path, and relative path to package root
	packageName, err := PackageName(bundlePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting the package name: %s\n",
			err.Error())
		os.Exit(1)
	}
	if packageName == "" {
		fmt.Fprintln(os.Stderr, "Could not find package name in config")
		os.Exit(1)
	}
	packagePath := path.Dir(bundlePath)
	relativeFromPackageRoot, err := filepath.Rel(packagePath, pwd)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Error getting relative path to package root: %s\n", err.Error())
		os.Exit(1)
	}
	// Create go workspace in temporary dir with this package in it
	tempDir, err := ioutil.TempDir("", "gobndl")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get temp dir: %s\n", err.Error())
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)
	makeTo := path.Join(tempDir, "src", path.Dir(packageName))
	if err := os.MkdirAll(makeTo, 0755); err != nil {
		fmt.Fprintf(os.Stderr,
			"Could not make temporary workspace for package: %s\n", err.Error())
		os.Exit(1)
	}
	cpCmd := exec.Command("cp", "-rf", packagePath, makeTo)
	if err := cpCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr,
			"Could not copy package to temporary workspace: %s\n", err.Error())
		os.Exit(1)
	}
	// Set environment variables
	origGoPath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", origGoPath)
	os.Setenv("GOPATH", fmt.Sprintf("%s%u%s", bundlePath, os.PathListSeparator,
		tempDir))
	origGoBin := os.Getenv("GOBIN")
	defer os.Setenv("GOBIN", origGoBin)
	os.Setenv("GOBIN", path.Join(bundlePath, "bin"))
	origPath := os.Getenv("GOPATH")
	defer os.Setenv("PATH", origPath)
	os.Setenv("PATH", fmt.Sprintf("%s%u%s", path.Join(bundlePath, "bin"),
		os.PathListSeparator, os.Getenv("PATH")))
	c := exec.Command(cmd, args...)
	c.Dir = path.Join(makeTo, path.Base(packagePath), relativeFromPackageRoot)
	// Redirect stdout and stderr to user
	outPipe, err := c.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting the output pipe: %s\n",
			err.Error())
	}
	go io.Copy(os.Stdout, outPipe)
	errPipe, err := c.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting the error pipe: %s\n",
			err.Error())
	}
	go io.Copy(os.Stderr, errPipe)
	// Run command
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %s\n", err.Error())
		os.Exit(1)
	}
}
