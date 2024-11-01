package main

import (
	"fmt"
	"strings"
	"unicode"
)

type UnicodeWriterObject struct {
	idx       int
	remaining int
	result    strings.Builder
}

func (uw *UnicodeWriterObject) UnicodeWriter(s string, cleanedStr strings.Builder) strings.Builder {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleanedStr.WriteRune(r)
		}
	}
	return cleanedStr
}

func (uw *UnicodeWriterObject) Append(n int, cleaned string) {
	for i := 0; i < (n / 3); i++ {
		if uw.result.Len() > 0 {
			uw.result.WriteString(" ")
		}
		uw.result.WriteString(cleaned[uw.idx : uw.idx+3])
		uw.idx += 3
	}
}

func (uw *UnicodeWriterObject) Handle(remaining int, cleaned string) string {
	if remaining == 1 {
		// If remaining is 1, append it to the last block to make it of length 2.
		lastBlockStart := len(uw.result.String()) - 4
		lastBlock := uw.result.String()[lastBlockStart:] + cleaned[uw.idx:]
		resultStr := uw.result.String()[:lastBlockStart]
		uw.result.Reset()
		uw.result.WriteString(resultStr)
		uw.result.WriteString(" ")
		uw.result.WriteString(lastBlock)
	} else if remaining > 0 {
		// If remaining is 2, append it as a new block.
		if uw.result.Len() > 0 {
			uw.result.WriteString(" ")
		}
		uw.result.WriteString(cleaned[uw.idx:])
	}
	return uw.result.String()

}

func FormatString(s string) string {
	// Step 1: Remove all spaces and dashes, and collect alphanumeric characters.
	var cleanedStr strings.Builder

	uw := UnicodeWriterObject{
		idx:    0,
		result: strings.Builder{},
	}
	cleanedStr = uw.UnicodeWriter(s, cleanedStr)

	// Step 2: Split the cleaned string into blocks of length 3.
	cleaned := cleanedStr.String()
	n := len(cleaned)
	remaining := n % 3

	// Step 3: Append blocks of length 3 with a space.
	uw.Append(n, cleaned)

	// Step 4: Handle the remaining characters.
	return uw.Handle(remaining, cleaned)
}

func main() {
	fmt.Println(FormatString("00-44  48 5555 8361"))       // Expected output: "004 448 555 583 61"
	fmt.Println(FormatString("0 - 22 1985--324"))          // Expected output: "022 198 532 4"
	fmt.Println(FormatString("ABC372654"))                 // Expected output: "ABC 372 654"
	fmt.Println(FormatString("AA-44  BB  5CD  85C D83FG")) // Expected output: "AA4 4BB 5C5 D83 FG"
}
