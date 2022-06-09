package tokenizer

import (
	"strconv"
	"strings"
)

const (
	OrderedList     TokenType = "OrderedList"
	OrderedListItem           = "OrderedListItem"
)

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
