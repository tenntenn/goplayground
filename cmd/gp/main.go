package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/txtar"
	"github.com/tenntenn/goplayground"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "subcomand (run/share/format/help) should be given")
		os.Exit(1)
	}

	cmdname := strings.Join(os.Args[:1], " ")
	fset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	fset.Usage = usage

	var asJSON, imports bool
	fset.BoolVar(&asJSON, "json", false, "output as JSON for run or format")
	fset.BoolVar(&imports, "imports", false, "use goimports for format")
	fset.Parse(os.Args[2:])

	switch os.Args[1] {
	case "run":
		if err := run(asJSON, fset.Args()...); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "fmt", "format":
		if err := format(asJSON, imports, fset.Args()...); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "share":
		if err := share(fset.Args()...); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "dl", "download":
		var r io.Reader
		if fset.NArg() <= 0 {
			r = os.Stdin
		} else {
			r = strings.NewReader(fset.Arg(0))
		}
		if err := download(r); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "-h", "help":
		help(fset.Arg(0))
	default:
		fmt.Fprintln(os.Stderr, "does not support subcomand", os.Args[1])
		fset.Usage()
		os.Exit(1)
	}
}

func toReader(paths ...string) (io.Reader, error) {
	if len(paths) == 0 {
		return os.Stdin, nil
	}

	if len(paths) == 1 {
		data, err := ioutil.ReadFile(paths[0])
		if err != nil {
			return nil, errors.Wrapf(err, "cannot read file", paths[0])
		}
		return bytes.NewReader(data), nil
	}

	var a txtar.Archive
	for _, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot read file", path)
		}
		a.Files = append(a.Files, txtar.File{
			Name: filepath.Clean(path),
			Data: data,
		})
	}

	return bytes.NewReader(txtar.Format(&a)), nil
}

func run(asJSON bool, paths ...string) error {
	src, err := toReader(paths...)
	if err != nil {
		return err
	}

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

func format(asJSON, imports bool, paths ...string) error {
	src, err := toReader(paths...)
	if err != nil {
		return err
	}

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

func share(paths ...string) error {

	src, err := toReader(paths...)
	if err != nil {
		return err
	}

	var cli goplayground.Client
	shareURL, err := cli.Share(src)
	if err != nil {
		return errors.Wrap(err, "share is failed")
	}

	fmt.Println(shareURL)

	return nil
}

func download(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "download is failed")
	}

	hashOrURL := strings.TrimSpace(string(b))
	var cli goplayground.Client
	if err := cli.Download(os.Stdout, hashOrURL); err != nil {
		return errors.Wrap(err, "download is failed")
	}
	return nil
}

func help(cmd string) {
	switch cmd {
	case "run":
		usageRun()
	case "format":
		usageFormat()
	case "share":
		usageShare()
	case "download":
		usageDownload()
	default:
		usage()
	}
}
