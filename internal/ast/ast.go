package ast

import "github.com/mcmelontv/chorus/internal/source"

type Node interface {
	Span() source.Span
}

type node struct {
	SourceSpan source.Span
}

func (n *node) Span() source.Span {
	return n.SourceSpan
}

// Expression produces a value
type Expression interface {
	Node
	isExpression()
}

// Item is a node that may appear in a block
type Item interface {
	Node
	isItem()
}

// Statement performs an action
type Statement interface {
	Item
	isStatement()
}

// Declaration introduces a name
type Declaration interface {
	Item
	isDeclaration()
}

// TopLevelDeclaration is a declaration that may appear at the file scope
type TopLevelDeclaration interface {
	Declaration
	isTopLevelDeclaration()
}

// TypeSpec represents a type as written in source code
type TypeSpec interface {
	Node
	isTypeSpec()
}

// Pattern represents destructuring or targeting syntax
type Pattern interface {
	Node
	isPattern()
}

// BindingPattern represents a pattern that introduces new names
type BindingPattern interface {
	Pattern
	isBindingPattern()
}

// AssignmentPattern represents a pattern that refers to existing targets
type AssignmentPattern interface {
	Pattern
	isAssignmentPattern()
}

// ClassMember is a declaration that may appear directly in a class body
type ClassMember interface {
	Declaration
	isClassMember()
}

// InterfaceMember is a declaration that may appear directly in an interface body
type InterfaceMember interface {
	Declaration
	isInterfaceMember()
}

type File struct {
	Source *source.File

	// TODO:
	// module/package
	// imports

	Declaration TopLevelDeclaration
}

// ExpressionStatement evaluates an expression without using its result
type ExpressionStatement struct {
	node

	Expression Expression
}

func (*ExpressionStatement) isItem()      {}
func (*ExpressionStatement) isStatement() {}

// CallExpression invokes a callable expression with a list of arguments
type CallExpression struct {
	node

	Callee    Expression
	Arguments []Expression
}

func (*CallExpression) isExpression() {}

// ReturnStatement returns control to the caller, optionally with a value
type ReturnStatement struct {
	node

	Value Expression // nil for `return;`
}

func (*ReturnStatement) isItem()      {}
func (*ReturnStatement) isStatement() {}
