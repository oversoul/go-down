package tokenizer

import (
	"testing"
)

func tokenValid(t *Token, ttype TokenType, value string) bool {
	return (t.Ttype == ttype && t.Value == value)
}

func TestHeadingLevel1(t *testing.T) {
	tokens := NewParser("# Hello world").Tokenize()

	if len(tokens) < 1 {
		t.Error("No tokens.")
	}
	if !tokenValid(tokens[0], Heading1, "Hello world") {
		t.Fail()
	}
}

func TestHeadingLevel2(t *testing.T) {
	tokens := NewParser("## Hello world").Tokenize()

	if len(tokens) != 1 {
		t.Fail()
	}
	if !tokenValid(tokens[0], Heading2, "Hello world") {
		t.Fail()
	}
}

func TestHeadingLevel3(t *testing.T) {
	tokens := NewParser("### Hello world").Tokenize()

	if len(tokens) != 1 {
		t.Fail()
	}
	if !tokenValid(tokens[0], Heading3, "Hello world") {
		t.Fail()
	}
}

func TestHeadingLevel4(t *testing.T) {
	tokens := NewParser("#### Hello world").Tokenize()

	if len(tokens) != 1 {
		t.Fail()
	}
	if !tokenValid(tokens[0], Heading4, "Hello world") {
		t.Fail()
	}
}

func TestHeadingLevel5(t *testing.T) {
	tokens := NewParser("##### Hello world").Tokenize()

	if len(tokens) != 1 {
		t.Fail()
	}
	if !tokenValid(tokens[0], Heading5, "Hello world") {
		t.Fail()
	}
}

func TestHeadingLevel6(t *testing.T) {
	tokens := NewParser("###### Hello world").Tokenize()

	if len(tokens) != 1 {
		t.Fail()
	}

	if !tokenValid(tokens[0], Heading6, "Hello world") {
		t.Fail()
	}
}

func TestHeadingIgnoredIfSpaceDoesntFollowAfterHash(t *testing.T) {
	tokens := NewParser("###Hello world").Tokenize()

	if len(tokens) != 1 {
		t.Fail()
	}
	if !tokenValid(tokens[0], Paragraph, "###Hello world") {
		t.Error("Should not be able to parse header.")
	}
}

func TestCodeBlock(t *testing.T) {
	tokens := NewParser("```text\nHello world\n```").Tokenize()

	if len(tokens) != 1 {
		t.Fail()
	}

	if !tokenValid(tokens[0], CodeBloc, "Hello world") {
		t.Fail()
	}
}

func TestCodeBlockWithLanguage(t *testing.T) {
	tokens := NewParser("```go\nconst msg := \"Hello world\"\n```").Tokenize()

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
	tokens := NewParser("> Hello world").Tokenize()

	if len(tokens) < 1 {
		t.Error("Not enough tokens")
	}
	if !tokenValid(tokens[0], Blockquote, "Hello world") {
		t.Error("Not valid Blockquote")
	}
}

func TestBlocquoteMultipleLines(t *testing.T) {
	tokens := NewParser("> Hello world\n> Something else").Tokenize()

	if len(tokens) < 1 {
		t.Error("Not enough tokens")
	}

	if !tokenValid(tokens[0], Blockquote, "Hello world") {
		t.Error("Not valid Blockquote")
	}
}

func TestUnorederList(t *testing.T) {
	tokens := NewParser("- First item\n- Second item\n- Third item\n- Fourth item").Tokenize()

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
	tokens := NewParser(`
- First item
- Second item
- Third item
  - Indented item
  - Indented item
- Fourth item`).Tokenize()

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
	tokens := NewParser("* Winter\n  ```jsx\n  const Snow = <Snowflake amount=20 />;\n  ```\n* Frost").Tokenize()

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
	tokens := NewParser("* Winter\n* Frost\n  hello").Tokenize()

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
	tokens := NewParser("1. First item\n2. Second item\n").Tokenize()

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

/*
func TestSpanParserForBold(t *testing.T) {
	token := parseSpans("Hello **world**")

	if len(token.Children) < 1 {
		t.Error("Not enough tokens")
	}
	if !tokenValid(token.Children[0], TextNormal, "Hello ") {
		t.Error("Not valid Text normal")
	}
	if !tokenValid(token.Children[1], TextBold, "world") {
		t.Error("Not valid Bold text")
	}
}

func TestSpanParserForWrongBold(t *testing.T) {
	token := parseSpans("Hello *_world")
	if len(token.Children) < 1 {
		t.Error("Not enough tokens")
	}

	if !tokenValid(token.Children[0], TextNormal, "Hello *_world") {
		t.Error("Not valid Text normal")
	}
}

func TestSpanParserForBoldWithUnderscore(t *testing.T) {
	token := parseSpans("Hello __world__")

	if len(token.Children) < 1 {
		t.Error("Not enough tokens")
	}
	if !tokenValid(token.Children[0], TextNormal, "Hello ") {
		t.Error("Not valid Text normal")
	}
	if !tokenValid(token.Children[1], TextBold, "world") {
		t.Error("Not valid Bold text")
	}
}

func TestSpanParserForItalic(t *testing.T) {
	token := parseSpans("Hello *world*")
	if len(token.Children) < 2 {
		t.Error("Not enough tokens")
	}
	if !tokenValid(token.Children[0], TextNormal, "Hello ") {
		t.Error("Not valid Text normal")
	}
	if !tokenValid(token.Children[1], TextItalic, "world") {
		t.Error("Not valid Bold text")
	}
}

func TestSpanParserForWrongItalic(t *testing.T) {
	token := parseSpans("Hello *world")
	if len(token.Children) < 1 {
		t.Error("Not enough tokens")
	}
	if !tokenValid(token.Children[0], TextNormal, "Hello *world") {
		t.Error("Not valid Text normal")
	}
}

func TestSpanParserForItalicWithUnderscore(t *testing.T) {
	token := parseSpans("Hello _world_")
	if len(token.Children) < 2 {
		t.Error("Not enough tokens")
	}
	if !tokenValid(token.Children[0], TextNormal, "Hello ") {
		t.Error("Not valid Text normal")
	}
	if !tokenValid(token.Children[1], TextItalic, "world") {
		t.Error("Not valid Bold text")
	}
}

func TestSpanParserForWrongItalicUndercore(t *testing.T) {
	token := parseSpans("Hello _world")
	if len(token.Children) < 1 {
		t.Error("Not enough tokens")
	}
	if !tokenValid(token.Children[0], TextNormal, "Hello _world") {
		t.Error("Not valid Text normal")
	}
}

func TestMultipleSpans(t *testing.T) {
	token := parseSpans("Hello **world _of_**")
	if len(token.Children) < 1 {
		t.Error("Not enough tokens")
	}
	fmt.Println(token.Children[1])
	if !tokenValid(token.Children[0], TextNormal, "Hello ") {
		t.Error("Not valid Text normal")
	}
}
*/

func TestSpanSimpleText(t *testing.T) {
	spans := parseSpans("world")

	if len(spans) == 0 {
		t.Error("Not enough spans.")
	}
	if spans[0].ttype != "normal" {
		t.Error("Span should be normal type.")
	}
	if spans[0].value != "world" {
		t.Errorf("Span value should be `world` == `%s`.", spans[0].value)
	}
}

func TestSpanSimpleTextBold(t *testing.T) {
	spans := parseSpans("world *hello*")
	if len(spans) < 2 {
		t.Error("Not enough spans.")
	}
	if spans[0].ttype != "normal" || spans[0].value != "world " {
		t.Error("Span should be normal type of value world.")
	}
	if spans[1].ttype != "italic" {
		t.Error("Span should be italic")
	}
	if spans[2].ttype != "normal" || spans[2].value != "hello" {
		t.Error("Span should be normal with value hello")
	}
}

func TestParseLink(t *testing.T) {
	spans := parseSpans("[example](https://example.com)")
	if len(spans) < 2 {
		t.Errorf("Not enough spans. %d", len(spans))
	}
	span := spans[0]
	if span.ttype != "link-txt" || span.value != "example" {
		t.Error("Not finding alt text of link.")
		return
	}
	span = spans[1]
	if span.ttype != "link-url" || span.value != "https://example.com" {
		t.Error("Not finding url of link.")
		return
	}
}

func TestParseTwoLinks(t *testing.T) {
	spans := parseSpans("[example](https://example.com)[example](https://example.com)")
	if len(spans) < 4 {
		t.Errorf("Not enough spans. %d", len(spans))
	}

	values := []span{
		{"example", "link-txt"},
		{"https://example.com", "link-url"},
		{"example", "link-txt"},
		{"https://example.com", "link-url"},
	}

	for i, span := range spans {
		if span.ttype != values[i].ttype || span.value != values[i].value {
			t.Errorf("Span not correct %s", span)
			return
		}
	}
}

func TestParseTwoLinksWithSpace(t *testing.T) {
	spans := parseSpans("[example](https://example.com) [example](https://example.com)")
	if len(spans) < 5 {
		t.Errorf("Not enough spans. %d", len(spans))
	}

	values := []span{
		{"example", "link-txt"},
		{"https://example.com", "link-url"},
		{" ", "normal"},
		{"example", "link-txt"},
		{"https://example.com", "link-url"},
	}

	for i, span := range spans {
		if span.ttype != values[i].ttype || span.value != values[i].value {
			t.Errorf("Span not correct %s", span)
			return
		}
	}
}

func TestParseTwoImages(t *testing.T) {
	spans := parseSpans("![example](https://example.com)![example](https://example.com)")
	if len(spans) < 4 {
		t.Errorf("Not enough spans. %d", len(spans))
	}

	values := []span{
		{"example", "img-alt"},
		{"https://example.com", "img-src"},
		{"example", "img-alt"},
		{"https://example.com", "img-src"},
	}

	for i, span := range spans {
		if span.ttype != values[i].ttype || span.value != values[i].value {
			t.Errorf("Span not correct %s", span)
			return
		}
	}
}

func TestParseTwoImgsWithSpace(t *testing.T) {
	spans := parseSpans("![example](https://example.com) ![example](https://example.com)")
	if len(spans) < 5 {
		t.Errorf("Not enough spans. %d", len(spans))
	}

	values := []span{
		{"example", "img-alt"},
		{"https://example.com", "img-src"},
		{" ", "normal"},
		{"example", "img-alt"},
		{"https://example.com", "img-src"},
	}

	for i, span := range spans {
		if span.ttype != values[i].ttype || span.value != values[i].value {
			t.Errorf("Span not correct %s", span)
			return
		}
	}
}

func TestSpanParseMultipleSpans(t *testing.T) {
	spans := parseSpans("*world **hello*** ![world](https://example.com/img.jpg) [example](https://example.com)")
	if len(spans) < 10 {
		t.Errorf("Not enough spans. %d", len(spans))
	}

	values := []span{
		{"", "italic"},
		{"world ", "normal"},
		{"", "bold"},
		{"hello", "normal"},
		{"", "endbold"},
		{"", "enditalic"},
		{" ", "normal"},
		{"world", "img-alt"},
		{"https://example.com/img.jpg", "img-src"},
		{" ", "normal"},
		{"example", "link-txt"},
		{"https://example.com", "link-url"},
	}

	if len(values) != len(spans) {
		t.Error("Spans and values not equal.")
		return
	}

	for i, span := range spans {
		if i >= len(values) {
			t.Errorf("Spans not found: %s", span)
			return
		}
		if values[i].ttype != span.ttype || values[i].value != span.value {
			t.Errorf("Spans do not match: %s | %s", span, values[i])
			return
		}
	}
}
