package main

import (
	"bufio"
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
		fmt.Fprintln(os.Stderr, "There is already a bundle config in this directory")
	} else {
		fmt.Print("Please enter your package name (eg. github.com/beefsack/gobndl): ")
		br := bufio.NewReader(os.Stdin)
		lineBytes, _, _ := br.ReadLine()
		if err := ioutil.WriteFile(configPath, lineBytes,
			0644); err != nil {
			fmt.Fprintf(os.Stderr, "Could not create %s file: %s\n",
				configPath, err.Error())
			os.Exit(1)
		}
	}
	fmt.Printf(`Created %s
`, BUNDLE_DIR)
}
