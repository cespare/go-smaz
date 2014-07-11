package smaz

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

var antirezTestStrings = []string{"",
	"This is a small string",
	"foobar",
	"the end",
	"not-a-g00d-Exampl333",
	"Smaz is a simple compression library",
	"Nothing is more difficult, and therefore more precious, than to be able to decide",
	"this is an example of what works very well with smaz",
	"1000 numbers 2000 will 10 20 30 compress very little",
	"and now a few italian sentences:",
	"Nel mezzo del cammin di nostra vita, mi ritrovai in una selva oscura",
	"Mi illumino di immenso",
	"L'autore di questa libreria vive in Sicilia",
	"try it against urls",
	"http://google.com",
	"http://programming.reddit.com",
	"http://github.com/antirez/smaz/tree/master",
	"/media/hdb1/music/Alben/The Bla",
}

func TestCorrectness(t *testing.T) {
	// Set up our slice of test strings.
	inputs := make([][]byte, 0)
	for _, s := range antirezTestStrings {
		inputs = append(inputs, []byte(s))
	}
	// An array with every possible byte value in it.
	allBytes := make([]byte, 256)
	for i := 0; i < 256; i++ {
		allBytes[i] = byte(i)
	}
	inputs = append(inputs, allBytes)
	// A long array of all 0s (the longest continuous string that can be represented is 256; any longer than
	// this and the compressor will need to split it into chunks)
	allZeroes := make([]byte, 300)
	for i := 0; i < 300; i++ {
		allZeroes[i] = byte(0)
	}
	inputs = append(inputs, allZeroes)

	for _, input := range inputs {
		compressed := Compress(input)
		decompressed, err := Decompress(compressed)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(input, decompressed) {
			t.Fatal("want %q after decompression; got %q\n", input, decompressed)
		}

		if len(input) > 1 && len(input) < 50 {
			compressionLevel := 100 - ((100.0 * len(compressed)) / len(input))
			if compressionLevel < 0 {
				t.Logf("%q enlarged by %d%%\n", input, -compressionLevel)
			} else {
				t.Logf("%q compressed by %d%%\n", input, compressionLevel)
			}
		}
	}
}

func loadTestData(t testing.TB) [][]byte {
	f, err := os.Open("./testdata/pg5200.txt")
	if err != nil {
		t.Fatal(err)
	}

	var lines [][]byte
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, []byte(scanner.Text())) // Note that .Bytes() would require us to manually copy
	}
	return lines
}

func BenchmarkCompression(b *testing.B) {
	b.StopTimer()
	inputs := loadTestData(b)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			Compress(input)
		}
	}
}

func BenchmarkDecompression(b *testing.B) {
	b.StopTimer()
	inputs := loadTestData(b)
	compressedStrings := make([][]byte, len(inputs))
	for i, input := range inputs {
		compressedStrings[i] = Compress(input)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, compressed := range compressedStrings {
			Decompress(compressed)
		}
	}
}
