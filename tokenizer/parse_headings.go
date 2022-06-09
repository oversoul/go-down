package tokenizer

const (
	Heading1 TokenType = "Heading1"
	Heading2           = "Heading2"
	Heading3           = "Heading3"
	Heading4           = "Heading4"
	Heading5           = "Heading5"
	Heading6           = "Heading6"
)

func parseHeading(lines []string, index int) ([]*Token, int) {
	line := lines[index]
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

	headings := map[int]TokenType{
		1: Heading1,
		2: Heading2,
		3: Heading3,
		4: Heading4,
		5: Heading5,
		6: Heading6,
	}

	if value, found := headings[i]; found {
		return []*Token{newToken(value, line[i+1:])}, 1
	}

	return nil, 0
}
