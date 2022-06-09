package tokenizer

import "strings"

const (
	CodeBloc TokenType = "CodeBloc"
)

func parseCodeBlock(lines []string, index int) ([]*Token, int) {
	return parseCodeBlockWithSpaces(lines, index, 0)
}

func parseCodeBlockWithSpaces(lines []string, index int, spaces int) ([]*Token, int) {
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
