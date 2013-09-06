package main

import (
	"fmt"
	"os"
	"path"
)

func Get(bundlePath string, packages ...string) {
	if len(packages) == 0 {
		packagePath := path.Dir(bundlePath)
		packages = GetImports(packagePath)
	}
	if err := UseBndl(bundlePath, true, func() error {
		return RunCommand("go", append([]string{"get"}, packages...)...)
	}); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	if err := CleanVcs(bundlePath); err != nil {
		fmt.Fprintf(os.Stderr, "Could not clean up vcs dirs: %s", err.Error())
		os.Exit(1)
	}
}
