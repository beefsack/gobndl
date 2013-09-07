package main

import (
	"fmt"
	"os"
)

func Exec(bundlePath, pwd, cmd string, args ...string) {
	if err := UnmangleVcsDirs(bundlePath); err != nil {
		fmt.Fprintf(os.Stderr, "Could not unmangle vcs dirs: %s", err.Error())
		os.Exit(1)
	}
	defer MangleVcsDirs(bundlePath)
	if err := UseBndl(bundlePath, false, func() error {
		return RunCommand(cmd, args...)
	}); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
