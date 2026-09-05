package lexer

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

type decodedRune struct {
	r    rune
	size int
}

// discardWhile advances the reader while predicate is true and returns amount of runes discarded
func (l *Lexer) discardWhile(predicate func(rune) bool) (int, error) {
	n := 0

	for {
		br, err := l.peek()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return n, nil
			}
			return n, err
		}

		if !predicate(br.r) {
			return n, nil
		}

		if _, err := l.advance(); err != nil {
			return n, err
		}

		n++
	}
}

// readWhile advances the reader while predicate is true and returns the consumed runes
func (l *Lexer) readWhile(predicate func(rune) bool) ([]rune, error) {
	runes := []rune{}

	for {
		br, err := l.peek()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return runes, nil
			}
			return nil, err
		}

		if !predicate(br.r) {
			return runes, nil
		}

		br, err = l.advance()
		if err != nil {
			return nil, err
		}

		runes = append(runes, br.r)
	}
}

func (l *Lexer) advance() (decodedRune, error) {
	br, err := l.peek()
	if err != nil {
		return decodedRune{}, err
	}

	l.offset += uint32(br.size)
	return br, nil
}

func (l *Lexer) advanceN(n int) error {
	for range n {
		if _, err := l.advance(); err != nil {
			return err
		}
	}

	return nil
}

func (l *Lexer) peek() (decodedRune, error) {
	if uint64(l.offset) >= uint64(len(l.file.Data)) {
		return decodedRune{}, io.EOF
	}

	r, size := utf8.DecodeRune(l.file.Data[l.offset:])
	return decodedRune{r: r, size: size}, nil
}

func (l *Lexer) peekN(n int) ([]decodedRune, error) {
	if n < 0 {
		return nil, fmt.Errorf("peekN: negative count %d", n)
	}

	offset := l.offset
	runes := make([]decodedRune, 0, n)

	for range n {
		if uint64(offset) >= uint64(len(l.file.Data)) {
			return nil, io.EOF
		}

		r, size := utf8.DecodeRune(l.file.Data[offset:])
		runes = append(runes, decodedRune{r: r, size: size})
		offset += uint32(size)
	}

	return runes, nil
}

func (l *Lexer) peekSkip(skip int) (decodedRune, error) {
	runes, err := l.peekNSkip(1, skip)
	if err != nil {
		return decodedRune{}, err
	}

	return runes[0], nil
}

func (l *Lexer) peekNSkip(n, skipCount int) ([]decodedRune, error) {
	if n < 0 {
		return nil, fmt.Errorf("peekNSkip: negative count %d", n)
	}

	if skipCount < 0 {
		return nil, fmt.Errorf("peekNSkip: negative skip count %d", skipCount)
	}

	runes, err := l.peekN(skipCount + n)
	if err != nil {
		return nil, err
	}

	return runes[skipCount : skipCount+n], nil
}
