# No Go

This repo contains the smallest binary I've been able to make using the go
toolchain as an assembler for Plan9 style assembly on non-Plan9 platforms.

It compiles a binary which exits with status code 33.

The package "runtime" provides definitions for everything that "go
tool link" references from Go's runtime while compiling (even though
we explicitly set the entry point something defined in assembly which
doesn't use any of the Go runtime, go tool link still makes assumptions
that it's compiling a Go program, so we need to humour it by creating a
"runtime" with no code in it.)

Other than that, the bulk of the work is the go tool link parameters in
the Makefile.

- "-g" tells is to ignore the checks that it usually does for Go
  packages. This is required because otherwise Go will complain about not
  being a package main. 
- "-L ." tells the linker to search the current directory for packages. 
  This is required in order to make sure it finds our fake runtime, and
  not the one from $GOROOT 
- "-w" disables DWARF generation. This is required because otherwise we'd 
  have to add more symbols to our fake runtime that
  `go tool link` assumes should be present.
- "-E main" tells `go tool link` to use the text segment named "main"
  from our assembly as the program entry point.

The result of this is an 18kb binary on my Linux system (your mileage
may vary depending on Go version and OS), most of which is the type
information we defined in our fake runtime but don't use
 (by contrast, an empty "func main()" defined in Go is 980kb with the
 overhead of the Go runtime, and when I compile similar code on Plan 9
 with 6a/6l the result is ~2kb).
