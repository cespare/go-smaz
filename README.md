# go-smaz

go-smaz is a pure Go implementation of [smaz](https://github.com/antirez/smaz), a library for compressing
short strings (particularly containing English words).

## Installation

    $ go get github.com/cespare/go-smaz

## Usage

``` go
import "smaz"
s := "Now is the time for all good men to come to the aid of the party."
compressed := smaz.Compress([]byte(s))           // type is []byte
decompressed, err := smaz.Decompress(compressed) // type is []byte; string(decompressed) == s
```

Also see the [API documentation](go.pkgdoc.org/github.com/cespare/go-smaz).

## Notes

go-smaz is not a direct port of the C version. It is not guaranteed that the output of `smaz.Compress` will be
precisely the same as the C library. However, the output should be decompressible by the C library, and the
output of the C library should be decompressible by `smaz.Decompress`.

Right now go-smaz is very slow -- a very hasty benchmark on my beefy quad-core desktop showed compression
running at < 2MB/s and decompression at ~14MB/s. Initially I just hacked up a working implementation with no
thought for performance; I'll profile it and make it a lot faster when I get a chance.

## Author

Caleb Spare ([cespare](github.com/cespare))

## License

MIT Licensed.

## Other implementations

* [The original C implementation](https://github.com/antirez/smaz)
* [Javascript](https://npmjs.org/package/smaz)
