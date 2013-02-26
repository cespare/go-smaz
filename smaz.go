// Package smaz is an implementation of the smaz library (https://github.com/antirez/smaz) for compressing
// small strings.
package smaz

import (
	"bytes"
	"errors"
)

var codes = []string{" ",
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

var codeArrays = make([][]byte, len(codes))
var prefixToCode = make(map[string]byte)
var maxCodeLength = 0 // TODO: Unnecessary when we switch to a trie implementation

// Library initialization.
func init() {
	for i, code := range codes {
		codeArrays[i] = []byte(code)
		prefixToCode[code] = byte(i)
		if len(code) > maxCodeLength {
			maxCodeLength = len(code)
		}
	}
}

// BUG(cespare): Compress is written in an extremely naive manner for the time being and is very slow. I will
// reimplement (and then profile/optimize it) after I get done with go-trie, which will be a natural fit for
// this problem.

// Compress compresses a byte slice and returns the compressed data.
func Compress(input []byte) []byte {
	var outputBuffer bytes.Buffer
	var verbatim bytes.Buffer
	remaining := len(input)
	position := 0

	flushVerbatim := func() {
		// We can write a max of 255 continuous verbatim characters, because the length of the continous verbatim
		// section is represented by a single byte.
		for verbatim.Len() > 0 {
			chunk := verbatim.Next(255)
			if len(chunk) == 1 {
				// 254 is code for a single verbatim byte
				outputBuffer.WriteByte(byte(254))
			} else {
				// 255 is code for a verbatim string. It is followed by a byte containing the length of the string.
				outputBuffer.WriteByte(byte(255))
				outputBuffer.WriteByte(byte(len(chunk)))
			}
			outputBuffer.Write(chunk)
		}
		verbatim.Reset()
	}

	for remaining > 0 {
		// Find the longest matching substring, brute-force
		longestPossibleMatch := maxCodeLength
		if remaining < longestPossibleMatch {
			longestPossibleMatch = remaining
		}
		matchFound := false
		for i := longestPossibleMatch; i > 0; i-- {
			prefix := input[position : position+i]
			/*fmt.Printf("Prefix: %v\n", string(prefix))*/
			if code, ok := prefixToCode[string(prefix)]; ok {
				// Match found
				remaining -= i
				position += i
				flushVerbatim()
				outputBuffer.WriteByte(code)
				matchFound = true
				break
			}
		}
		if !matchFound {
			verbatim.WriteByte(input[position])
			remaining -= 1
			position += 1
		}
	}
	flushVerbatim()

	return outputBuffer.Bytes()
}

var decompressionError = errors.New("Invalid or corrupted compressed data.")

// Decompress decompresses a smaz-compressed byte slice and return a new slice with the decompressed data. err
// is nil if and only if decompression fails for any reason (e.g., corrupted data).
func Decompress(compressed []byte) ([]byte, error) {
	var decompressed bytes.Buffer
	var remaining = len(compressed)
	var position = 0

	for remaining > 0 {
		switch compressed[position] {
		case 254:
			// Verbatim byte
			if remaining < 2 {
				return nil, decompressionError
			}
			decompressed.WriteByte(compressed[position+1])
			remaining -= 2
			position += 2
		case 255:
			// Verbatim string
			if remaining < 2 {
				return nil, decompressionError
			}
			length := int(compressed[position+1])
			if remaining < length+2 {
				return nil, decompressionError
			}
			decompressed.Write(compressed[position+2 : position+length+2])
			remaining -= length + 2
			position += length + 2
		default:
			// Look up encoded value
			decompressed.Write([]byte(codes[int(compressed[position])]))
			remaining--
			position++
		}
	}

	return decompressed.Bytes(), nil
}
