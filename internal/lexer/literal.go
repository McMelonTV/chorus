package lexer

import (
	"errors"
	"io"
)

func isLiteralTokenStart(r rune) bool {
	_, ok := literalRoot.children[r]
	return ok
}

func (l *Lexer) lexLiteral() (Token, error) {
	start := l.pos

	node := literalRoot

	traversed := 0
	matchedRunes := 0
	var matchedKind TokenKind

	for depth := 0; ; depth++ {
		br, err := l.peekSkip(depth)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return Token{}, err
		}

		next := node.children[br.r]
		if next == nil {
			break
		}

		node = next
		traversed++

		if node.terminal {
			matchedRunes = traversed
			matchedKind = node.kind
		}
	}

	if matchedRunes == 0 {
		if err := l.advanceN(traversed); err != nil {
			return Token{}, err
		}
		return Token{}, &Error{
			Span:    Span{Start: start, End: l.pos},
			Message: "invalid token",
		}
	}

	if err := l.advanceN(matchedRunes); err != nil {
		return Token{}, err
	}

	return newToken(matchedKind, "", start, l.pos), nil
}
