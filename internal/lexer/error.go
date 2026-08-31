package lexer

import "fmt"

type Error struct {
	Span    Span
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"%d:%d: %s",
		e.Span.Start.Line,
		e.Span.Start.Column,
		e.Message,
	)
}
