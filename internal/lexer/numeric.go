package lexer

import (
	"errors"
	"io"

	"github.com/mcmelontv/chorus/internal/token"
)

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (l *Lexer) lexNumeric() (token.Token, error) {
	start := l.offset

	var runes []rune

	left, err := l.readWhile(isDigit)
	if err != nil {
		return token.Token{}, err
	}
	runes = append(runes, left...)

	decimalPoint, err := l.readNumericDecimalPoint()
	if err != nil {
		return token.Token{}, err
	}

	if decimalPoint != nil {
		right, err := l.readWhile(isDigit)
		if err != nil {
			return token.Token{}, err
		}

		runes = append(runes, *decimalPoint)
		runes = append(runes, right...)
	}

	exponent, err := l.readNumericExponent()
	if err != nil {
		return token.Token{}, err
	}
	runes = append(runes, exponent...)

	return l.tokenEndCurrent(token.TOKEN_VALUE_NUMERIC, string(runes), start)
}

func (l *Lexer) readNumericDecimalPoint() (*rune, error) {
	peek, err := l.peek()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, err
	}

	if peek.r != '.' {
		return nil, nil
	}

	next, err := l.peekSkip(1)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, err
	}

	if !isDigit(next.r) {
		return nil, nil
	}

	if _, err = l.advance(); err != nil {
		return nil, err
	}

	return &peek.r, nil
}

func (l *Lexer) readNumericExponent() ([]rune, error) {
	var runes []rune

	peek, err := l.peek()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, err
	}

	if peek.r != 'e' && peek.r != 'E' {
		return nil, nil
	}

	advance := 1

	next, err := l.peekSkip(1)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, err
	}

	nextIsPlusOrMinus := next.r == '+' || next.r == '-'

	if !isDigit(next.r) && !nextIsPlusOrMinus {
		return nil, nil
	}

	if nextIsPlusOrMinus {
		after, err := l.peekSkip(2)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, nil
			}
			return nil, err
		}

		if !isDigit(after.r) {
			return nil, nil
		}

		advance = 2
	}

	for range advance {
		br, err := l.advance()
		if err != nil {
			return nil, err
		}

		runes = append(runes, br.r)
	}

	read, err := l.readWhile(isDigit)
	if err != nil {
		return nil, err
	}

	runes = append(runes, read...)

	return runes, nil
}
