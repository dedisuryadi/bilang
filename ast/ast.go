package ast

import (
	"bytes"
	"strings"

	"github.com/dedisuryadi/bilang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
	Type() token.Type
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Type() token.Type { return token.PROGRAM }
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (id *Identifier) expressionNode()      {}
func (id *Identifier) Type() token.Type     { return id.Token.Type }
func (id *Identifier) TokenLiteral() string { return id.Token.Literal }
func (id *Identifier) String() string       { return id.Value }

type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) expressionNode()      {}
func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) Type() token.Type     { return vs.Token.Type }
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")
	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type KonstStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ks *KonstStatement) statementNode()       {}
func (ks *KonstStatement) Type() token.Type     { return ks.Token.Type }
func (ks *KonstStatement) TokenLiteral() string { return ks.Token.Literal }
func (ks *KonstStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ks.TokenLiteral() + " ")
	out.WriteString(ks.Name.String())
	out.WriteString(" = ")
	if ks.Value != nil {
		out.WriteString(ks.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type PilihStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (ps *PilihStatement) statementNode()       {}
func (ps *PilihStatement) Type() token.Type     { return ps.Token.Type }
func (ps *PilihStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PilihStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.TokenLiteral() + " ")
	if ps.ReturnValue != nil {
		out.WriteString(ps.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) Type() token.Type     { return es.Token.Type }
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) Type() token.Type     { return il.Token.Type }
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) Type() token.Type     { return pe.Token.Type }
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) Type() token.Type     { return ie.Token.Type }
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) Type() token.Type     { return b.Token.Type }
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type JikaExpression struct {
	Token       token.Token // The 'jika' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (je *JikaExpression) expressionNode()      {}
func (je *JikaExpression) Type() token.Type     { return je.Token.Type }
func (je *JikaExpression) TokenLiteral() string { return je.Token.Literal }
func (je *JikaExpression) String() string {
	var out bytes.Buffer
	out.WriteString("jika")
	out.WriteString(je.Condition.String())
	out.WriteString(" ")
	out.WriteString(je.Consequence.String())
	if je.Alternative != nil {
		out.WriteString("atau ")
		out.WriteString(je.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) Type() token.Type     { return bs.Token.Type }
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) Type() token.Type     { return fl.Token.Type }
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

type PilahExpression struct {
	Token      token.Token
	Target     *ExpressionStatement
	Conditions []*ExpressionStatement
	Values     []Expression
}

func (pl *PilahExpression) expressionNode()      {}
func (pl *PilahExpression) Type() token.Type     { return pl.Token.Type }
func (pl *PilahExpression) TokenLiteral() string { return pl.Token.Literal }
func (pl *PilahExpression) String() string {
	var out bytes.Buffer
	out.WriteString(pl.Target.String())
	for i := range pl.Conditions {
		out.WriteString(pl.Conditions[i].String())
		out.WriteString(token.FATARROW)
		out.WriteString(pl.Values[i].String())
	}
	return out.String()
}

type Wildcard struct {
	Token token.Token
}

func (w *Wildcard) expressionNode()      {}
func (w *Wildcard) Type() token.Type     { return w.Token.Type }
func (w *Wildcard) TokenLiteral() string { return w.Token.Literal }
func (w *Wildcard) String() string       { return w.Token.Literal }

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) Type() token.Type     { return ce.Token.Type }
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) Type() token.Type     { return sl.Token.Type }
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) Type() token.Type     { return al.Token.Type }
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) Type() token.Type     { return ie.Token.Type }
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (h *HashLiteral) expressionNode()      {}
func (h *HashLiteral) Type() token.Type     { return h.Token.Type }
func (h *HashLiteral) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for k, v := range h.Pairs {
		pairs = append(pairs, k.String()+":"+v.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type RegExLiteral struct {
	Token token.Token
	Value string
}

func (rel *RegExLiteral) expressionNode()      {}
func (rel *RegExLiteral) Type() token.Type     { return rel.Token.Type }
func (rel *RegExLiteral) TokenLiteral() string { return rel.Token.Literal }
func (rel *RegExLiteral) String() string       { return rel.Value }

type MethodCallExpression struct {
	Token  token.Token
	Object Expression
	Call   Expression
}

func (mc *MethodCallExpression) expressionNode()      {}
func (mc *MethodCallExpression) Type() token.Type     { return mc.Token.Type }
func (mc *MethodCallExpression) TokenLiteral() string { return mc.Token.Literal }
func (mc *MethodCallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(mc.Object.String())
	out.WriteString(".")
	out.WriteString(mc.Call.String())
	return out.String()
}

type Pipe struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (p *Pipe) expressionNode()      {}
func (p *Pipe) Type() token.Type     { return p.Token.Type }
func (p *Pipe) TokenLiteral() string { return p.Token.Literal }
func (p *Pipe) String() string {
	var out bytes.Buffer
	out.WriteString(p.Left.String())
	out.WriteString(" |> ")
	out.WriteString(p.Right.String())
	return out.String()
}

type NihilLiteral struct {
	Token token.Token
}

func (n *NihilLiteral) expressionNode()      {}
func (n *NihilLiteral) Type() token.Type     { return n.Token.Type }
func (n *NihilLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *NihilLiteral) String() string       { return n.Token.Literal }
