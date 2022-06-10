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
	token := parseSpans("world")

	if len(token.Children) == 0 {
		t.Error("Not enough spans.")
	}
	if token.Children[0].Ttype != Text {
		t.Error("Span should be normal type.")
	}
	if token.Children[0].Value != "world" {
		t.Errorf("Span value should be `world` == `%s`.", token.Children[0].Value)
	}
}

func TestSpanSimpleTextBold(t *testing.T) {
	token := parseSpans("world *hello*")
	if len(token.Children) < 2 {
		t.Error("Not enough token.")
		return
	}
	if token.Children[0].Ttype != Text || token.Children[0].Value != "world " {
		t.Error("Span should be normal type of value world.")
		return
	}
	if token.Children[1].Ttype != Italic {
		t.Error("Span should be italic")
		return
	}
	if token.Children[2].Ttype != Text || token.Children[2].Value != "hello" {
		t.Error("Span should be normal with value hello")
		return
	}
}

func TestParseLink(t *testing.T) {
	token := parseSpans("[example](https://example.com)")
	if len(token.Children) < 1 {
		t.Errorf("Not enough token. %d", len(token.Children))
		return
	}
	span := token.Children[0]
	if span.Ttype != Link || span.Value != "example" {
		t.Error("Not finding value of link.")
		return
	}

	if span.Attrs["url"] != "https://example.com" {
		t.Error("url not matching.")
		return
	}
}

func TestParseTwoLinks(t *testing.T) {
	token := parseSpans("[example](https://example.com)[example](https://example.com)")
	if len(token.Children) < 2 {
		t.Errorf("Not enough token. %d", len(token.Children))
		return
	}

	for _, span := range token.Children {
		if span.Ttype != Link || span.Value != "example" || span.Attrs["url"] != "https://example.com" {
			t.Errorf("Span not correct %s", span)
			return
		}
	}
}

func TestParseTwoLinksWithSpace(t *testing.T) {
	token := parseSpans("[example1](https://example.com) [example2](https://example.com)")
	if len(token.Children) < 3 {
		t.Errorf("Not enough token. %d", len(token.Children))
		return
	}
	if token.Children[0].Ttype != Link || token.Children[2].Ttype != Link {
		t.Error("Not links.")
		return
	}
	if token.Children[1].Ttype != Text {
		t.Error("Not space between links.")
		return
	}
	if token.Children[0].Value != "example1" || token.Children[2].Value != "example2" {
		t.Error("Not valid text.")
		return
	}

	url := "https://example.com"
	if token.Children[0].Attrs["url"] != url || token.Children[2].Attrs["url"] != url {
		t.Error("Not valid url.")
		return
	}
}

func TestParseTwoImages(t *testing.T) {
	token := parseSpans("![example](https://example.com)![example](https://example.com)")
	if len(token.Children) < 2 {
		t.Errorf("Not enough token. %d", len(token.Children))
		return
	}

	img := newToken(Image, "")
	img.Attrs["alt"] = "example"
	img.Attrs["src"] = "https://example.com"
	imgs := []*Token{img, img}

	for i, span := range token.Children {
		if span.Ttype != Image {
			t.Errorf("Span not an image %s", span)
			return
		}
		if span.Attrs["alt"] != imgs[i].Attrs["alt"] {
			t.Errorf("Span url not correct %s", span)
			return
		}
		if span.Attrs["src"] != imgs[i].Attrs["src"] {
			t.Errorf("Span src not correct %s", span)
			return
		}
	}
}

func TestParseTwoImgsWithSpace(t *testing.T) {
	token := parseSpans("![example](https://example.com) ![example](https://example.com)")
	if len(token.Children) < 3 {
		t.Errorf("Not enough token.Children. %d", len(token.Children))
		return
	}

	img := newToken(Image, "")
	img.Attrs["alt"] = "example"
	img.Attrs["src"] = "https://example.com"
	tags := []*Token{img, newToken(Text, " "), img}

	for i, span := range token.Children {
		if span.Ttype != tags[i].Ttype {
			t.Errorf("Span not %s", span)
			return
		}
		if span.Attrs["alt"] != tags[i].Attrs["alt"] {
			t.Errorf("Span url not correct %s", span)
			return
		}
		if span.Attrs["src"] != tags[i].Attrs["src"] {
			t.Errorf("Span src not correct %s", span)
			return
		}
	}
}

func TestSpanParseMultipleSpans(t *testing.T) {
	token := parseSpans("*world **hello*** ![world](https://example.com/img.jpg) [example](https://example.com)")
	if len(token.Children) < 10 {
		t.Errorf("Not enough token. %d", len(token.Children))
		return
	}

	img := newToken(Image, "")
	img.Attrs["alt"] = "world"
	img.Attrs["src"] = "https://example.com/img.jpg"
	link := newToken(Link, "example")
	link.Attrs["url"] = "https://example.com"
	tags := []*Token{
		newToken(Italic, ""),
		newToken(Text, "world "),
		newToken(Bold, ""),
		newToken(Text, "hello"),
		newToken(EndBold, ""),
		newToken(EndItalic, ""),
		newToken(Text, " "),
		img,
		newToken(Text, " "),
		link,
	}

	if len(tags) != len(token.Children) {
		t.Error("Spans and values not equal.")
		return
	}

	for i, span := range token.Children {
		if span.Ttype != tags[i].Ttype {
			t.Errorf("Span not %s", span)
			return
		}
		if span.Value != tags[i].Value {
			t.Errorf("Span value not correct %s", span)
			return
		}
		if span.Attrs["alt"] != tags[i].Attrs["alt"] {
			t.Errorf("Span url not correct %s", span)
			return
		}
		if span.Attrs["src"] != tags[i].Attrs["src"] {
			t.Errorf("Span src not correct %s", span)
			return
		}
	}
}
