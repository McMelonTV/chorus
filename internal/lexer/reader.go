package lexer

import (
	"errors"
	"fmt"
	"io"
)

type bufferedRune struct {
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

func (l *Lexer) advance() (bufferedRune, error) {
	var br bufferedRune

	if len(l.buffer) > 0 {
		br = l.buffer[0]
		l.buffer = l.buffer[1:]
	} else {
		r, size, err := l.reader.ReadRune()
		if err != nil {
			return bufferedRune{}, err
		}

		br = bufferedRune{
			r:    r,
			size: size,
		}
	}

	l.advancePosition(br)

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

func (l *Lexer) advancePosition(br bufferedRune) {
	l.pos.Offset += br.size

	if br.r == '\n' {
		l.pos.Line++
		l.pos.Column = 1
	} else {
		l.pos.Column++
	}
}

func (l *Lexer) peek() (bufferedRune, error) {
	runes, err := l.peekN(1)
	if err != nil {
		return bufferedRune{}, err
	}

	return runes[0], nil
}

func (l *Lexer) peekN(n int) ([]bufferedRune, error) {
	if n < 0 {
		return nil, fmt.Errorf("peekN: negative count %d", n)
	}

	for len(l.buffer) < n {
		r, size, err := l.reader.ReadRune()
		if err != nil {
			return nil, err
		}

		l.buffer = append(l.buffer, bufferedRune{
			r:    r,
			size: size,
		})
	}

	return l.buffer[:n], nil
}

func (l *Lexer) peekSkip(skip int) (bufferedRune, error) {
	runes, err := l.peekNSkip(1, skip)
	if err != nil {
		return bufferedRune{}, err
	}

	return runes[0], nil
}

func (l *Lexer) peekNSkip(n, skipCount int) ([]bufferedRune, error) {
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
