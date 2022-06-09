package tokenizer

import (
	"strconv"
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
	OrderedList                 = "OrderedList"
	OrderedListItem             = "OrderedListItem"
)

type Attribute map[string]any

type Token struct {
	Ttype    TokenType `json:"type"`
	Value    string    `json:"value"`
	Children []*Token  `json:"children"`
	Attrs    Attribute `json:"attributes"`
}

type Renderable interface {
	Render()
}

func newToken(ttype TokenType, value string) *Token {
	return &Token{
		Ttype:    ttype,
		Value:    value,
		Attrs:    Attribute{},
		Children: []*Token{},
	}
}

func isEmpty(line string) bool {
	return len(strings.TrimSpace(line)) == 0
}

func countSpaces(line string) int {
	spaces := 0
	for line[spaces] == ' ' {
		spaces++
	}
	return spaces
}

func parseBlocquote(lines []string, index int) ([]*Token, int) {
	if isEmpty(lines[index]) {
		return nil, 0
	}

	blockLines := []*Token{}

	for index < len(lines) {
		if isEmpty(lines[index]) {
			break
		}

		spaces := countSpaces(lines[index])

		if lines[index][spaces:spaces+2] != "> " {
			break
		}

		blockLines = append(blockLines, newToken(Blockquote, lines[index][spaces+2:]))
		index++
	}

	return blockLines, len(blockLines)
}

func parseHeading(line string) ([]*Token, int) {
	if isEmpty(line) {
		return nil, 0
	}
	i := 0
	for line[i] == '#' {
		i++
	}

	if line[i] != ' ' {
		return nil, 0
	}

	var ttype TokenType

	switch i {
	case 1:
		ttype = Heading1
	case 2:
		ttype = Heading2
	case 3:
		ttype = Heading3
	case 4:
		ttype = Heading4
	case 5:
		ttype = Heading5
	case 6:
		ttype = Heading6
	default:
		return nil, 0
	}

	return []*Token{newToken(ttype, line[i+1:])}, 1
}

func parseHr(line string) (*Token, int) {
	if isEmpty(line) {
		return nil, 0
	}
	if line == "---" || line == "===" {
		return newToken(Hr, ""), 1
	}
	return nil, 0
}

func parseCodeBloc(lines []string, index int, spaces int) ([]*Token, int) {
	if len(lines[index]) < 4+spaces {
		return nil, 0
	}
	if lines[index][spaces:spaces+3] != "```" {
		return nil, 0
	}

	language := strings.TrimSpace(lines[index][spaces+3:])
	blocLines := []string{}
	i := index + 1
	for i < len(lines) {
		if len(lines[i]) >= 3 && lines[i][spaces:spaces+3] == "```" {
			break
		}
		blocLines = append(blocLines, strings.TrimSpace(lines[i]))
		i++
	}

	token := newToken(CodeBloc, strings.Join(blocLines, "\n"))

	token.Attrs["language"] = language
	return []*Token{token}, len(blocLines) + 2
}

func parseUnorderedList(lines []string, index int) ([]*Token, int) {
	if isEmpty(lines[index]) {
		return nil, 0
	}

	isList := func(line string) bool {
		return line == "- " || line == "+ " || line == "* "
	}

	skip := 0
	list := newToken(UnorderedList, "")
	current := list
	for index < len(lines) {
		if isEmpty(lines[index]) {
			break
		}

		spaces := countSpaces(lines[index])

		subI := spaces / 2

		firstChars := lines[index][spaces : spaces+2]

		if !isList(firstChars) && spaces == 0 {
			break
		}

		for subI > 0 {
			if current != nil && len(current.Children) > 0 {
				current = current.Children[len(current.Children)-1]
			}
			subI--
		}

		inc_value := 1
		skip_value := 1

		if tokens, skip_bloc := parseCodeBloc(lines, index, spaces); tokens != nil {
			current.Children = append(current.Children, tokens...)
			inc_value = skip_bloc
			skip_value = skip_bloc
		} else {
			if isList(firstChars) {
				current.Children = append(
					current.Children,
					newToken(UnorderedListItem, lines[index][spaces+2:]),
				)
			} else {
				current.Children = append(
					current.Children,
					newToken(Paragraph, lines[index][spaces:]),
				)
			}
		}

		current = list
		skip += skip_value
		index += inc_value
	}

	return []*Token{list}, skip
}

func parseOrderedList(lines []string, index int) ([]*Token, int) {
	if isEmpty(lines[index]) {
		return nil, 0
	}

	skip := 0
	list := newToken(OrderedList, "")
	current := list
	for index < len(lines) {
		if isEmpty(lines[index]) {
			break
		}

		spaces := countSpaces(lines[index])
		slices := strings.SplitN(lines[index][spaces:], ". ", 2)

		_, err := strconv.Atoi(slices[0])
		if err != nil && spaces == 0 {
			break
		}

		token := newToken(OrderedListItem, slices[1])
		token.Attrs["id"] = skip + 1

		current.Children = append(current.Children, token)
		skip += 1
		index += 1
	}

	return []*Token{list}, skip
}

func Tokenize(content string) []*Token {
	lines := strings.Split(content, "\n")
	i := 0
	tokens := []*Token{}
	for i < len(lines) {
		if token, skip := parseHr(lines[i]); token != nil {
			tokens = append(tokens, token)
			i += skip
			continue
		}

		if blocks, skip := parseHeading(lines[i]); skip > 0 {
			tokens = append(tokens, blocks...)
			i += skip
			continue
		}

		if blocks, skip := parseCodeBloc(lines, i, 0); skip > 0 {
			tokens = append(tokens, blocks...)
			i += skip
			continue
		}

		if blocks, skip := parseBlocquote(lines, i); skip > 0 {
			tokens = append(tokens, blocks...)
			i += skip
			continue
		}

		if blocks, skip := parseUnorderedList(lines, i); skip > 0 {
			tokens = append(tokens, blocks...)
			i += skip
			continue
		}

		if blocks, skip := parseOrderedList(lines, i); skip > 0 {
			tokens = append(tokens, blocks...)
			i += skip
			continue
		}

		if i >= len(lines) {
			break
		}

		if !isEmpty(lines[i]) {
			tokens = append(tokens, newToken(Paragraph, lines[i]))
		}

		i++
	}
	return tokens
}
