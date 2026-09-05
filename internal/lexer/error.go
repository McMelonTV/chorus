package lexer

import (
	"fmt"

	"github.com/mcmelontv/chorus/internal/source"
)

type Error struct {
	File    *source.File
	Span    source.Span
	Message string
}

func (e *Error) Error() string {
	p, ok := e.File.FilePos(e.Span.Start)
	if !ok {
		return "lexer/Error.Error(): invalid FilePos"
	}

	return fmt.Sprintf(
		"%d:%d: %s",
		p.Line,
		p.Column,
		e.Message,
	)
}
