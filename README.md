HTTP Here
=========

Very simple golang webserver to serve the current directory. For when you just
want a web server, here, now.

    $ httphere
    2014/05/06 23:17:47 Opening server on :8080 available at:
      - http://127.0.0.1:8080/
      - http://192.0.2.36:8080/
      - http://[::1]:8080/
      - http://[2001:db8:11a9:c:cafe:4eee:eeee:eee]:8080/
      - http://[fe80::5e51:4eee:eeee:eee]:8080/

If you trust me (and I have a build for your system):

    curl -O https://dgl.cx/httphere/$(uname -sm | tr 'A-Z ' 'a-z-')/httphere
    chmod +x httphere
    ./httphere

This will serve the current directory, doing this in your home directory may be
a bad idea. Better to store it in your ~/bin, then run it in a specific
directory.

Some of these builds may not work for various reasons (e.g. go statically links
on OpenBSD, this will only work for a specific version because OpenBSD does not
guarantee compatibility between versions).


Why?
====

I was using the simple:

    python3 -m http.server

However this is single threaded and will block if you have a download in
progress. Also it didn't support IPv6 by default. Both of which are a bit lame.
So this is a very simple go program.

Building
========

Install go, if you haven't, see http://golang.org/doc/install

Then:

    go build -ldflags -s httphere.go

(Using -ldflags -s cuts the executable size in half.)
