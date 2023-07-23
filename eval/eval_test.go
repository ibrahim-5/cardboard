package eval

import (
	"cardboard/lexer"
	"cardboard/object"
	"cardboard/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10}}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
func TestEvalComplexIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"+5;", 5},
		{"-5;", -5},
		{"-50 + 100 + -50", 0},
		{"20 - 5;", 15},
		{"-15 - 100;", -115},
		{"5 + (5 - 10);", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.CreateLexer(input)
	p := parser.CreateParser(l)
	program := p.ParseCardBoard()
	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("object is not Integer. got=%T", obj)
		return false
	}
	if result.Value != expected {
		t.Fatalf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}
