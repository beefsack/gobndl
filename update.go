package main

import (
	"fmt"
	"os"
	"path"
)

func Update(bundlePath string, packages ...string) {
	var err error
	if err := UnmangleVcsDirs(bundlePath); err != nil {
		fmt.Fprintf(os.Stderr, "Could not unmangle vcs dirs: %s", err.Error())
		os.Exit(1)
	}
	defer MangleVcsDirs(bundlePath)
	if len(packages) == 0 {
		packagePath := path.Dir(bundlePath)
		packages, err = GetImports(packagePath)
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}
	}
	if err := UseBndl(bundlePath, true, func() error {
		return RunCommand("go", append([]string{"get", "-u"}, packages...)...)
	}); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
