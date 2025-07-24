package ast

import (
	"bytes"
	"morty/token"
	"strings"
)

type Noder interface {
	TokenLiteral() string
	ToString() string
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
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) ToString() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.ToString())
	}

	return out.String()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

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

func (ls *LetStatement) ToString() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.ToString())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.ToString())
	}
	out.WriteString(";")

	return out.String()
}

func (rs *ReturnStatement) ToString() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.ToString())
	}
	out.WriteString(";")

	return out.String()
}

func (es *ExpressionStatement) ToString() string {
	if es.Expression != nil {
		return es.Expression.ToString()
	}
	return ""
}

func (i *Identifier) ToString() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) ToString() string     { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.ToString())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.ToString())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.ToString())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) ToString() string     { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Concequence *BlockStatement
	Alternative *BlockStatement
}

func (ife *IfExpression) expressionNode()      {}
func (ife *IfExpression) TokenLiteral() string { return ife.Token.Literal }
func (ife *IfExpression) ToString() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ife.Condition.ToString())
	out.WriteString(" ")
	out.WriteString(ife.Concequence.ToString())

	if ife.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ife.Alternative.ToString())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) ToString() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.ToString())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // fn token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) ToString() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.ToString())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.ToString())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // ( token after the identifier
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) ToString() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.ToString())
	}

	out.WriteString(ce.Function.ToString())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (str *StringLiteral) expressionNode()      {}
func (str *StringLiteral) TokenLiteral() string { return str.Token.Literal }
func (str *StringLiteral) ToString() string     { return str.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (arr *ArrayLiteral) expressionNode()      {}
func (arr *ArrayLiteral) TokenLiteral() string { return arr.Token.Literal }
func (arr *ArrayLiteral) ToString() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range arr.Elements {
		elements = append(elements, el.ToString())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
