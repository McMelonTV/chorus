package lexer

import (
	"fmt"

	"github.com/mcmelontv/chorus/internal/token"
)

type Error struct {
	Span    token.Span
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
