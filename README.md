gobndl - Go package bundling
============================

**gobndl** is a tool to bundle Go packages inside your project so you don't have
to worry about external dependencies disappearing, and allows you to lock to a
specific version of a dependency.

Features
--------

*   Storage of dependencies in local bundle inside package (.bndl directory)
*   Execution of commands in bundled environment using `gobndl exec`
*   Recursively find and install all dependencies using `gobndl get`
*   Ability to check bundle into version control, locking dependency versions

Installation
------------

```
-> go get github.com/beefsack/gobndl
-> gobndl help
```

Usage
-----

*   Initialise a bundle in your project root directory by running `gobndl init`
*   Get all package dependencies in your project using `gobndl get`
*   Get a specific package using `gobndl get github.com/robfig/revel/revel`
*   Execute commands using your bundle environment and binaries in your bundle
    using `gobndl exec revel run github.com/robfig/revel/samples/chat`