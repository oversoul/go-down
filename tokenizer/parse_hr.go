package tokenizer

const (
	Hr TokenType = "Hr"
)

func parseHr(lines []string, index int) ([]*Token, int) {
	if isEmpty(lines[index]) {
		return nil, 0
	}
	if lines[index] == "---" || lines[index] == "===" {
		return []*Token{newToken(Hr, "")}, 1
	}
	return nil, 0
}
