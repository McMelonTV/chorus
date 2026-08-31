package lexer

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

func (l *Lexer) lexIdentifier() (Token, error) {
	start := l.pos

	runes, err := l.readWhile(isIdentifierContinue)
	if err != nil {
		return Token{}, err
	}

	if !isValidIdentifier(runes) {
		return Token{}, &Error{
			Span:    Span{Start: start, End: l.pos},
			Message: "invalid identifier",
		}
	}

	return newToken(TOKEN_IDENTIFIER, string(runes), start, l.pos), nil
}
