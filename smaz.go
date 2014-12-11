// Package smaz is an implementation of the smaz library (https://github.com/antirez/smaz) for compressing
// small strings.
package smaz

import (
	"errors"

	"github.com/kjk/smaz/trie"
)

var codeStrings = []string{" ",
	"the", "e", "t", "a", "of", "o", "and", "i", "n", "s", "e ", "r", " th",
	" t", "in", "he", "th", "h", "he ", "to", "\r\n", "l", "s ", "d", " a", "an",
	"er", "c", " o", "d ", "on", " of", "re", "of ", "t ", ", ", "is", "u", "at",
	"   ", "n ", "or", "which", "f", "m", "as", "it", "that", "\n", "was", "en",
	"  ", " w", "es", " an", " i", "\r", "f ", "g", "p", "nd", " s", "nd ", "ed ",
	"w", "ed", "http://", "for", "te", "ing", "y ", "The", " c", "ti", "r ", "his",
	"st", " in", "ar", "nt", ",", " to", "y", "ng", " h", "with", "le", "al", "to ",
	"b", "ou", "be", "were", " b", "se", "o ", "ent", "ha", "ng ", "their", "\"",
	"hi", "from", " f", "in ", "de", "ion", "me", "v", ".", "ve", "all", "re ",
	"ri", "ro", "is ", "co", "f t", "are", "ea", ". ", "her", " m", "er ", " p",
	"es ", "by", "they", "di", "ra", "ic", "not", "s, ", "d t", "at ", "ce", "la",
	"h ", "ne", "as ", "tio", "on ", "n t", "io", "we", " a ", "om", ", a", "s o",
	"ur", "li", "ll", "ch", "had", "this", "e t", "g ", "e\r\n", " wh", "ere",
	" co", "e o", "a ", "us", " d", "ss", "\n\r\n", "\r\n\r", "=\"", " be", " e",
	"s a", "ma", "one", "t t", "or ", "but", "el", "so", "l ", "e s", "s,", "no",
	"ter", " wa", "iv", "ho", "e a", " r", "hat", "s t", "ns", "ch ", "wh", "tr",
	"ut", "/", "have", "ly ", "ta", " ha", " on", "tha", "-", " l", "ati", "en ",
	"pe", " re", "there", "ass", "si", " fo", "wa", "ec", "our", "who", "its", "z",
	"fo", "rs", ">", "ot", "un", "<", "im", "th ", "nc", "ate", "><", "ver", "ad",
	" we", "ly", "ee", " n", "id", " cl", "ac", "il", "</", "rt", " wi", "div",
	"e, ", " it", "whi", " ma", "ge", "x", "e c", "men", ".com",
}

var codes = make([][]byte, len(codeStrings))
var codeTrie = trie.New()

func init() {
	for i, code := range codeStrings {
		codes[i] = []byte(code)
		codeTrie.Put([]byte(code), i)
	}
}

func next(d []byte, n int) ([]byte, []byte) {
	if len(d) < n {
		return d, nil
	}
	return d[:n], d[n:]
}

func flushVerb(outBufPtr, verbBufPtr *[]byte) {
	// We can write a max of 255 continuous verbatim characters, because the
	// length of the continous verbatim section is represented by a single byte.
	outBuf := *outBufPtr
	verbBuf := *verbBufPtr
	var chunk []byte
	for len(verbBuf) > 0 {
		chunk, verbBuf = next(verbBuf, 255)
		if len(chunk) == 1 {
			// 254 is code for a single verbatim byte
			outBuf = append(outBuf, byte(254))
		} else {
			// 255 is code for a verbatim string. It is followed by a byte
			// containing the length of the string.
			outBuf = append(outBuf, byte(255))
			outBuf = append(outBuf, byte(len(chunk)))
		}
		outBuf = append(outBuf, chunk...)
	}
	*outBufPtr = outBuf
	*verbBufPtr = verbBuf[:0]
}

// Compress compresses a byte slice and returns the compressed data.
func Compress(input []byte) []byte {
	var dst []byte
	var verbBuf []byte
	root := codeTrie.Root()

	for len(input) > 0 {
		prefixLen := 0
		code := 0
		node := root
		for i, c := range input {
			next, ok := node.Walk(c)
			if !ok {
				break
			}
			node = next
			if node.Terminal() {
				prefixLen = i + 1
				code = node.Val()
			}
		}

		if prefixLen > 0 {
			input = input[prefixLen:]
			flushVerb(&dst, &verbBuf)
			dst = append(dst, byte(code))
		} else {
			verbBuf = append(verbBuf, input[0])
			input = input[1:]
		}
	}
	flushVerb(&dst, &verbBuf)

	return dst
}

// ErrDecompression is returned when decompressing invalid smaz-encoded data.
var ErrDecompression = errors.New("Invalid or corrupted compressed data.")

// Decompress decompresses a smaz-compressed byte slice and return a new slice with the decompressed data. err
// is nil if and only if decompression fails for any reason (e.g., corrupted data).
func Decompress(compressed []byte) ([]byte, error) {
	decompressed := make([]byte, 0, len(compressed))
	for len(compressed) > 0 {
		switch compressed[0] {
		case 254: // Verbatim byte
			if len(compressed) < 2 {
				return nil, ErrDecompression
			}
			decompressed = append(decompressed, compressed[1])
			compressed = compressed[2:]
		case 255: // Verbatim string
			if len(compressed) < 2 {
				return nil, ErrDecompression
			}
			n := int(compressed[1])
			if len(compressed) < n+2 {
				return nil, ErrDecompression
			}
			decompressed = append(decompressed, compressed[2:n+2]...)
			compressed = compressed[n+2:]
		default: // Look up encoded value
			d := codes[int(compressed[0])]
			decompressed = append(decompressed, d...)
			compressed = compressed[1:]
		}
	}

	return decompressed, nil
}
