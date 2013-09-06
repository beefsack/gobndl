package main

import (
	"fmt"
)

func Help() {
	fmt.Println(`gobndl is a tool to bundle package dependencies

Usage:

	gobndl command [arguments]

The commands are:

    init        initialise a new bundle in this directory
    get         get a package or packages into your bundle, if called without
                arguments it will parse all go files in your package directory
                and get all dependencies
    exec        run the following commands using the bundled environment

Common usages:

    Install all dependencies in your package directory into the bundle
        gobndl get
    Install a specific go package into the bundle
        gobndl get github.com/robfig/revel/revel
    Run a specific command from the bundle bin directory
        gobndl exec revel run github.com/robfig/revel/samples/chat`)
}
