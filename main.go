package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	BUNDLE_DIR    = ".bndl"
	CONFIG_FILE   = "config"
	MANGLE_SUFFIX = "gobndl"
)

var VcsPaths = []string{
	".git",
	".svn",
	".hg",
	".bzr",
}

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
	case "get":
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, `Could not get pwd: %s`, err.Error())
			os.Exit(1)
		}
		bundlePath, err := FindBundlePath(pwd)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Find args
		var args []string
		for i, a := range os.Args {
			if a == "get" {
				args = os.Args[i+1:]
				break
			}
		}
		Get(bundlePath, args...)
	case "exec":
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, `Could not get pwd: %s`, err.Error())
			os.Exit(1)
		}
		bundlePath, err := FindBundlePath(pwd)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
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
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Please specify a command to execute")
			os.Exit(1)
		}
		Exec(bundlePath, pwd, args[0], args[1:]...)
	case "help":
		Help()
	default:
		Help()
	}
}
