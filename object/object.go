package object

import (
	"bytes"
	"cardboard/parser/ast"
	"strconv"
	"strings"
)

// Objects
type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Types
const (
	INTEGER_OBJ  ObjectType = "INTEGER_OBJ"
	UNBOX_OBJ    ObjectType = "UNBOX_OBJ"
	NULL         ObjectType = "NULL"
	FUNCTION_OBJ ObjectType = "FUNCTION_OBJ"
)

// Integer
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return strconv.Itoa(int(i.Value)) }

// Null
type Null struct{}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "null" }

// Unbox
type Unbox struct {
	Value Object
}

func (un *Unbox) Type() ObjectType { return UNBOX_OBJ }
func (un *Unbox) Inspect() string  { return un.Value.Inspect() }

// Function (box)
type Box struct {
	Env           *Environment
	ParameterList []*ast.Identifier
	Body          *ast.BlockStatement
}

func (f *Box) Type() ObjectType { return FUNCTION_OBJ }
func (f *Box) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.ParameterList {
		params = append(params, p.String())
	}
	out.WriteString("box")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}
