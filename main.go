package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mcmelontv/chorus/internal/lexer"
	"github.com/mcmelontv/chorus/internal/source"
	"github.com/mcmelontv/chorus/internal/token"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("no path provided")
		return
	}

	path := filepath.Clean(os.Args[1])

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("failed to open file: %s\n", err)
		return
	}

	fs := source.NewFileSet()

	file, ok := fs.AddFile(path, data)
	if !ok {
		fmt.Printf("failed to add file %s\n", path)
		return
	}

	l := lexer.New(file)
	for {
		t, err := l.Next()
		if err != nil {
			fmt.Printf("failed to lex next token: %s\n", err)
			break
		}

		fmt.Print(t.String(), " ")

		if t.Kind == token.TOKEN_EOF {
			break
		}
	}
	fmt.Println()
}
