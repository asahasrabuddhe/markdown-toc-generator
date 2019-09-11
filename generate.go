package main

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var (
	hashPatternHeader   = regexp.MustCompile("^(?P<indent>#+) ?(?P<title>.+)$")
	underscorePatternH1 = regexp.MustCompile("^=+$")
	underscorePatternH2 = regexp.MustCompile("^-+$")
	bullet              = map[int]string{
		0: "*",
		1: "-",
		2: "+",
	}
)

func generateToc(input []byte, depth, skipHeaders int) ([]byte, error) {
	var builder bytes.Buffer

	builder.WriteString("<!-- toc -->\n")

	scanner := bufio.NewScanner(bytes.NewReader(input))

	var previousLine string
	parsedHeaders := make(map[string]int)
	indentDiff := skipHeaders
	for scanner.Scan() {
		switch {
		case hashPatternHeader.Match(scanner.Bytes()):
			matches := hashPatternHeader.FindStringSubmatch(scanner.Text())
			if depth > 0 && len(matches[1]) > depth {
				continue
			}
			if strings.Contains(previousLine, "`") {
				continue
			}
			appendToToc(&builder, matches[2], len(matches[1])-1 - indentDiff, parsedHeaders, &skipHeaders)
		case underscorePatternH1.Match(scanner.Bytes()):
			appendToToc(&builder, previousLine, 0, parsedHeaders, &skipHeaders)
		case underscorePatternH2.Match(scanner.Bytes()):
			if depth > 0 && depth < 2 {
				continue
			}
			appendToToc(&builder, previousLine, 1 - indentDiff, parsedHeaders, &skipHeaders)
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

	builder.WriteString(fmt.Sprintf("%s%s [%s](#%s)\n", strings.Repeat("   ", indent), bullet[indent%len(bullet)], title, link))
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

	return strings.ToLower(out)
}
