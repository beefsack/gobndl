package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func Exec(bundlePath string, cmd string, args ...string) {
	// Create go workspace in temporary dir with this package in it
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
	// Set environment variable
	origGoPath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", origGoPath)
	os.Setenv("GOPATH", fmt.Sprintf("%s:%s", bundlePath, tempDir))
	origPath := os.Getenv("GOPATH")
	defer os.Setenv("PATH", origPath)
	os.Setenv("PATH", fmt.Sprintf("%s:%s", path.Join(bundlePath, "bin"),
		os.Getenv("PATH")))
	c := exec.Command(cmd, args...)
	c.Dir = path.Join(makeTo, path.Base(packagePath))
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
