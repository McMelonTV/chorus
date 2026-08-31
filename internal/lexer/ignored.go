package lexer

import (
	"errors"
	"io"

	"github.com/mcmelontv/chorus/internal/token"
)

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

func (l *Lexer) consumeWhitespace() (bool, error) {
	n, err := l.discardWhile(isWhitespace)
	return n > 0, err
}

func (l *Lexer) consumeComment() (bool, error) {
	start := l.pos

	p, err := l.peekN(2)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false, nil
		}
		return false, err
	}

	if p[0].r != '/' {
		return false, nil
	}

	switch p[1].r {
	case '/':
		if err := l.advanceN(2); err != nil {
			return false, err
		}

		for {
			br, err := l.advance()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return true, nil
				}
				return true, err
			}

			if br.r == '\n' {
				return true, nil
			}
		}

	case '*':
		if err := l.advanceN(2); err != nil {
			return false, err
		}

		for {
			br, err := l.advance()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return true, &Error{
						Span: token.Span{
							Start: start,
							End:   l.pos,
						},
						Message: "unterminated block comment",
					}
				}
				return true, err
			}

			if br.r != '*' {
				continue
			}

			next, err := l.peek()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return true, &Error{
						Span: token.Span{
							Start: start,
							End:   l.pos,
						},
						Message: "unterminated block comment",
					}
				}
				return true, err
			}

			if next.r != '/' {
				continue
			}

			if _, err := l.advance(); err != nil {
				return true, err
			}

			return true, nil
		}

	default:
		return false, nil
	}
}
