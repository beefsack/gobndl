package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	BUNDLE_DIR  = ".bndl"
	CONFIG_FILE = "config"
)

func main() {
	flag.Parse()
	command := flag.Arg(0)
	switch command {
	case "init":
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, `Could not get pwd: %s`, err.Error())
			os.Exit(1)
		}
		Init(pwd)
	case "exec":
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, `Could not get pwd: %s`, err.Error())
			os.Exit(1)
		}
		bundlePath, err := FindBundlePath(pwd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not find bundle dir: %s\n",
				err.Error())
			os.Exit(1)
		}
		// Find args
		var args []string
		for i, a := range os.Args {
			if a == "exec" {
				args = os.Args[i+1:]
				break
			}
		}
		Exec(bundlePath, pwd, args[0], args[1:]...)
	case "help":
		Help()
	default:
		Help()
	}
}

func FindBundlePath(from string) (string, error) {
	bundlePath := path.Join(from, BUNDLE_DIR)
	if _, err := os.Stat(bundlePath); err == nil {
		return bundlePath, nil
	}
	if from == "" || path.Base(from) == "/" {
		return "", errors.New("Could not find bundle directory")
	}
	return FindBundlePath(path.Dir(from))
}

func PackageName(bundlePath string) (string, error) {
	file, err := os.Open(path.Join(bundlePath, CONFIG_FILE))
	if err != nil {
		return "", err
	}
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(contents)), nil
}
