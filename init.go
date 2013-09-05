package main

import (
	"fmt"
	"io/ioutil"
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
	// Make sure the config exists
	configPath := path.Join(bundlePath, CONFIG_FILE)
	if _, err := os.Stat(configPath); err == nil {
		fmt.Fprintln(os.Stderr, "There is already a bundle in this directory")
	} else {
		if err := ioutil.WriteFile(configPath, []byte("your/package/name/here"),
			0644); err != nil {
			fmt.Fprintf(os.Stderr, "Could not create %s file: %s\n",
				configPath, err.Error())
			os.Exit(1)
		}
	}
	fmt.Printf(`Created %s
Make sure to edit %s and replace the package name with your own
`, BUNDLE_DIR, configPath)
}
