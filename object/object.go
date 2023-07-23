package object

import "strconv"

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Types
const (
	INTEGER_OBJ ObjectType = "INTEGER_OBJ"
	NULL        ObjectType = "NULL"
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
