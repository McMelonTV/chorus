package parser

import (
	"github.com/mcmelontv/chorus/internal/ast"
	"github.com/mcmelontv/chorus/internal/token"
)

type Parser struct {
}

// topdown LL recursive descent pratt parser

func Parse(tokens token.Token) ast.File {
	return ast.File{}
}

// TODO: Stream() with Lexer input later (calls Lexer.Next())

func (p *Parser) parseX() {

}
