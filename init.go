package main

import (
	"fmt"
	"os"
	"path"
)

func Init(p string) {
	bundlePath := path.Join(p, BUNDLE_DIR)
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
You should consider ignoring .bndl/bin and .bndl/pkg in version control
`, BUNDLE_DIR)
}
