// Package smaz is an implementation of the smaz library (https://github.com/antirez/smaz) for compressing
// small strings.
package smaz

import (
	"bytes"
	"errors"
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
var prefixToCode = make(map[string]byte)
var maxCodeLen = 0 // TODO: Unnecessary when we switch to a trie implementation

func init() {
	for i, code := range codeStrings {
		codes[i] = []byte(code)
		prefixToCode[code] = byte(i)
		if len(code) > maxCodeLen {
			maxCodeLen = len(code)
		}
	}
}

// BUG(cespare): Compress is written in an extremely naive manner for the time being and is very slow. I will
// reimplement (and then profile/optimize it) after I get done with go-trie, which will be a natural fit for
// this problem.

// Compress compresses a byte slice and returns the compressed data.
func Compress(input []byte) []byte {
	var outBuf bytes.Buffer
	var verbatim bytes.Buffer

	flushVerbatim := func() {
		// We can write a max of 255 continuous verbatim characters, because the length of the continous verbatim
		// section is represented by a single byte.
		for verbatim.Len() > 0 {
			chunk := verbatim.Next(255)
			if len(chunk) == 1 {
				// 254 is code for a single verbatim byte
				outBuf.WriteByte(byte(254))
			} else {
				// 255 is code for a verbatim string. It is followed by a byte containing the length of the string.
				outBuf.WriteByte(byte(255))
				outBuf.WriteByte(byte(len(chunk)))
			}
			outBuf.Write(chunk)
		}
		verbatim.Reset()
	}

	for len(input) > 0 {
		// Find the longest matching substring, brute-force
		maxPossibleMatch := maxCodeLen
		if len(input) < maxPossibleMatch {
			maxPossibleMatch = len(input)
		}
		matchFound := false
		for matchLen := maxPossibleMatch; matchLen > 0; matchLen-- {
			prefix := input[:matchLen]
			if code, ok := prefixToCode[string(prefix)]; ok {
				// Match found
				input = input[matchLen:]
				flushVerbatim()
				outBuf.WriteByte(code)
				matchFound = true
				break
			}
		}
		if !matchFound {
			verbatim.WriteByte(input[0])
			input = input[1:]
		}
	}
	flushVerbatim()

	return outBuf.Bytes()
}

// DecompressionError is returned when decompressing invalid smaz-encoded data.
var DecompressionError = errors.New("Invalid or corrupted compressed data.")

// Decompress decompresses a smaz-compressed byte slice and return a new slice with the decompressed data. err
// is nil if and only if decompression fails for any reason (e.g., corrupted data).
func Decompress(compressed []byte) ([]byte, error) {
	decompressed := bytes.NewBuffer(make([]byte, 0, len(compressed))) // Estimate initial size

	for len(compressed) > 0 {
		switch compressed[0] {
		case 254: // Verbatim byte
			if len(compressed) < 2 {
				return nil, DecompressionError
			}
			decompressed.WriteByte(compressed[1])
			compressed = compressed[2:]
		case 255: // Verbatim string
			if len(compressed) < 2 {
				return nil, DecompressionError
			}
			n := int(compressed[1])
			if len(compressed) < n+2 {
				return nil, DecompressionError
			}
			decompressed.Write(compressed[2 : n+2])
			compressed = compressed[n+2:]
		default: // Look up encoded value
			decompressed.Write(codes[int(compressed[0])])
			compressed = compressed[1:]
		}
	}

	return decompressed.Bytes(), nil
}
