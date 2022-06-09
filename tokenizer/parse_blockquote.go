package tokenizer

const (
	Blockquote TokenType = "Blockquote"
)

func parseBlockquote(lines []string, index int) ([]*Token, int) {
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
