package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	BUNDLE_DIR  = ".bndl"
	CONFIG_FILE = "config"
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

func CleanVcs(root string) error {
	return filepath.Walk(root,
		func(p string, info os.FileInfo, err error) error {
			for _, vcsP := range VcsPaths {
				if path.Base(p) == vcsP {
					err := os.RemoveAll(p)
					if err != nil {
						return err
					}
					return filepath.SkipDir
				}
			}
			return nil
		})
}

func CopyGoDir(src, dest string) error {
	fi, err := os.Stat(dest)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return errors.New("Destination is not a directory")
	}
	destDir := path.Join(dest, path.Base(src))
	if err := os.Mkdir(destDir, 0755); err != nil {
		return err
	}
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, f := range files {
		fPath := path.Join(src, f.Name())
		if f.Name() == BUNDLE_DIR {
			continue
		}
		if f.IsDir() {
			// Recurse dir
			if err := CopyGoDir(fPath, destDir); err != nil {
				return err
			}
		} else {
			// Copy file
			sf, err := os.Open(fPath)
			if err != nil {
				return err
			}
			defer sf.Close()
			df, err := os.Create(path.Join(destDir, f.Name()))
			if err != nil {
				return err
			}
			defer df.Close()
			if _, err := io.Copy(df, sf); err != nil {
				return err
			}
		}
	}
	return nil
}
