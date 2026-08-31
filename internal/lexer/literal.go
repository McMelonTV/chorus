package lexer

import (
	"errors"
	"io"

	"github.com/mcmelontv/chorus/internal/token"
)

func isLiteralTokenStart(r rune) bool {
	_, ok := token.literalRoot.children[r]
	return ok
}

func (l *Lexer) lexLiteral() (token.Token, error) {
	start := l.pos

	node := token.literalRoot

	traversed := 0
	matchedRunes := 0
	var matchedKind token.TokenKind

	for depth := 0; ; depth++ {
		br, err := l.peekSkip(depth)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return token.Token{}, err
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
			return token.Token{}, err
		}
		return token.Token{}, &Error{
			token.Span: token.Span{Start: start, End: l.pos},
			Message:    "invalid token",
		}
	}

	if err := l.advanceN(matchedRunes); err != nil {
		return token.Token{}, err
	}

	return token.newToken(matchedKind, "", start, l.pos), nil
}
