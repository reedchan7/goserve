# goserve

A lightweight, high-performance static file server similar to `http-server` and `serve`, but with better performance and no external dependencies. Unlike `http-server` and `serve` which require a Node.js environment, this runs directly as a standalone executable binary.

## Features

- [x] Fast and lightweight static file server
- [x] Cross-platform support (Windows, macOS, Linux)
- [x] URL prefix support
- [x] CORS enabled by default
- [x] Graceful shutdown
- [x] Zero external dependencies for the binary
- [ ] Hot reloading (in progress)

## Installation

### Option 1: Download Binary

Download the pre-compiled binary from the [Releases](https://github.com/reedchan7/goserve/releases) page.

### Option 2: Install with Go

```shell
go install github.com/reedchan7/goserve@main
```

## Usage

```
goserve [options] [directory]
```

### Options

- `-p <port>`: Specify the port to use (default: 3000)
- `-prefix <path>`: Add URL prefix to all requests (default: none)

### Examples

Serve current directory on default port (3000):
```
goserve
```

Serve specific directory on port 8080:
```
goserve -p 8080 ./public
```

Serve with URL prefix (e.g., http://localhost:3000/app/index.html):
```
goserve -prefix /app ./public
```

Combine options:
```
goserve -p 8080 -prefix /api ./data
```
