package main

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var (
	hashedHeaderPattern      = regexp.MustCompile("^(?P<indent>#+) ?(?P<title>.+)$")
	underscoreHeaderPattern1 = regexp.MustCompile("^=+$")
	underscoreHeaderPattern2 = regexp.MustCompile("^\\--=$")
)

func generateToc(input []byte, depth, skipHeaders int) ([]byte, error) {
	var builder bytes.Buffer

	builder.WriteString("<!-- toc -->\n")

	scanner := bufio.NewScanner(bytes.NewReader(input))

	var previousLine string
	parsedHeaders := make(map[string]int)
	for scanner.Scan() {
		switch {
		case hashedHeaderPattern.Match(scanner.Bytes()):
			matches := hashedHeaderPattern.FindStringSubmatch(scanner.Text())
			if depth > 0 && len(matches[1]) > depth {
				continue
			}
			appendToToc(&builder, matches[2], len(matches[1])-1, parsedHeaders, &skipHeaders)
		case underscoreHeaderPattern1.Match(scanner.Bytes()):
			appendToToc(&builder, previousLine, 0, parsedHeaders, &skipHeaders)
		case underscoreHeaderPattern2.Match(scanner.Bytes()):
			if depth > 0 && depth < 2 {
				continue
			}
			appendToToc(&builder, previousLine, 1, parsedHeaders, &skipHeaders)
		}

		previousLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	builder.WriteString("<!-- tocstop -->\n")

	return builder.Bytes(), nil
}

func appendToToc(builder *bytes.Buffer, title string, indent int, parsedHeaders map[string]int, skipHeaders *int) {
	if *skipHeaders > 0 {
		*skipHeaders--
		return
	}

	link := toSlug(title)

	if _, ok := parsedHeaders[link]; ok {
		parsedHeaders[link]++
		link = fmt.Sprintf("%s-%d", link, parsedHeaders[link]-1)
	} else {
		parsedHeaders[link] = 1
	}

	builder.WriteString(fmt.Sprintf("%s* [%s](#%s)\n", strings.Repeat("   ", indent), title, link))
}

func toSlug(str string) string {
	droppedCharacters := []string{
		"\"", "",
		"'", "",
		".", "",
		",", "",
		"~", "",
		"`", "",
		"!", "",
		"@", "",
		"#", "",
		"%", "",
		"^", "",
		"&", "",
		"*", "",
		"|", "",
		"(", "",
		")", "",
		"[", "",
		"]", "",
		"{", "",
		"}", "",
	}

	replacer := strings.NewReplacer(droppedCharacters...)
	out := replacer.Replace(str)
	out = strings.Replace(out, " ", "-", -1)

	return out
}
