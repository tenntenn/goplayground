# Go Playground Client

This is a client of [The Go Playground](https://go.dev/play).

## Use as CLI

### Install

```
$ go install github.com/tenntenn/goplayground/cmd/gp@latest
```

### Usage

```
$ gp help
```

### Run

```
$ gp run main.go
```

```
$ gp run a.go b.go
```

```
$ find . -type f | xargs gp run
```

```
$ find . -type f -not -path '*/\.*' | xargs gp run
```

### Format

```
$ gp format [-imports] main.go
```

```
$ gp format [-imports] -output main.go main.go
```

```
$ gp format a.go b.go
```

```
$ find . -type f | xargs gp format
```

```
$ find . -type f -not -path '*/\.*' | xargs gp format
```

### Share

```
$ gp share main.go
```

```
$ gp share a.go b.go
```

```
$ find . -type f | xargs gp share
```

```
$ find . -type f -not -path '*/\.*' | xargs gp share
```

### Download

```
$ gp download https://go.dev/play/p/sTkdodLtokQ
```

```
$ gp dl https://play.golang.org/p/sTkdodLtokQ
```

```
$ gp dl -dldir=output https://go.dev/play/p/sTkdodLtokQ
```

## Version

`version` prints Go version of playground.

```
$ gp version
Version: go1.17.5
Release: go1.17
Name: Go 1.17
```

```
$ gp version -backend gotip
Version: devel go1.18-2c58bb2e42 Wed Jan 5 09:50:29 2022 +0000
Release: go1.18
Name: Go dev branch
```

## With Go dev branch

```
$ gp format -backend gotip example.go2
$ gp run -backend gotip example.go2
$ gp share -backend gotip example.go2
$ gp download -backend gotip hYtdQPeKUC3
```

## Use as a libary

See: https://pkg.go.dev/github.com/tenntenn/goplayground
