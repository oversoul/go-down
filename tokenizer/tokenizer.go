package tokenizer

import (
	"strings"
)

type TokenType string

const (
	Paragraph TokenType = "Paragraph"
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

type ParserFunc func([]string, int) ([]*Token, int)

type parser struct {
	lines   []string
	parsers []ParserFunc
}

func NewParser(content string) *parser {
	return &parser{
		lines: strings.Split(content, "\n"),
		parsers: []ParserFunc{
			parseHr,
			parseHeading,
			parseCodeBlock,
			parseBlockquote,
			parseUnorderedList,
			parseOrderedList,
			parseParagraph,
		},
	}
}

func parseParagraph(lines []string, index int) ([]*Token, int) {
	if isEmpty(lines[index]) {
		return nil, 0
	}
	paragraph := newToken(Paragraph, "")
	paragraph.Children = parseSpans(lines[index])
	return []*Token{paragraph}, 1
}

func (p *parser) Tokenize() []*Token {
	i := 0
	tokens := []*Token{}

out:
	for i < len(p.lines) {
		for _, parser := range p.parsers {
			blocks, skip := parser(p.lines, i)
			if skip > 0 {
				tokens = append(tokens, blocks...)
				i += skip
				continue out
			}
		}
		i++
	}

	return tokens
}
