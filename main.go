package main

import (
	"fmt"
	"os"
	"path"

	"github.com/mcmelontv/chorus/internal/lexer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("no path provided")
		return
	}

	fileName := os.Args[1]
	filePath := path.Clean(fileName)

	if _, err := os.Stat(filePath); err != nil {
		fmt.Printf("failed to read file info: %s", err)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("failed to open file: %s", err)
		return
	}
	defer file.Close()

	tokens, err := lexer.Lex(file)
	if err != nil {
		fmt.Printf("failed to lex file: %s", err)
		return
	}

	for _, token := range tokens {
		fmt.Print(token.String(), " ")
	}
	fmt.Println()
}
