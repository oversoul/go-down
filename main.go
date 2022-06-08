package main

import (
	"encoding/json"
	"fmt"
	"os"
	"oversoul/godown/tokenizer"
)

type Token struct {
	*tokenizer.Token
}

func (root *Token) Render() {
	switch root.Ttype {
	case tokenizer.Heading1:
		fmt.Print("# ")
	case tokenizer.Heading2:
		fmt.Print("## ")
	case tokenizer.Heading3:
		fmt.Print("### ")
	case tokenizer.Heading4:
		fmt.Print("#### ")
	case tokenizer.Heading5:
		fmt.Print("##### ")
	case tokenizer.Heading6:
		fmt.Print("###### ")
	}
	fmt.Println(root.Value)
	for _, item := range root.Children {
		t := Token{item}
		t.Render()
	}
}

func main() {
	// content := "# Welcome to StackEdit!\n\nHi! I'm your first Markdown.\n\n- first item\n- second item\n"
	content, err := os.ReadFile("example.md")
	if err != nil {
		panic(err)
	}
	tokens := tokenizer.Tokenize(string(content))

	// for _, token := range tokens {
	// 	t := Token{token}
	// 	t.Render()
	// }

	data, err := json.Marshal(tokens)
	if err != nil {
		panic(err)
	}

	// _ = data
	fmt.Println(string(data))
}
