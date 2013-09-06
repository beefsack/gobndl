package main

import (
	"fmt"
	"os"
	"path"
)

func Init(p string) {
	bundlePath := path.Join(p, BUNDLE_DIR)
	// Try to create the directory
	if _, err := os.Stat(bundlePath); err == nil {
		fmt.Fprintln(os.Stderr, "There is already a bundle in this directory")
	} else {
		if err := os.Mkdir(bundlePath, 0775); err != nil {
			fmt.Fprintf(os.Stderr, "Could not create %s directory: %s\n",
				bundlePath, err.Error())
			os.Exit(1)
		}
	}
	fmt.Printf(`Created %s
`, BUNDLE_DIR)
}
