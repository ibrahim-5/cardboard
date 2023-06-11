package ast

import (
	"cardboard/lexer/token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Every Cardboard Program is a list of statements,
// therefore the AST Root Node is the list of statements of the program
type Program struct {
	Statements []Statement
}

// Identifiers are Expressions.
type Identifier struct {
	NodeToken token.Token
	Value     string
}

func (ident *Identifier) expressionNode()      {}
func (ident *Identifier) TokenLiteral() string { return ident.Value }

// 'put' statement
// put <identifier> = <expression>
type PutStatement struct {
	NodeToken      token.Token
	NodeIdentifier Identifier
	NodeExpression Expression
}

func (p *PutStatement) statementNode()       {}
func (p *PutStatement) TokenLiteral() string { return p.NodeToken.TokenLiteral }

// 'unbox' statement. Basically return statement!
// unbox <expression>;
type UnboxStatement struct {
	NodeToken      token.Token
	NodeExpression Expression
}

func (p *UnboxStatement) statementNode()       {}
func (p *UnboxStatement) TokenLiteral() string { return p.NodeToken.TokenLiteral }
