package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/tenntenn/goplayground"
	"golang.org/x/tools/txtar"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "subcomand (run/share/format/version/help) should be given")
		os.Exit(1)
	}

	cmdname := strings.Join(os.Args[:1], " ")
	fset := flag.NewFlagSet(cmdname, flag.ExitOnError)

	fset.Usage = usage

	var (
		go2go, asJSON, imports, open bool
		dldir                        string
		backend                      goplayground.Backend
	)
	fset.BoolVar(&go2go, "go2", false, "Deprecated: use go2goplay.golang.org")
	fset.BoolVar(&asJSON, "json", false, "output as JSON for run or format")
	fset.BoolVar(&imports, "imports", false, "use goimports for format")
	fset.BoolVar(&open, "open", false, "open url in browser for share")
	fset.StringVar(&dldir, "dldir", "", "output directory for download")
	fset.Var(&backend, "backend", `go version: empty is release version and "gotip" is the developer branch`)
	fset.Parse(os.Args[2:])

	p := &playground{
		cli: &goplayground.Client{
			Backend: backend,
		},
		asJSON:  asJSON,
		imports: imports,
		open:    open,
		dldir:   dldir,
	}

	if go2go {
		fmt.Fprintln(os.Stderr, "The option -go2 is deprecated.")
		fmt.Fprintln(os.Stderr, "Please use -v=gotip instead of it.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		if err := p.run(fset.Args()...); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "fmt", "format":
		if err := p.format(fset.Args()...); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "share":
		if err := p.share(fset.Args()...); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "dl", "download":
		if err := p.download(fset.Args()); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "version":
		if err := p.version(); err != nil {
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
			return nil, fmt.Errorf("cannot read file (%s): %w", paths[0], err)
		}
		return bytes.NewReader(data), nil
	}

	var a txtar.Archive
	for _, p := range paths {
		data, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("cannot read file (%s): %w", p, err)
		}
		a.Files = append(a.Files, txtar.File{
			Name: filepath.ToSlash(filepath.Clean(p)),
			Data: data,
		})
	}

	return bytes.NewReader(txtar.Format(&a)), nil
}

type playground struct {
	cli     *goplayground.Client
	asJSON  bool
	imports bool
	open    bool
	dldir   string
}

func (p *playground) run(paths ...string) error {
	src, err := toReader(paths...)
	if err != nil {
		return err
	}

	r, err := p.cli.Run(src)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	if p.asJSON {
		if err := json.NewEncoder(os.Stdout).Encode(r); err != nil {
			return fmt.Errorf("result of run cannot encode as JSON: %w", err)
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

func (p *playground) format(paths ...string) error {
	src, err := toReader(paths...)
	if err != nil {
		return err
	}

	r, err := p.cli.Format(src, p.imports)
	if err != nil {
		return fmt.Errorf("format: %w", err)
	}

	if p.asJSON {
		if err := json.NewEncoder(os.Stdout).Encode(r); err != nil {
			return fmt.Errorf("result of format cannot encode as JSON: %w", err)
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

func (p *playground) share(paths ...string) error {

	src, err := toReader(paths...)
	if err != nil {
		return err
	}

	shareURL, err := p.cli.Share(src)
	if err != nil {
		return fmt.Errorf("share: %w", err)
	}

	if p.cli.Backend != goplayground.BackendDefault {
		params := shareURL.Query()
		params.Set("v", p.cli.Backend.String())
		shareURL.RawQuery = params.Encode()
	}

	if p.open {
		if err = browser.OpenURL(shareURL.String()); err != nil {
			return err
		}
	}

	fmt.Println(shareURL)

	return nil
}

func toHashOrURL(r io.Reader) (string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("cannot read hash or URL: %w", err)
	}
	return strings.TrimSpace(string(b)), nil
}

func (p *playground) download(args []string) error {
	var hashOrURL string
	if len(args) <= 0 {
		s, err := toHashOrURL(os.Stdin)
		if err != nil {
			return fmt.Errorf("download: %w", err)
		}
		hashOrURL = s
	} else {
		hashOrURL = args[0]
	}

	var buf bytes.Buffer
	if err := p.cli.Download(&buf, hashOrURL); err != nil {
		return fmt.Errorf("download: %w", err)
	}

	if p.dldir == "" {
		if _, err := io.Copy(os.Stdout, &buf); err != nil {
			return fmt.Errorf("download: %w", err)
		}
		return nil
	}

	data := buf.Bytes()
	a := txtar.Parse(data)
	if len(a.Files) == 0 {
		fname := hashOrURL
		dlURL, err := url.Parse(fname)
		if err == nil { // hashOrURL is URL
			fname = path.Base(dlURL.Path)
			if !strings.HasSuffix(fname, ".go") {
				fname += ".go"
			}
		}

		f, err := os.Create(filepath.Join(p.dldir, fname))
		if err != nil {
			return fmt.Errorf("download: %w", err)
		}

		if _, err := io.Copy(f, bytes.NewReader(data)); err != nil {
			return fmt.Errorf("download: %w", err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("download: %w", err)
		}

		return nil
	}

	if v := bytes.TrimSpace(a.Comment); len(v) > 0 {
		a.Files = append([]txtar.File{txtar.File{
			Name: "prog.go",
			Data: a.Comment,
		}}, a.Files...)
		a.Comment = nil
	}

	for _, f := range a.Files {
		fpath := filepath.Join(p.dldir, filepath.FromSlash(f.Name))
		fmt.Printf("output %s ... ", fpath)

		if err := os.MkdirAll(filepath.Dir(fpath), 0o777); err != nil {
			return fmt.Errorf("download: %w", err)
		}

		dst, err := os.Create(fpath)
		if err != nil {
			return fmt.Errorf("download: %w", err)
		}

		if _, err := io.Copy(dst, bytes.NewReader(f.Data)); err != nil {
			return fmt.Errorf("download: %w", err)
		}

		if err := dst.Close(); err != nil {
			return fmt.Errorf("download: %w", err)
		}

		fmt.Println("ok")
	}

	return nil
}

func (p *playground) version() error {
	r, err := p.cli.Version()
	if err != nil {
		return fmt.Errorf("version: %w", err)
	}

	fmt.Println("Version:", r.Version)
	fmt.Println("Release:", r.Release)
	fmt.Println("Name:", r.Name)

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
