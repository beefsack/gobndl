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
    exec        run the following commands using the bundled environment

Common usages:

    Install dependant packages for the current directory into the bundle
        gobndl exec go get
    Install a specific go package into the bundle
        gobndl exec go get github.com/robfig/revel/revel
    Run a specific command from the bundle bin
        gobndl exec revel run github.com/robfig/revel/samples/chat`)
}
