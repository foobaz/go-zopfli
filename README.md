go-zopfli
=========

Go port of Zopfli, a zlib-compatible compression library.

Zopfli compresses data more effectively than zlib does, at the expense of
compression speed. The go-zopfli port is 2-3 times slower than the C version,
and compresses at approximately 100 kB/s.

It can be used to compress files that will not change often or if you need a
lower-level interface to Deflate compression.

[See the package documentation.](http://godoc.org/github.com/foobaz/go-zopfli/zopfli)
