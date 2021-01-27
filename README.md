# httkey

httkey is a simple web server that serves static files but with a twist. The host and path of an incoming HTTP request are hashed and the result of this is used to determine the file to serve. httkey was in part inspired by Google's SFFE (Static File Front End) web server.

## Design Specification

* Binaries should be statically-linked with any assets (e.g. error pages) embedded.
* Files should be served from a single directory. There should be no sub-hierarchy.
* The hash should be a 128-bit MurmurHash3 of the request path, seeded with the contents of the `Host` header.
* The binary output of the hash should be encoded in hexadecimal form.
* If a file exists with a name that matches the hex-encoded output of the hash and the request method is `GET` then we should serve it.
* If a file exists but the requests method is not `GET` then we should serve a `405` error page.
* If a file does not exist then we should serve a `404` error page.

It's expected that httkey will be run behind a reverse proxy such as nginx which will take care of the TLS termination. As such, URLs are often handled internally as `http` rather than `https` however this makes no practical difference. Port numbers also have no bearing on the overall hashing.

## Build

To compile httkey you can run the following in a terminal. This will produce 64-bit binaries for Windows and macOS in addition to Linux binaries for the `arm`, `arm64`, `i386` and `amd64` architectures.

```bash
make clean && make build
```

## Quickstart

This quick example demonstrates how you can deploy a static [MTA-STS policy](https://en.wikipedia.org/wiki/MTA-STS) ([RFC 8461](https://tools.ietf.org/html/rfc8461)) using httkey:

```
$ httkey hash http://mta-sts.example.com/.well-known/mta-sts.txt
http://mta-sts.example.com/.well-known/mta-sts.txt	c18eefdae3a2b5d69501ac89b6c5a9e69087e601

$ cat > /tmp/c18eefdae3a2b5d69501ac89b6c5a9e69087e601 <<EOF
version: STSv1
mode: enforce
mx: mail.example.com
max_age: 86400
EOF

$ httkey serve -p /tmp
```

And in a different terminal tab:

```
$ curl --resolve mta-sts.example.com:8080:127.0.0.1 http://mta-sts.example.com:8080/.well-known/mta-sts.txt
version: STSv1
mode: enforce
mx: mail.example.com
max_age: 86400
```

A production MTA-STS deployment would need DNS & TLS provisiong and setting up correctly but this is out of scope for this simple demonstration.

## Usage

```
Usage:
  httkey [command]

Available Commands:
  hash        Hash a URL
  help        Help about any command
  license     Print copyright license
  serve       Run web server
  version     Print version information

Use "httkey [command] --help" for more information about a command.
```

## Copyright

Copyright (c) 2021 Charlie Jonas.\
The code here is released under the [BSD 2-Clause License](https://opensource.org/licenses/BSD-2-Clause).\
See the LICENSE file for full details.
