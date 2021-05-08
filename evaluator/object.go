package evaluator

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/dedisuryadi/bilang/ast"
)

type Type string

const (
	INTEGER  = "INTEGER"
	BOOLEAN  = "BOOLEAN"
	NULL     = "NULL"
	RETURN   = "RETURN"
	ERROR    = "ERROR"
	FUNCTION = "FUNCTION"
	STRING   = "STRING"
	BUILTIN  = "BUILTIN"
	ARRAY    = "ARRAY"
	HASH     = "HASH"
	benar    = "benar"
	salah    = "salah"
	VOID     = "VOID"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Iterable interface {
	Iter() bool
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type      { return BUILTIN }
func (b *Builtin) Inspect() string { return "builtin function" }

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type      { return INTEGER }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type { return BOOLEAN }
func (b *Boolean) Inspect() string {
	if b.Value {
		return benar
	}
	return salah
}

type Null struct{}

func (n *Null) Type() Type      { return NULL }
func (n *Null) Inspect() string { return "nihil" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type      { return RETURN }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() Type      { return ERROR }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() Type { return FUNCTION }
func (f *Function) Inspect() string {
	var (
		out    bytes.Buffer
		params = []string{}
	)
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Iter() bool      { return true }
func (s *String) Type() Type      { return STRING }
func (s *String) Inspect() string { return s.Value }

type Array struct {
	Elements []Object
}

func (a *Array) Iter() bool { return true }
func (a *Array) Type() Type { return ARRAY }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  Type
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
func (s *String) HashKey() HashKey {
	h := fnv.New64()
	_, _ = h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Iter() bool { return true }
func (h *Hash) Type() Type { return HASH }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Break struct{}

func (b *Break) Inspect() string { return "usai" }
func (b *Break) Type() Type      { return BREAK }

type Continue struct{}

func (c *Continue) Inspect() string { return "lanjut" }
func (c *Continue) Type() Type      { return CONTINUE }
