package main

import (
	"bufio"
	"bytes"
	"regexp"
)

var (
	tocStartRegex = regexp.MustCompile("^<!\\-\\- toc \\-\\->")
	tocEndRegex   = regexp.MustCompile("^<!\\-\\- tocstop \\-\\->")
)

func replaceToc(input []byte, output []byte) []byte {
	var (
		buf                           bytes.Buffer
		tocBlockFound, insideTocBlock bool
	)
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		if tocStartRegex.Match(scanner.Bytes()) {
			tocBlockFound = true
			insideTocBlock = true
			buf.Write(output)
		}

		if !insideTocBlock {
			buf.Write(scanner.Bytes())
			buf.WriteString("\n")
		}

		if tocEndRegex.Match(scanner.Bytes()) {
			insideTocBlock = false
		}
	}
	if !tocBlockFound {
		buf.Reset()
		buf.Write(output)
		buf.Write(input)
	}
	output = buf.Bytes()
	return output
}
