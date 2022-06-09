package tokenizer

const (
	UnorderedList     TokenType = "UnorderedList"
	UnorderedListItem           = "UnorderedListItem"
)

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

		if tokens, skip_bloc := parseCodeBlockWithSpaces(lines, index, spaces); tokens != nil {
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
