package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tenntenn/goplayground"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "subcomand (run/share/format) should be given")
		os.Exit(1)
	}

	cmdname := strings.Join(os.Args[:1], " ")
	fset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	var asJSON, imports bool

	fset.BoolVar(&asJSON, "json", false, "output as JSON for run or format")
	fset.BoolVar(&imports, "imports", false, "use goimports for format")
	fset.Parse(os.Args[2:])

	switch os.Args[1] {
	case "run":
		if err := run(asJSON, fset.Arg(0)); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "fmt", "format":
		if err := format(asJSON, imports, fset.Arg(0)); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "share":
		if err := share(fset.Arg(0)); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "does not support subcomand", os.Args[1])
		os.Exit(1)
	}
}

func toReader(path string) (io.Reader, func() error, error) {
	if path == "" {
		return os.Stdin, func() error { return nil }, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open given file")
	}

	return f, f.Close, nil
}

func run(asJSON bool, path string) error {
	src, closeFunc, err := toReader(path)
	if err != nil {
		return err
	}
	defer closeFunc()

	var cli goplayground.Client
	r, err := cli.Run(src)
	if err != nil {
		return errors.Wrap(err, "run is failed")
	}

	if asJSON {
		if err := json.NewEncoder(os.Stdout).Encode(r); err != nil {
			return errors.Wrap(err, "result of run cannot encode as JSON")
		}
		return nil
	}

	if r.Errors != "" {
		fmt.Fprintln(os.Stderr, r.Errors)
		return nil
	}

	for i := range r.Events {
		time.Sleep(r.Events[i].Delay)
		switch r.Events[i].Kind {
		case "stdout":
			fmt.Print(r.Events[i].Message)
		case "stderr":
			fmt.Fprint(os.Stderr, r.Events[i].Message)
		}
	}

	return nil
}

func format(asJSON, imports bool, path string) error {
	src, closeFunc, err := toReader(path)
	if err != nil {
		return err
	}
	defer closeFunc()

	var cli goplayground.Client
	r, err := cli.Format(src, imports)
	if err != nil {
		return errors.Wrap(err, "format is failed")
	}

	if asJSON {
		if err := json.NewEncoder(os.Stdout).Encode(r); err != nil {
			return errors.Wrap(err, "result of format cannot encode as JSON")
		}
		return nil
	}

	if r.Error != "" {
		fmt.Fprintln(os.Stderr, r.Error)
	} else {
		fmt.Println(r.Body)
	}

	return nil
}

func share(path string) error {
	src, closeFunc, err := toReader(path)
	if err != nil {
		return err
	}
	defer closeFunc()

	var cli goplayground.Client
	shareURL, err := cli.Share(src)
	if err != nil {
		return errors.Wrap(err, "share is failed")
	}

	fmt.Println(shareURL)

	return nil
}
