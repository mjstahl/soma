package cmd

import (
	"fmt"
)

var InfoUsage = `Usage:
    soma info

    Displays information about the Social Machines
    language runtime.  If in a project directory
    the output includes information specific to
    the project.

Example:
    $ soma info
    Social Machines vA.B.C
    Copyright (C) YYYY Mark Stahl

    Portions of this source code are:
    Copyright (C) YYYY The Go Authors
`

func RuntimeInfo(ver string) {
	fmt.Printf("Social Machines v%s\n", ver)
	fmt.Printf("Copyright (C) 2014 Mark Stahl\n")
	fmt.Printf("\n")
	fmt.Printf("Portions of this source code are:\n")
	fmt.Printf("Copyright (C) 2012 The Go Authors\n")
}
