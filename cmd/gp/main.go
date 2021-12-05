package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/browser"
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

	var (
		go2go, asJSON, imports, open bool
		dldir, outputpath            string
	)
	fset.BoolVar(&go2go, "go2", false, "use "+goplayground.Go2BaseURL)
	fset.BoolVar(&asJSON, "json", false, "output as JSON for run or format")
	fset.BoolVar(&imports, "imports", false, "use goimports for format")
	fset.BoolVar(&open, "open", false, "open url in browser for share")
	fset.StringVar(&dldir, "dldir", "", "output directory for download")
	fset.StringVar(&outputpath, "output", "", "output file path for format")
	fset.Parse(os.Args[2:])

	p := &playground{
		asJSON: asJSON,
		imports: imports,
		open: open,
		path: outputpath,
	}

	if go2go {
		p.cli.BaseURL = goplayground.Go2BaseURL
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
		var hashOrURL string
		if fset.NArg() <= 0 {
			s, err := toHashOrURL(os.Stdin)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
			hashOrURL = s
		} else {
			hashOrURL = fset.Arg(0)
		}

		var buf bytes.Buffer
		if err := p.download(&buf, hashOrURL); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		if dldir == "" {
			if _, err := io.Copy(os.Stdout, &buf); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
			return
		}

		data := buf.Bytes()
		a := txtar.Parse(data)
		if len(a.Files) == 0 {
			fname := path.Base(hashOrURL)
			fname = fname[:len(fname)-len(filepath.Ext(fname))] + ".go"
			f, err := os.Create(filepath.Join(dldir, fname))
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}

			if _, err := io.Copy(f, bytes.NewReader(data)); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
		}

		if err := txtar.Write(a, dldir); err != nil {
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
	for _, p := range paths {
		data, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot read file", p)
		}
		a.Files = append(a.Files, txtar.File{
			Name: filepath.Clean(p),
			Data: data,
		})
	}

	return bytes.NewReader(txtar.Format(&a)), nil
}

type playground struct {
	cli     goplayground.Client
	asJSON  bool
	imports bool
	open    bool
	path    string
}

func (p *playground) run(paths ...string) error {
	src, err := toReader(paths...)
	if err != nil {
		return err
	}

	r, err := p.cli.Run(src)
	if err != nil {
		return errors.Wrap(err, "run is failed")
	}

	if p.asJSON {
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

func (p *playground) format(paths ...string) error {
	src, err := toReader(paths...)
	if err != nil {
		return err
	}

	r, err := p.cli.Format(src, p.imports)
	if err != nil {
		return errors.Wrap(err, "format is failed")
	}

	if p.asJSON {
		if err := json.NewEncoder(os.Stdout).Encode(r); err != nil {
			return errors.Wrap(err, "result of format cannot encode as JSON")
		}
		return nil
	}

	if r.Error != "" {
		fmt.Fprintln(os.Stderr, r.Error)
		return nil
	}

	if p.path != "" {
		if err := ioutil.WriteFile(p.path, []byte(r.Body), os.ModePerm); err != nil {
			return errors.Wrap(err, "failed to write file")
		}
	}
	fmt.Println(r.Body)
	return nil
}

func (p *playground) share(paths ...string) error {

	src, err := toReader(paths...)
	if err != nil {
		return err
	}

	shareURL, err := p.cli.Share(src)
	if err != nil {
		return errors.Wrap(err, "share is failed")
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
		return "", errors.Wrap(err, "cannot read hash or URL")
	}
	return strings.TrimSpace(string(b)), nil
}

func (p *playground) download(w io.Writer, hashOrURL string) error {
	if err := p.cli.Download(w, hashOrURL); err != nil {
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
