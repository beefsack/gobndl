gobndl - Go package dependency bundler
======================================

**gobndl** is a tool to bundle Go dependencies inside your project so you don't
have to worry about external dependencies disappearing, and allowing you to lock
to a specific version of a dependency for your package.

Features
--------

*   Storage of dependencies in local bundle inside package (.bndl directory)
*   Execution of commands in bundled environment using `gobndl exec`
*   Recursively find and install all dependencies using `gobndl get`
*   Ability to check bundle into version control, locking dependency versions

Installation
------------

```bash
~  go get github.com/beefsack/gobndl
```

Usage
-----

```bash
# Initialise your bundle and install dependencies
~  cd github.com/beefsack/my-go-package
~  gobndl init
~  gobndl get

# Use the bundle with exec
~  gobndl exec go build

# You can also install binaries to your bundle
~  gobndl get github.com/robfig/revel/revel
~  gobndl exec revel run github.com/robfig/revel/samples/chat
```
