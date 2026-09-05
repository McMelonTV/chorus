package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mcmelontv/chorus/internal/lexer"
	"github.com/mcmelontv/chorus/internal/source"
	"github.com/mcmelontv/chorus/internal/token"
)

func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "no path provided\n")
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get working directory: %s\n", err)
		os.Exit(1)
	}

	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to resolve path: %s\n", err)
		os.Exit(1)
	}

	relative, err := filepath.Rel(cwd, path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get relative path: %s\n", err)
		os.Exit(1)
	}

	if relative == ".." ||
		strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		_, _ = fmt.Fprintf(os.Stderr, "file is outside project directory\n")
		os.Exit(1)
	}

	name := filepath.ToSlash(relative)

	data, err := os.ReadFile(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err)
		os.Exit(1)
	}

	fs := source.NewFileSet()

	file, ok := fs.AddFile(name, data)
	if !ok {
		_, _ = fmt.Fprintf(os.Stderr, "failed to add file %s\n", name)
		os.Exit(1)
	}

	l := lexer.New(file)
	for {
		t, err := l.Next()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to lex next token: %s\n", err)
			os.Exit(1)
		}

		fmt.Print(t.String(), " ")

		if t.Kind == token.TOKEN_EOF {
			break
		}
	}
	fmt.Println()
}
