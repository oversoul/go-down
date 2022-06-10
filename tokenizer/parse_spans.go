package tokenizer

const (
	TextNormal TokenType = "TextNormal"
	TextBold   TokenType = "TextBold"
	TextItalic TokenType = "TextItalic"
)

type SpanType string

const (
	Normal SpanType = "Normal"
	Bold   SpanType = "Bold"
	Italic SpanType = "Italic"
)

type gap struct {
	start int
	end   int
	ttype string
}

type span struct {
	value string
	ttype string
}

func last(arr []gap) *gap {
	return &arr[len(arr)-1]
}

func parseImage(line string, i int) []gap {
	i += 1
	if i >= len(line) || line[i] != '[' {
		return []gap{}
	}

	i++
	alt := gap{i, i, "img-alt"}
	for i < len(line) && line[i] != ']' {
		alt.end = i
		i++
	}

	i++
	if line[i] != '(' {
		return []gap{}
	}

	i++
	src := gap{i, 0, "img-src"}
	for i < len(line) && line[i] != ')' {
		src.end = i
		i++
	}

	return []gap{alt, src}
}

func parseLink(line string, i int) []gap {
	i += 1

	txt := gap{i, 0, "link-txt"}
	for i < len(line) && line[i] != ']' {
		txt.end = i
		i++
	}

	i++
	if line[i] != '(' {
		return []gap{}
	}

	i++
	url := gap{i, 0, "link-url"}
	for i < len(line) && line[i] != ')' {
		url.end = i
		i++
	}

	return []gap{txt, url}
}

func addOrCloseGap(gaps *[]gap, ttype string, i int, count int) int {
	if len(*gaps) > 0 && last(*gaps).ttype == ttype && last(*gaps).end+1 == i {
		last(*gaps).end += 1
		return 0
	}

	if count%2 == 0 {
		*gaps = append(*gaps, gap{i + 1, i, ttype})
		return 1
	}

	*gaps = append(*gaps, gap{i + 1, i, "end" + ttype})
	return 0
}

func parseSpans(line string) []span {
	i := 0
	gaps := []gap{}
	italics := 0
	und_bolds := 0
	star_bolds := 0

	for i < len(line) {
		if line[i] == '!' {
			if img := parseImage(line, i); len(img) > 1 {
				gaps = append(gaps, img...)
				i = img[1].end + 2
				continue
			}
		}
		if line[i] == '[' {
			if link := parseLink(line, i); len(link) > 1 {
				gaps = append(gaps, link...)
				i = link[1].end + 2
				continue
			}
		}
		if line[i] == '*' || line[i] == '_' {
			if i+1 < len(line) && line[i] == '*' && line[i+1] == '*' {
				star_bolds += addOrCloseGap(&gaps, "bold", i, star_bolds)
				i += 2
				continue
			}
			if i+1 < len(line) && line[i] == '_' && line[i+1] == '_' {
				und_bolds += addOrCloseGap(&gaps, "bold", i, und_bolds)
				i += 2
				continue
			}
			italics += addOrCloseGap(&gaps, "italic", i, italics)
			i++
			continue
		}

		if len(gaps) > 0 && last(gaps).ttype == "normal" && last(gaps).end+1 == i {
			last(gaps).end += 1
		} else {
			gaps = append(gaps, gap{i, i, "normal"})
		}

		i++
	}

	spans := []span{}
	for _, gap := range gaps {
		spans = append(
			spans, span{value: line[gap.start : gap.end+1], ttype: gap.ttype},
		)
	}

	return spans
}
