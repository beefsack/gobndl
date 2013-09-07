package main

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func FindBundlePath(from string) (string, error) {
	bundlePath := path.Join(from, BUNDLE_DIR)
	if _, err := os.Stat(bundlePath); err == nil {
		return bundlePath, nil
	}
	if from == "" || path.Base(from) == "/" {
		return "", errors.New(
			`Run "gobndl init" to initialise a new bundle in this directory`)
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

func MangleVcsDirs(root string) error {
	return filepath.Walk(root,
		func(p string, info os.FileInfo, err error) error {
			for _, vcsP := range VcsPaths {
				if path.Base(p) == vcsP {
					if err := os.Rename(p, fmt.Sprintf("%s%s", p,
						MANGLE_SUFFIX)); err != nil {
						return err
					}
					return filepath.SkipDir
				}
			}
			return nil
		})
}

func UnmangleVcsDirs(root string) error {
	return filepath.Walk(root,
		func(p string, info os.FileInfo, err error) error {
			for _, vcsP := range VcsPaths {
				if path.Base(p) == fmt.Sprintf("%s%s", vcsP, MANGLE_SUFFIX) {
					if err := os.Rename(p,
						p[:len(p)-len(MANGLE_SUFFIX)]); err != nil {
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

func UseBndl(bundlePath string, replacePath bool, cb func() error) error {
	// Set GOPATH
	var newGoPath string
	origGoPath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", origGoPath)
	if replacePath {
		newGoPath = bundlePath
	} else {
		newGoPath = fmt.Sprintf("%s%c%s", bundlePath, os.PathListSeparator,
			origGoPath)
	}
	os.Setenv("GOPATH", newGoPath)
	// Set GOBIN
	origGoBin := os.Getenv("GOBIN")
	defer os.Setenv("GOBIN", origGoBin)
	os.Setenv("GOBIN", path.Join(bundlePath, "bin"))
	// Set PATH
	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)
	os.Setenv("PATH", fmt.Sprintf("%s%c%s", path.Join(bundlePath, "bin"),
		os.PathListSeparator, origPath))
	return cb()
}

func GetImports(packagePath string) ([]string, error) {
	// Find all external imports in the packagePath dir
	packageName, err := GetPackageName(packagePath)
	if err != nil {
		return []string{}, err
	}
	packageMap := map[string]bool{}
	fset := token.NewFileSet()
	packageReg := regexp.MustCompile("^[`\"]?(.+?)[`\"]?$")
	filepath.Walk(packagePath,
		func(p string, info os.FileInfo, err error) error {
			if path.Base(p) == BUNDLE_DIR {
				return filepath.SkipDir
			} else if !info.IsDir() && path.Ext(p) == ".go" {
				f, err := parser.ParseFile(fset, p, nil, parser.ImportsOnly)
				if err == nil {
					for _, s := range f.Imports {
						if matches := packageReg.FindStringSubmatch(
							s.Path.Value); matches != nil &&
							matches[1] != "" && matches[1][0] != '.' {
							importPath := strings.Replace(matches[1], "/",
								string(os.PathSeparator), -1)
							if !strings.Contains(importPath, packageName) {
								packageMap[matches[1]] = true
							}
						}
					}
				}
			}
			return nil
		})
	imports := make([]string, len(packageMap))
	i := 0
	for p, _ := range packageMap {
		imports[i] = p
		i += 1
	}
	return imports, nil
}

func RunCommand(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	// Redirect stdout and stderr to user
	outPipe, err := c.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error getting the output pipe: %s\n",
			err.Error())
	}
	go io.Copy(os.Stdout, outPipe)
	errPipe, err := c.StderrPipe()
	if err != nil {
		return fmt.Errorf("Error getting the error pipe: %s\n", err.Error())
	}
	go io.Copy(os.Stderr, errPipe)
	return c.Run()
}

func GetPackageName(packagePath string) (string, error) {
	for _, gopath := range strings.Split(os.Getenv("GOPATH"),
		string(os.PathListSeparator)) {
		evalGopath, err := filepath.EvalSymlinks(gopath)
		if err != nil {
			return "", err
		}
		if relPath, err := filepath.Rel(path.Join(evalGopath, "src"),
			packagePath); err == nil {
			return relPath, nil
		}
	}
	return "", errors.New("Directory does not appear to be in GOPATH")
}
