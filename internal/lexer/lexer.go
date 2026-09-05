package lexer

import (
	"errors"
	"fmt"
	"io"

	"github.com/mcmelontv/chorus/internal/source"
	"github.com/mcmelontv/chorus/internal/token"
)

type Lexer struct {
	file   *source.File
	offset uint32
}

func New(f *source.File) *Lexer {
	return &Lexer{
		file:   f,
		offset: 0,
	}
}

func Lex(f *source.File) ([]token.Token, error) {
	return New(f).Lex()
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
		return l.tokenEOF()
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
		start := l.offset

		br, err := l.advance()
		if err != nil {
			return token.Token{}, err
		}

		return token.Token{}, l.error(start, l.offset, fmt.Sprintf("unexpected character %q", br.r))
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

func (l *Lexer) span(start, end uint32) (source.Span, bool) {
	return source.NewSpan(l.file.Pos(start), l.file.Pos(end))
}

func (l *Lexer) currentFilePos() (source.FilePos, bool) {
	return l.file.FilePos(l.currentPos())
}

func (l *Lexer) currentPos() source.Pos {
	return l.file.Pos(l.offset)
}

func (l *Lexer) error(start, end uint32, message string) error {
	span, ok := l.span(start, end)

	if !ok {
		span = source.Span{}
	}

	return &Error{
		File:    l.file,
		Span:    span,
		Message: message,
	}
}

func (l *Lexer) errorEndCurrent(start uint32, message string) error {
	return l.error(start, l.offset, message)
}

func (l *Lexer) token(kind token.TokenKind, value string, start, end uint32) (token.Token, error) {
	span, ok := l.span(start, end)
	if !ok {
		return token.Token{}, fmt.Errorf("lexer/Lexer.token(): invalid span")
	}

	return token.Token{
		Kind:  kind,
		Value: value,
		Span:  span,
	}, nil
}

func (l *Lexer) tokenEmptyValue(kind token.TokenKind, start, end uint32) (token.Token, error) {
	return l.token(kind, "", start, end)
}

func (l *Lexer) tokenEndCurrent(kind token.TokenKind, value string, start uint32) (token.Token, error) {
	return l.token(kind, value, start, l.offset)
}

func (l *Lexer) tokenEmptyValueEndCurrent(kind token.TokenKind, start uint32) (token.Token, error) {
	return l.token(kind, "", start, l.offset)
}

func (l *Lexer) tokenEOF() (token.Token, error) {
	return l.token(token.TOKEN_EOF, "", l.offset, l.offset)
}

func (l *Lexer) tokenUnexpected(offset uint32) (token.Token, error) {
	return l.token(token.TOKEN_UNEXPECTED, "", offset, offset)
}
