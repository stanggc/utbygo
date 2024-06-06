# utbygo

utbygo is a collection/library of functionalities provided by Go and its ecosystem,
that can be used from C/C++ code. It leverages Go's cgo feature to achieve.

## Features

Currently, utbygo provides the following:

* HTTP client (from Go standard library).
* HTTP server (from Go standard library).
	- Currently implemented as a singleton.

## Building

```shell
$ go build -buildmode=c-archive
```

`utbygo.a` and `utbygo.h` will be generated.

If you want to copy out the build artifacts, make sure to copy
`utbygo.a`, `utbygo.h` and `preamble.h`.

## Usage

See `test_server.c` for an example on usage.

## Usage on MacOS

On MacOS, the CoreFoundation and Security frameworks must be specified
during compilation:

```shell
$ gcc -o testserver test_server.c utbygo.a -framework CoreFoundation -framework Security
```

## Testing

To test the library, build the test server and client:

Server:

```shell
$ gcc -o test_server test_server.c utbygo.a
```

Client:

```shell
$ gcc -o test_client test_client.c utbygo.a
```

Once built, run the test server first, on port 8080 HTTPS only.
Run the test client once server is up and running.

Both should print `OK:`, if tests are OK.
