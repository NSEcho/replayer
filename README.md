# replayer

Simple tool to replay HTTP requests from the file.

Currently supports:
* Count of how many times to replay the request
* Print the response for each request on the stdout
* Automatically gunzip the responses which are gzip compressed

# Installation

```bash
$ git clone https://github.com/lateralusd/replayer.git
$ cd replayer/ && go build
$ ./replayer --help
Usage of ./replayer:
  -c int
        how many times to send request (default 1)
  -p string
        proxy for requests in format https?://ip:port
  -s    print response on the stdout
  -t int
        timeout for request (default 10)
```

# Usage

We will first create the file containing our HTTP request, in this case for `neverssl.com`.

```bash
$ cat req
GET / HTTP/1
Host: neverssl.com
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:99.0) Gecko/20100101 Firefox/99.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
$ ./replayer -s=true ./req
<html>
        <head>
                <title>NeverSSL - Connecting ... </title>
                <style>
                body {
                        font-family: Montserrat, helvetica, arial, sans-serif;
                        font-size: 16x;
                        color: #444444;
...
```