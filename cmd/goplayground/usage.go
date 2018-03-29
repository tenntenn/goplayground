package main

import "fmt"

func usage() {
	fmt.Println(`goplayground is client of Go Playground.

Usage:

	goplayground command [arguments]

The commands are:
	run		compiles and runs on Go Playground
	format	formats Go code on Go Playground
	share	generates share URL on Go Playground
	help	print this help

Use "go help [command]" for more information about a command.`)
}

func usageRun() {
	fmt.Println(`usage: goplayground run [-json] [gofile]

run compiles and runs on Go Playground.
If [gofile] is not specify, it compiles and runs from stdin.

The flags are:
	-json	output result of run as JSON`)
}

func usageFormat() {
	fmt.Println(`usage: goplayground format [-json] [gofile]

format formats Go code on Go Playground.
If [gofile] is not specify, it compiles and runs from stdin.

The flags are:
	-json	output result of run as JSON`)
}

func usageShare() {
	fmt.Println(`usage: goplayground share [gofile]

share generates share URL on Go Playground.
If [gofile] is not specify, it compiles and runs from stdin.`)
}
