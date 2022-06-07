package tokenizer

import (
	"strings"
)

type TokenType string

const (
	Hr                TokenType = "Hr"
	Heading1                    = "Heading1"
	Heading2                    = "Heading2"
	Heading3                    = "Heading3"
	Heading4                    = "Heading4"
	Heading5                    = "Heading5"
	Heading6                    = "Heading6"
	CodeBloc                    = "CodeBloc"
	Paragraph                   = "Paragraph"
	Blockquote                  = "Blockquote"
	UnorderedList               = "UnorderedList"
	UnorderedListItem           = "UnorderedListItem"
)

type Attribute map[string]any

type Token struct {
	parent   *Token
	Ttype    TokenType `json:"type"`
	Value    string    `json:"value"`
	Children []*Token  `json:"children"`
	Attrs    Attribute `json:"attributes"`
}

type Renderable interface {
	Render()
}

func newToken(ttype TokenType, value string) *Token {
	return &Token{Ttype: ttype, Value: value, Children: []*Token{}, Attrs: Attribute{}}
}

func isEmpty(line string) bool {
	return len(strings.TrimSpace(line)) == 0
}

func parseBlocquote(lines []string, index int) ([]*Token, int) {
	if isEmpty(lines[index]) {
		return nil, 0
	}

	blockLines := []*Token{}

	for index < len(lines) {
		if !isEmpty(lines[index]) && lines[index][0:2] == "> " {
			blockLines = append(blockLines, newToken(Blockquote, lines[index][2:]))
			index++
			continue
		}
		index++
	}

	return blockLines, len(blockLines)
}

func parseHeading(line string) *Token {
	if isEmpty(line) {
		return nil
	}
	i := 0
	for line[i] == '#' {
		i++
	}

	if line[i] != ' ' {
		return nil
	}

	switch i {
	case 1:
		return newToken(Heading1, line[i+1:])
	case 2:
		return newToken(Heading2, line[i+1:])
	case 3:
		return newToken(Heading3, line[i+1:])
	case 4:
		return newToken(Heading4, line[i+1:])
	case 5:
		return newToken(Heading5, line[i+1:])
	case 6:
		return newToken(Heading6, line[i+1:])
	default:
		return nil
	}
}

func parseHr(line string) (*Token, int) {
	if isEmpty(line) {
		return nil, 0
	}
	if line[0:3] == "---" || line[0:3] == "===" {
		return newToken(Hr, ""), 1
	}
	return nil, 0
}

func parseCodeBloc(lines []string, index int) (*Token, int) {
	if len(strings.TrimSpace(lines[index])) == 0 {
		return nil, 0
	}
	if lines[index][0:3] != "```" {
		return nil, 0
	}

	language := strings.TrimSpace(lines[index][3:])
	blocLines := 0
	i := index + 1
	for i < len(lines) {
		blocLines++
		if lines[i][0:3] == "```" {
			i++
			break
		}
		i++
	}

	token := newToken(CodeBloc, strings.Join(lines[index+1:blocLines], "\n"))
	token.Attrs["language"] = language
	return token, blocLines
}

func parseUnorderedList(lines []string, index int) (*Token, int) {
	if isEmpty(lines[index]) {
		return nil, 0
	}

	isList := func(line string) bool {
		return line == "- " || line == "+ " || line == "* "
	}

	i := index
	skip := 0
	list := newToken(UnorderedList, "")
	current := list
	for i < len(lines) {
		if isEmpty(lines[i]) {
			i++
			continue
		}
		subI := 0
		for lines[i][subI] == '\t' {
			subI++
		}

		firstChars := lines[i][subI : subI+2]

		if isList(firstChars) {
			start := subI
			for subI > 0 {
				if current != nil {
					current = current.Children[len(current.Children)-1]
				}
				subI--
			}
			current.Children = append(current.Children, newToken(UnorderedListItem, lines[i][start+2:]))
			current = list
			skip++
		}
		i++
	}

	if skip == 0 {
		return nil, 0
	}

	return list, skip
}

func parseParagraph(lines []string, index int) (*Token, int) {
	return nil, 0
}

func Tokenize(content string) []*Token {
	lines := strings.Split(content, "\n")
	i := 0
	tokens := []*Token{}
	for i < len(lines) {
		if token, skip := parseHr(lines[i]); token != nil {
			tokens = append(tokens, token)
			i += skip
		}

		if token := parseHeading(lines[i]); token != nil {
			tokens = append(tokens, token)
			i++
		}

		if token, skip := parseCodeBloc(lines, i); token != nil {
			tokens = append(tokens, token)
			i += skip
		}

		if blocks, skip := parseBlocquote(lines, i); len(blocks) > 0 {
			tokens = append(tokens, blocks...)
			i += skip
		}

		if token, skip := parseUnorderedList(lines, i); token != nil {
			tokens = append(tokens, token)
			i += skip
		}

		// parse_ordered_list
		// parse_link_reference

		if token, skip := parseParagraph(lines, i); token != nil {
			tokens = append(tokens, token)
			i += skip
		}
		// tokens = append(tokens, newToken(Paragraph, lines[i]))

		i++
	}
	return tokens
}
