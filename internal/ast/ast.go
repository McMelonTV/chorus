package ast

type Node interface {
	Span()
}

// Expression computes a value
type Expression interface {
	Node
	isExpression()
}

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

type TypeSpec interface {
	Node
	isTypeSpec()
}

type Pattern interface {
	Node
	isPattern()
}

type BindingPattern interface {
	Pattern
	isBindingPattern()
}

type AssignmentPattern interface {
	Pattern
	isAssignmentPattern()
}

// ClassMember, InterfaceMember
