package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/mcmelontv/chorus/internal/token"
)

type Lexer struct {
	reader *bufio.Reader
	buffer []bufferedRune
	pos    token.FilePos
}

func New(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
		pos: token.FilePos{
			Offset: 0,
			Line:   1,
			Column: 1,
		},
	}
}

func Lex(reader io.Reader) ([]token.Token, error) {
	return New(reader).Lex()
}

func (l *Lexer) Lex() ([]token.Token, error) {
	var tokens []token.Token

	for {
		t, err := l.Next()
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, t)

		if t.Kind == token.TOKEN_EOF {
			return tokens, nil
		}
	}
}

func (l *Lexer) Next() (token.Token, error) {
	for {
		whitespace, err := l.consumeWhitespace()
		if err != nil {
			return token.Token{}, err
		}

		comment, err := l.consumeComment()
		if err != nil {
			return token.Token{}, err
		}

		if !whitespace && !comment {
			break
		}
	}

	lk, err := l.classifyNext()
	if err != nil {
		return token.Token{}, err
	}

	switch lk {
	case lexKindEOF:
		return token.EofToken(l.pos), nil
	case lexKindIdentifier:
		return l.lexIdentifier()
	case lexKindNumeric:
		return l.lexNumeric()
	case lexKindString:
		return l.lexString()
	case lexKindRune:
		return l.lexRune()
	case lexKindLiteral:
		return l.lexLiteral()
	default:
		start := l.pos

		br, err := l.advance()
		if err != nil {
			return token.Token{}, err
		}

		return token.Token{}, &Error{
			Span:    token.Span{Start: start, End: l.pos},
			Message: fmt.Sprintf("unexpected character %q", br.r),
		}
	}
}

type lexKind int

const (
	lexKindUnexpected lexKind = iota
	lexKindEOF
	lexKindIdentifier
	lexKindNumeric
	lexKindString
	lexKindRune
	lexKindLiteral
)

func (l *Lexer) classifyNext() (lexKind, error) {
	pbr, err := l.peek()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return lexKindEOF, nil
		}

		return lexKindUnexpected, err
	}

	switch {
	case isIdentifierStart(pbr.r):
		return lexKindIdentifier, nil
	case isDigit(pbr.r):
		return lexKindNumeric, nil
	case pbr.r == '.':
		next, err := l.peekSkip(1)
		if err == nil && isDigit(next.r) {
			return lexKindNumeric, nil
		}
		if err != nil && !errors.Is(err, io.EOF) {
			return lexKindUnexpected, err
		}
		return lexKindLiteral, nil
	case pbr.r == '"':
		return lexKindString, nil
	case pbr.r == '\'':
		return lexKindRune, nil
	case isLiteralTokenStart(pbr.r):
		return lexKindLiteral, nil
	}

	return lexKindUnexpected, err
}
