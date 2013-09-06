package main

import (
	"fmt"
	"os"
)

func Exec(bundlePath, pwd, cmd string, args ...string) {
	if err := UseBndl(bundlePath, false, func() error {
		return RunCommand(cmd, args...)
	}); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	// Clean up the bundle dir in case new things were added
	if err := CleanVcs(bundlePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error cleaning bundle: %s\n", err.Error())
		os.Exit(1)
	}
}
