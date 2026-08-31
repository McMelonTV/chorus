package lexer

import (
	"errors"
	"io"

	"github.com/mcmelontv/chorus/internal/token"
)

func isLiteralTokenStart(r rune) bool {
	_, ok := token.LiteralRoot.Children[r]
	return ok
}

func (l *Lexer) lexLiteral() (token.Token, error) {
	start := l.pos

	node := token.LiteralRoot

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

		next := node.Children[br.r]
		if next == nil {
			break
		}

		node = next
		traversed++

		if node.Terminal {
			matchedRunes = traversed
			matchedKind = node.Kind
		}
	}

	if matchedRunes == 0 {
		if err := l.advanceN(traversed); err != nil {
			return token.Token{}, err
		}
		return token.Token{}, &Error{
			Span:    token.Span{Start: start, End: l.pos},
			Message: "invalid token",
		}
	}

	if err := l.advanceN(matchedRunes); err != nil {
		return token.Token{}, err
	}

	return token.NewToken(matchedKind, "", start, l.pos), nil
}
