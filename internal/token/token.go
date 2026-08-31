package token

import (
	"fmt"
)

type literalNode struct {
	Children map[rune]*literalNode

	Kind     TokenKind
	Terminal bool
}

var LiteralRoot = &literalNode{
	Children: make(map[rune]*literalNode),
}

var (
	tokenLiterals map[TokenKind]string
	literalTokens map[string]TokenKind
)

func init() {
	tokenLiterals = make(map[TokenKind]string)
	literalTokens = make(map[string]TokenKind)

	for _, token := range literalTokenList {
		tokenLiterals[token.kind] = token.literal
		literalTokens[token.literal] = token.kind
	}

	for _, lt := range literalTokenList {
		node := LiteralRoot

		for _, r := range lt.literal {
			child := node.Children[r]

			if child == nil {
				child = &literalNode{
					Children: make(map[rune]*literalNode),
				}
				node.Children[r] = child
			}

			node = child
		}

		if node.Terminal {
			panic("duplicate literal token")
		}

		node.Kind = lt.kind
		node.Terminal = true
	}
}

type FilePos struct {
	Offset, Line, Column int
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

func UnexpectedToken(start, end FilePos) Token {
	return NewTokenWithoutValue(TOKEN_UNEXPECTED, start, end)
}

func EofToken(pos FilePos) Token {
	return NewTokenWithoutValue(TOKEN_EOF, pos, pos)
}

func NewTokenWithoutValue(kind TokenKind, start, end FilePos) Token {
	return Token{
		Kind:  kind,
		Value: "",
		Span: Span{
			Start: start,
			End:   end,
		},
	}
}

func NewToken(kind TokenKind, value string, start, end FilePos) Token {
	return Token{
		Kind:  kind,
		Value: value,
		Span: Span{
			Start: start,
			End:   end,
		},
	}
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
	TOKEN_ARROW_RIGHT   // ->
	TOKEN_ARROW_LEFT    // <-
	TOKEN_ARROW_BI      // <->

	// Logic
	TOKEN_LOGICAL_AND // &&
	TOKEN_LOGICAL_OR  // ||

	// Optional
	TOKEN_QUESTION       // ?
	TOKEN_BANG           // !
	TOKEN_NONE_COALESCE  // ??
	TOKEN_NONE_CONDITION // :?

	// Reference
	TOKEN_POUND // #

	// Non-static
	TOKEN_IDENTIFIER    // Token.Value
	TOKEN_VALUE_STRING  // Token.Value
	TOKEN_VALUE_RUNE    // Token.Value
	TOKEN_VALUE_NUMERIC // Token.Value
)

type literalToken struct {
	kind    TokenKind
	literal string
}

var literalTokenList = []literalToken{
	// Arithmetic
	{TOKEN_PLUS, "+"},
	{TOKEN_MINUS, "-"},
	{TOKEN_MULTIPLY, "*"},
	{TOKEN_DIVIDE, "/"},
	{TOKEN_MODULO, "%"},

	// Assignment
	{TOKEN_DECLARE_ASSIGN, ":="},
	{TOKEN_ASSIGN, "="},
	{TOKEN_PLUS_ASSIGN, "+="},
	{TOKEN_MINUS_ASSIGN, "-="},
	{TOKEN_MULTIPLY_ASSIGN, "*="},
	{TOKEN_DIVIDE_ASSIGN, "/="},
	{TOKEN_MODULO_ASSIGN, "%="},

	// Bitwise operations
	{TOKEN_BITWISE_AND, "&"},
	{TOKEN_BITWISE_OR, "|"},
	{TOKEN_BITWISE_XOR, "^"},
	{TOKEN_BITWISE_NOT, "~"},

	{TOKEN_BITWISE_SHIFT_LEFT, "<<"},
	{TOKEN_BITWISE_SHIFT_RIGHT, ">>"},

	{TOKEN_BITWISE_AND_ASSIGN, "&="},
	{TOKEN_BITWISE_OR_ASSIGN, "|="},
	{TOKEN_BITWISE_XOR_ASSIGN, "^="},

	// Comparison
	{TOKEN_EQUALS, "=="},
	{TOKEN_NOT_EQUALS, "!="},
	{TOKEN_GREATER_EQUALS, ">="},
	{TOKEN_LESS_EQUALS, "<="},

	{TOKEN_CHEVRON_LEFT, "<"},
	{TOKEN_CHEVRON_RIGHT, ">"},

	// Delimiters
	{TOKEN_PAREN_OPEN, "("},
	{TOKEN_PAREN_CLOSE, ")"},
	{TOKEN_BRACE_OPEN, "{"},
	{TOKEN_BRACE_CLOSE, "}"},
	{TOKEN_BRACKET_OPEN, "["},
	{TOKEN_BRACKET_CLOSE, "]"},
	{TOKEN_SEMICOLON, ";"},
	{TOKEN_COLON, ":"},
	{TOKEN_COMMA, ","},
	{TOKEN_DOT, "."},
	{TOKEN_AT, "@"},
	{TOKEN_ARROW_RIGHT, "->"},
	{TOKEN_ARROW_LEFT, "<-"},
	{TOKEN_ARROW_BI, "<->"},

	// Logic
	{TOKEN_LOGICAL_AND, "&&"},
	{TOKEN_LOGICAL_OR, "||"},

	// Optional
	{TOKEN_QUESTION, "?"},
	{TOKEN_BANG, "!"},
	{TOKEN_NONE_COALESCE, "??"},
	{TOKEN_NONE_CONDITION, ":?"},

	// Reference
	{TOKEN_POUND, "#"},
}
