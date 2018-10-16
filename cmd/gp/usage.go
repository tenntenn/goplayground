package main

import "fmt"

func usage() {
	fmt.Println(`gp is client of Go Playground.

Usage:

	gp command [arguments]

The commands are:
	run			compiles and runs on Go Playground
	format		formats Go code on Go Playground
	share		generates share URL on Go Playground
	download	download given hash or URL Go code
	help		print this help

Use "go help [command]" for more information about a command.`)
}

func usageRun() {
	fmt.Println(`usage: gp run [-json] [gofile]

"run" compiles and runs on Go Playground.
If [gofile] is not specify, it compiles and runs from stdin.

The flags are:
	-json	output result of run as JSON`)
}

func usageFormat() {
	fmt.Println(`usage: gp format [-json] [gofile]

"format" formats Go code on Go Playground.
If [gofile] is not specify, it compiles and runs from stdin.

The flags are:
	-json	output result of run as JSON`)
}

func usageShare() {
	fmt.Println(`usage: goplayground share [gofile]

"share" generates share URL on Go Playground.
If [gofile] is not specify, it compiles and runs from stdin.`)
}

func usageDownload() {
	fmt.Println(`usage: gp download [hash|share URL]

"download" downloads Go code corresponds to given hash or URL.
If hash or share URL is not specify, it compiles and runs from stdin.`)
}
