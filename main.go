package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mcmelontv/chorus/internal/lexer"
	"github.com/mcmelontv/chorus/internal/token"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("no path provided")
		return
	}

	file, err := os.Open(filepath.Clean(os.Args[1]))
	if err != nil {
		fmt.Printf("failed to open file: %s", err)
		return
	}
	defer file.Close()

	l := lexer.New(file)
	for {
		t, err := l.Next()
		if err != nil {
			fmt.Printf("failed to lex next token: %s", err)
			break
		}

		fmt.Print(t.String(), " ")

		if t.Kind == token.TOKEN_EOF {
			break
		}
	}
	fmt.Println()
}
