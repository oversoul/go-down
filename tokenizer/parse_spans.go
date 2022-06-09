package tokenizer

const (
	TextNormal TokenType = "TextNormal"
	TextBold   TokenType = "TextBold"
	TextItalic           = "TextItalic"
)

func findEndingBold(line string, start int, chars string) int {
	i := start
	for i < len(line)-1 {
		if line[i:i+2] == chars {
			return i
		}
		i++
	}
	return -1
}

func findEndingItalic(line string, start int, chars byte) int {
	i := start
	for i < len(line) {
		if line[i] == chars {
			return i
		}
		i++
	}
	return -1
}

func isBoldSpan(line string, i int) bool {
	return len(line) >= i+2 && (line[i:i+2] == "**" || line[i:i+2] == "__")
}

func isItalicSpan(line string, i int) bool {
	return line[i] == '*' || line[i] == '_'
}

func parseSpans(line string) *Token {
	token := newToken(Paragraph, "")
	i := 0
	lineSize := len(line)
	normalText := []byte{}
	for i < lineSize {

		if isBoldSpan(line, i) {
			end := findEndingBold(line, i+2, line[i:i+2])
			if end > 0 {
				if len(normalText) > 0 {
					token.Children = append(
						token.Children,
						newToken(TextNormal, string(normalText)),
					)
					normalText = []byte{}
				}
				token.Children = append(
					token.Children,
					newToken(TextBold, line[i+2:end]),
				)
				i += end
				continue
			}
		}

		if isItalicSpan(line, i) {
			end := findEndingItalic(line, i+1, line[i])
			if end > 0 {
				if len(normalText) > 0 {
					token.Children = append(
						token.Children,
						newToken(TextNormal, string(normalText)),
					)
					normalText = []byte{}
				}
				token.Children = append(
					token.Children,
					newToken(TextItalic, line[i+1:end]),
				)
				i += end
				continue
			}
		}

		normalText = append(normalText, line[i])
		i++
	}

	if len(normalText) > 0 {
		token.Children = append(
			token.Children,
			newToken(TextNormal, string(normalText)),
		)
	}

	return token
}
