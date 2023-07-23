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
		evaluated := testEval(tt.input, t)
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
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalUnboxStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"unbox 10 + 2; 100;", 12},
		{"5 + (5 - 10); unbox (90 + 10); -15 - 100;", 100},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalPutStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"put a = 5; a;", 5},
		{"put a = 5 + 5; a;", 10},
		{"put a = 5; put b = a; b;", 5},
		{"put a = 5; put b = a; put c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string, t *testing.T) object.Object {
	l := lexer.CreateLexer(input)
	p := parser.CreateParser(l)
	program := p.ParseCardBoard()
	checkParserErrors(t, p)
	env := object.CreateEnvironment()
	return Eval(program, env)
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

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errs := p.GetErrors()
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
		t.FailNow()
	}
}
