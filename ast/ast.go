package ast

import (
	"bytes"
	"morty/token"
)

type Noder interface {
	TokenLiteral() string
	toString() string
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

func (p *Program) TokenLiteral() string { // why only call tokenliteral on the first statement?
	if len(p.Statements) > 0 {
		// the implementation of this method is determined dynamically,
		// depending upon the TokenType of the first statement
		return p.Statements[0].TokenLiteral()
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

// the ident is treated as an expression so it can be evaluated as one
// Ex. a + 5
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (p *Program) toString() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.toString())
	}
	return out.String()
}

func (ls *LetStatement) toString() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.toString())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.toString())
	}
	out.WriteString(";")

	return out.String()
}

func (rs *ReturnStatement) toString() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.toString())
	}
	out.WriteString(";")

	return out.String()
}

func (es *ExpressionStatement) toString() string {
	if es.Expression != nil {
		return es.Expression.toString()
	}
	return ""
}

func (i *Identifier) toString() string {
	return i.Value
}
