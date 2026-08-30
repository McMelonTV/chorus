package lexer

import "fmt"

type FilePos struct {
	Line, Column int
}

type Span struct {
	Start, End FilePos
}

type Token struct {
	Kind  TokenKind
	Value string
	Span  Span
}

func (t *Token) String() string {
	return fmt.Sprintf("%s(%s)", t.Kind, t.Value)
}

type TokenKind int

//go:generate stringer -type=TokenKind -trimprefix=TOKEN_
const (
	TOKEN_UNEXPECTED TokenKind = iota
	TOKEN_EOF

	// Arithmetic
	TOKEN_PLUS     // +
	TOKEN_MINUS    // -
	TOKEN_MULTIPLY // *
	TOKEN_DIVIDE   // /
	TOKEN_MODULO   // %

	// Assignment
	TOKEN_DECLARE_ASSIGN  // :=
	TOKEN_ASSIGN          // =
	TOKEN_PLUS_ASSIGN     // +=
	TOKEN_MINUS_ASSIGN    // -=
	TOKEN_MULTIPLY_ASSIGN // *=
	TOKEN_DIVIDE_ASSIGN   // /=
	TOKEN_MODULO_ASSIGN   // %=

	// Bitwise operations
	TOKEN_BITWISE_AND // &
	TOKEN_BITWISE_OR  // |
	TOKEN_BITWISE_XOR // ^
	TOKEN_BITWISE_NOT // ~

	TOKEN_BITWISE_SHIFT_LEFT  // <<
	TOKEN_BITWISE_SHIFT_RIGHT // >>

	TOKEN_BITWISE_AND_ASSIGN // &=
	TOKEN_BITWISE_OR_ASSIGN  // |=
	TOKEN_BITWISE_XOR_ASSIGN // ^=
	TOKEN_BITWISE_NOT_ASSIGN // ~=

	// Comparison
	TOKEN_EQUALS         // ==
	TOKEN_NOT_EQUALS     // !=
	TOKEN_GREATER_EQUALS // >=
	TOKEN_LESS_EQUALS    // <=

	TOKEN_CHEVRON_LEFT  // <
	TOKEN_CHEVRON_RIGHT // >

	// Delimeters
	TOKEN_PAREN_OPEN    // (
	TOKEN_PAREN_CLOSE   // )
	TOKEN_BRACE_OPEN    // {
	TOKEN_BRACE_CLOSE   // }
	TOKEN_BRACKET_OPEN  // [
	TOKEN_BRACKET_CLOSE // ]
	TOKEN_SEMICOLON     // ;
	TOKEN_COLON         // :
	TOKEN_COMMA         // ,
	TOKEN_DOT           // .
	TOKEN_AT            // @

	// Logic
	TOKEN_LOGICAL_AND // &&
	TOKEN_LOGICAL_OR  // ||

	// Optional
	TOKEN_QUESTION       // ?
	TOKEN_BANG           // !
	TOKEN_NONE_COALESCE  // ??
	TOKEN_NONE_CONDITION // :?

	// Control Flow
	TOKEN_IF     // if
	TOKEN_ELSE   // else
	TOKEN_FOR    // for
	TOKEN_RETURN // return

	// Literals
	TOKEN_TRUE  // true
	TOKEN_FALSE // false
	TOKEN_NONE  // none

	// Qualifiers
	TOKEN_CONST  // const
	TOKEN_STABLE // stable
	TOKEN_EMBED  // embed
	TOKEN_STATIC // static

	// Reference
	TOKEN_POUND // #
	TOKEN_SELF  // self
	TOKEN_THIS  // this
	TOKEN_NEW   // new

	// Keywords
	TOKEN_CLASS     // class
	TOKEN_INTERFACE // interface
	TOKEN_CONTAINS  // contains
	TOKEN_SATISFIES // satisfies

	// Non-static
	TOKEN_IDENTIFIER    // Token.Value
	TOKEN_VALUE_STRING  // Token.Value
	TOKEN_VALUE_NUMERIC // Token.Value

	// Type tokens?
	TOKEN_MAP // map
)
