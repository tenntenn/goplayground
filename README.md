# Go Playground Client

This is a client of [Go Playground](https://play.golang.org).

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

### Download

```
$ gp download https://play.golang.org/p/sTkdodLtokQ
```

```
$ gp dl https://play.golang.org/p/sTkdodLtokQ
```

```
$ gp dl -dldir=output https://play.golang.org/p/sTkdodLtokQ
```

## Try generics codes with go2goplay.golang.org

```
$ gp format -go2 example.go2
$ gp run -go2 example.go2
$ gp share -go2 example.go2
$ gp download -go2 hYtdQPeKUC3
```

## Use as a libary

See: https://pkg.go.dev/github.com/tenntenn/goplayground
