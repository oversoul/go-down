package tokenizer

import (
	"testing"
)

func tokenValid(t *Token, ttype TokenType, value string) bool {
	return (t.Ttype == ttype && t.Value == value)
}

func TestHeadingLevel1(t *testing.T) {
	tokens := Tokenize("# Hello world")

	if len(tokens) < 1 {
		t.Error("No tokens.")
	}
	if !tokenValid(tokens[0], Heading1, "Hello world") {
		t.Fail()
	}
}

func TestHeadingLevel2(t *testing.T) {
	tokens := Tokenize("## Hello world")

	if len(tokens) != 1 {
		t.Fail()
	}
	if !tokenValid(tokens[0], Heading2, "Hello world") {
		t.Fail()
	}
}

func TestHeadingLevel3(t *testing.T) {
	tokens := Tokenize("### Hello world")

	if len(tokens) != 1 {
		t.Fail()
	}
	if !tokenValid(tokens[0], Heading3, "Hello world") {
		t.Fail()
	}
}

func TestHeadingLevel4(t *testing.T) {
	tokens := Tokenize("#### Hello world")

	if len(tokens) != 1 {
		t.Fail()
	}
	if !tokenValid(tokens[0], Heading4, "Hello world") {
		t.Fail()
	}
}

func TestHeadingLevel5(t *testing.T) {
	tokens := Tokenize("##### Hello world")

	if len(tokens) != 1 {
		t.Fail()
	}
	if !tokenValid(tokens[0], Heading5, "Hello world") {
		t.Fail()
	}
}

func TestHeadingLevel6(t *testing.T) {
	tokens := Tokenize("###### Hello world")

	if len(tokens) != 1 {
		t.Fail()
	}

	if !tokenValid(tokens[0], Heading6, "Hello world") {
		t.Fail()
	}
}

func TestHeadingIgnoredIfSpaceDoesntFollowAfterHash(t *testing.T) {
	tokens := Tokenize("###Hello world")

	if len(tokens) != 1 {
		t.Fail()
	}
	if !tokenValid(tokens[0], Paragraph, "###Hello world") {
		t.Error("Should not be able to parse header.")
	}
}

func TestCodeBlock(t *testing.T) {
	tokens := Tokenize("```text\nHello world\n```")

	if len(tokens) != 1 {
		t.Fail()
	}

	if !tokenValid(tokens[0], CodeBloc, "Hello world") {
		t.Fail()
	}
}

func TestCodeBlockWithLanguage(t *testing.T) {
	tokens := Tokenize("```go\nconst msg := \"Hello world\"\n```")

	if len(tokens) < 1 {
		t.Error("no tokens")
	}
	if !tokenValid(tokens[0], CodeBloc, "const msg := \"Hello world\"") {
		t.Error("Not codeblock token.")
	}
	if tokens[0].Attrs["language"] != "go" {
		t.Error("Language not detected.")
	}
}

func TestBlocquote(t *testing.T) {
	tokens := Tokenize("> Hello world")

	if len(tokens) < 1 {
		t.Error("Not enough tokens")
	}
	if !tokenValid(tokens[0], Blockquote, "Hello world") {
		t.Error("Not valid Blockquote")
	}
}

func TestBlocquoteMultipleLines(t *testing.T) {
	tokens := Tokenize("> Hello world\n> Something else")

	if len(tokens) < 1 {
		t.Error("Not enough tokens")
	}

	if !tokenValid(tokens[0], Blockquote, "Hello world") {
		t.Error("Not valid Blockquote")
	}
}

func TestUnorederList(t *testing.T) {
	tokens := Tokenize("- First item\n- Second item\n- Third item\n- Fourth item")

	if len(tokens) < 1 {
		t.Error("Not enough tokens")
	}

	if !tokenValid(tokens[0], UnorderedList, "") {
		t.Error("Not valid UnorderedList")
	}

	if len(tokens[0].Children) != 4 {
		t.Error("Should parse 4 items.")
	}

	if tokens[0].Children[2].Value != "Third item" {
		t.Error("Third item is not valid.")
	}
}

func TestNestedUnorederList(t *testing.T) {
	tokens := Tokenize(`
- First item
- Second item
- Third item
  - Indented item
  - Indented item
- Fourth item`)

	if len(tokens) < 1 {
		t.Error("Not enough tokens")
	}

	if !tokenValid(tokens[0], UnorderedList, "") {
		t.Error("Not valid UnorderedList")
	}
	value := tokens[0].Children[2].Children[0].Value
	if value != "Indented item" {
		t.Errorf("Not valid indented unordered list item. `%s`", value)
	}
}

func TestSimpleCodeBlocWithinUnorderedList(t *testing.T) {
	tokens := Tokenize("* Winter\n  ```jsx\n  const Snow = <Snowflake amount=20 />;\n  ```\n* Frost")

	if len(tokens) < 1 {
		t.Error("Not enough tokens")
	}

	if !tokenValid(tokens[0], UnorderedList, "") {
		t.Error("Not valid UnorderedList")
	}

	if tokens[0].Children[0].Value != "Winter" {
		t.Error("Not valid UnorderedListItem Winter")
	}

	if tokens[0].Children[0].Children[0].Ttype != CodeBloc {
		t.Error("Not valid CodeBloc")
	}

	if tokens[0].Children[1].Value != "Frost" {
		t.Error("Not valid UnorderedListItem Frost")
	}
}

func TestParagraphWithinUnorderedList(t *testing.T) {
	tokens := Tokenize("* Winter\n* Frost\n  hello")

	if len(tokens) < 1 {
		t.Error("Not enough tokens")
	}
	if !tokenValid(tokens[0], UnorderedList, "") {
		t.Error("Not valid UnorderedList")
	}
	if tokens[0].Children[0].Value != "Winter" {
		t.Error("Not valid UnorderedListItem Winter")
	}
	if tokens[0].Children[1].Value != "Frost" {
		t.Error("Not valid UnorderedListItem Frost")
	}
	if tokens[0].Children[1].Children[0].Ttype != Paragraph {
		t.Error("Not valid Paragraph")
	}
}

func TestOrderedList(t *testing.T) {
	tokens := Tokenize("1. First item\n2. Second item\n")

	if len(tokens) < 1 {
		t.Error("Not enough tokens")
	}
	if !tokenValid(tokens[0], OrderedList, "") {
		t.Error("Not valid OrderedList")
	}
	item := tokens[0].Children[0]
	if item.Value != "First item" || item.Attrs["id"] != 1 {
		t.Errorf("Not valid item. `%s`", item.Value)
	}
	item = tokens[0].Children[1]
	if item.Value != "Second item" || item.Attrs["id"] != 2 {
		t.Errorf("Not valid item. `%s`", item.Value)
	}
}

// func TestOrderedListWithSubItems(t *testing.T) {
// 	tokens := Tokenize("1. First item\n  1.1 \n2. Second item\n")

// 	if len(tokens) < 1 {
// 		t.Error("Not enough tokens")
// 	}
// 	if !tokenValid(tokens[0], OrderedList, "") {
// 		t.Error("Not valid OrderedList")
// 	}
// 	value := tokens[0].Children[0].Value
// 	if value != "First item" {
// 		t.Errorf("Not valid item. `%s`", value)
// 	}
// 	value = tokens[0].Children[1].Value
// 	if value != "Second item" {
// 		t.Errorf("Not valid item. `%s`", value)
// 	}
// }
