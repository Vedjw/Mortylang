package ast

import "morty/token"

type Noder interface {
	TokenLiteral() string
}

type Statement interface {
	Noder
	statementNode()
}

type Expression interface {
	Noder
	expressionNode()
}

type Program struct {
	Statements []Statement // This is an array of structs which can hold any values that implement the Statement interface
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral() // why only call tokenliteral on the first statement?
	} else {
		return ""
	}
} // the first statement sets up the initial conditions for the rest program...Maybe

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string      // Literal
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
