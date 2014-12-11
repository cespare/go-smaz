# smaz

This is a pure Go implementation of [antirez's](https://github.com/antirez)
[smaz](https://github.com/antirez/smaz), a library for compressing short strings
(particularly containing English words).

## Installation

    $ go get github.com/kjk/smaz

## Usage

``` go
import (
  "github.com/kjk/smaz"
)

func main() {
  s := "Now is the time for all good men to come to the aid of the party."
  compressed := smaz.Encode(nil, []byte(s))
  decompressed, err := smaz.Decode(compressed)
  if err != nil {
    fmt.Printf("decompressed: %s\n", string(decompressed))
    ...
}
```

Also see the [API documentation](http://godoc.org/github.com/kjk/smaz).

## Notes

smaz is not a direct port of the C version. It is not guaranteed that the output
of `smaz.Compress` will be precisely the same as the C library. However, the
output should be decompressible by the C library, and the output of the C
library should be decompressible by `smaz.Decompress`.

## Author

[Salvatore Sanfilippo](https://github.com/antirez) designed smaz and wrote
[C implementation]](https://github.com/antirez/smaz).

[Caleb Sparece](https://github.com/cespare) wrote initial
[Go port](https://github.com/cespare/go-smaz).

[Krzysztof Kowalczyk](http://blog.kowalczyk.info) improved speed of
decompression (2x faster) and compression.

## Contributors

[Antoine Grondin](https://github.com/aybabtme)

## License

MIT Licensed.

## Other implementations

* [The original C implementation](https://github.com/antirez/smaz)
* [Javascript](https://npmjs.org/package/smaz)
