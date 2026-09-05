package lexer

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mcmelontv/chorus/internal/token"
)

func (l *Lexer) lexString() (token.Token, error) {
	start := l.offset

	br, err := l.advance()
	if err != nil {
		return token.Token{}, err
	}

	if br.r != '"' {
		return l.tokenUnexpected(start)
	}

	var b strings.Builder

	for {
		runeStart := l.offset

		br, err = l.advance()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return token.Token{}, l.errorEndCurrent(start, "unterminated string")
			}
			return token.Token{}, err
		}

		switch br.r {
		case '"':
			return l.tokenEndCurrent(token.TOKEN_VALUE_STRING, b.String(), start)
		case '\\':
			r, err := l.readEscape(runeStart, '"')
			if err != nil {
				return token.Token{}, err
			}

			b.WriteRune(r)
		case '\n', '\r':
			// TODO: remove this if we wanna allow multiline strings
			return token.Token{}, l.errorEndCurrent(runeStart, "newline in string literal")
		default:
			b.WriteRune(br.r)
		}
	}
}

func (l *Lexer) lexRune() (token.Token, error) {
	start := l.offset

	br, err := l.advance()
	if err != nil {
		return token.Token{}, err
	}

	if br.r != '\'' {
		return l.tokenUnexpected(start)
	}

	runeStart := l.offset

	br, err = l.advance()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return token.Token{}, l.errorEndCurrent(start, "unterminated rune literal")
		}
		return token.Token{}, err
	}

	var r rune

	switch br.r {
	case '\'':
		return token.Token{}, l.errorEndCurrent(start, "empty rune literal")
	case '\\':
		r, err = l.readEscape(runeStart, '\'')
		if err != nil {
			return token.Token{}, err
		}
	case '\n', '\r':
		return token.Token{}, l.errorEndCurrent(runeStart, "newline in rune literal")
	default:
		r = br.r
	}

	br, err = l.advance()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return token.Token{}, l.errorEndCurrent(start, "unterminated rune literal")
		}
		return token.Token{}, err
	}

	if br.r != '\'' {
		return token.Token{}, l.errorEndCurrent(start, "rune literal must contain exactly one rune")
	}

	return l.tokenEndCurrent(token.TOKEN_VALUE_RUNE, string(r), start)
}

func (l *Lexer) readEscape(start uint32, delimiter rune) (rune, error) {
	br, err := l.advance()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return 0, l.errorEndCurrent(start, "unterminated escape sequence")
		}
		return 0, err
	}

	// TODO: implement \u{...} unicode escape sequence here if needed

	switch br.r {
	case 'n':
		return '\n', nil
	case 'r':
		return '\r', nil
	case 't':
		return '\t', nil
	case '\\':
		return '\\', nil
	case '"', '\'':
		if br.r != delimiter {
			return 0, l.errorEndCurrent(start, fmt.Sprintf("escape sequence \\%c is not valid in %s literal", br.r, literalKindName(delimiter)))
		}

		return br.r, nil
	case '0':
		return '\x00', nil
	case '\n', '\r':
		return 0, l.errorEndCurrent(start, "newline in escape sequence")
	default:
		return 0, l.errorEndCurrent(start, fmt.Sprintf("unknown escape sequence \\%c", br.r))
	}
}

func literalKindName(delimiter rune) string {
	switch delimiter {
	case '"':
		return "string"
	case '\'':
		return "rune"
	default:
		panic("invalid literal delimiter")
	}
}
