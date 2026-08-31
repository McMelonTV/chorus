package parser

import "github.com/mcmelontv/chorus/internal/lexer"

type Parser struct {
	lex *lexer.Lexer
}

// topdown LL recursive descent pratt parser

func New(l *lexer.Lexer) *Parser {
	return &Parser{lex: l}
}

func (p *Parser) parseX() {

}
