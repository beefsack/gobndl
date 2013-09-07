package main

import (
	"fmt"
	"os"
	"path"
)

func Get(bundlePath string, packages ...string) {
	var err error
	if err := UnmangleVcsDirs(bundlePath); err != nil {
		fmt.Fprintf(os.Stderr, "Could not unmangle vcs dirs: %s", err.Error())
		os.Exit(1)
	}
	defer MangleVcsDirs(bundlePath)
	nonFlagPackageCount := 0
	for _, p := range packages {
		if len(p) > 0 && p[0] != '-' {
			nonFlagPackageCount += 1
		}
	}
	if nonFlagPackageCount == 0 {
		origPackages := packages
		packagePath := path.Dir(bundlePath)
		packages, err = GetImports(packagePath)
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}
		packages = append(origPackages, packages...)
	}
	fmt.Println(packages)
	if err := UseBndl(bundlePath, true, func() error {
		return RunCommand("go", append([]string{"get"}, packages...)...)
	}); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
