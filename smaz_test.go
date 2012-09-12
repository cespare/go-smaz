package smaz

import (
	"bufio"
	"bytes"
	"fmt"
	. "launchpad.net/gocheck"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type SmazSuite struct{}

var _ = Suite(&SmazSuite{})

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

func (s *SmazSuite) TestCorrectness(c *C) {
	// Set up our slice of test strings.
	inputs := make([][]byte, 0)
	for _, testInput := range antirezTestStrings {
		inputs = append(inputs, []byte(testInput))
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

	for _, testInput := range inputs {
		compressed := Compress(testInput)
		decompressed, err := Decompress(compressed)
		c.Assert(err, IsNil)
		if len(testInput) == 0 {
			// Can't use DeepEquals for a nil slice and an empty slice -- they're different
			c.Assert(decompressed, HasLen, 0)
		} else {
			c.Assert(testInput, DeepEquals, decompressed)
		}

		if len(testInput) > 1 && len(testInput) < 50 {
			compressionLevel := 100 - ((100.0 * len(compressed)) / len(testInput))
			if compressionLevel < 0 {
				fmt.Printf("'%s' enlarged by %d%%\n", testInput, -compressionLevel)
			} else {
				fmt.Printf("'%s' compressed by %d%%\n", testInput, compressionLevel)
			}
		}
	}
}

func loadTestData() [][]byte {
	file, err := os.Open("./testdata/pg5200.txt")
	if err != nil {
		fmt.Printf("Error opening test data file: %v\n", err)
		os.Exit(1)
	}

	totalSize := 0
	testStrings := make([][]byte, 0)
	currentLine := new(bytes.Buffer)
	reader := bufio.NewReader(file)
	var part []byte
	var prefix bool

	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		currentLine.Write(part)
		totalSize += len(part)
		if !prefix {
			testStrings = append(testStrings, currentLine.Bytes())
			currentLine = new(bytes.Buffer)
		}
	}
	return testStrings
}

func (s *SmazSuite) BenchmarkCompression(c *C) {
	c.StopTimer()
	testStrings := loadTestData()
	/*fmt.Printf("The test corpus contains %d lines and %d bytes of text.", len(testStrings), totalSize)*/
	c.StartTimer()
	for i := 0; i < c.N; i++ {
		for _, testString := range testStrings {
			Compress(testString)
		}
	}
}

func (s *SmazSuite) BenchmarkDecompression(c *C) {
	c.StopTimer()
	testStrings := loadTestData()
	compressedStrings := make([][]byte, len(testStrings))
	for i, testString := range testStrings {
		compressedStrings[i] = Compress(testString)
	}
	c.StartTimer()
	for i := 0; i < c.N; i++ {
		for _, compressed := range compressedStrings {
			Decompress(compressed)
		}
	}
}
