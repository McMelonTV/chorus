package lexer

import "github.com/mcmelontv/chorus/internal/token"

func isIdentifierStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func isIdentifierContinue(r rune) bool {
	return isIdentifierStart(r) || isDigit(r)
}

func isValidIdentifier(runes []rune) bool {
	if len(runes) < 1 {
		return false
	}

	if !isIdentifierStart(runes[0]) {
		return false
	}

	for _, r := range runes[1:] {
		if !isIdentifierContinue(r) {
			return false
		}
	}

	return true
}

func (l *Lexer) lexIdentifier() (token.Token, error) {
	start := l.pos

	runes, err := l.readWhile(isIdentifierContinue)
	if err != nil {
		return token.Token{}, err
	}

	if !isValidIdentifier(runes) {
		return token.Token{}, &Error{
			token.Span: token.Span{Start: start, End: l.pos},
			Message:    "invalid identifier",
		}
	}

	return token.newToken(token.TOKEN_IDENTIFIER, string(runes), start, l.pos), nil
}
