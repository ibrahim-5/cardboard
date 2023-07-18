package parser

import (
	"cardboard/lexer"
	"cardboard/parser/ast"
	"fmt"
	"testing"
)

func TestPutStatementParsing1(t *testing.T) {
	input := `
	put x = 5;
	put y = 10;
	put z = 3;
	put a = 1;
	`

	l := lexer.CreateLexer(input)
	p := CreateParser(l)

	programStatements := p.ParseCardBoard()

	checkParserErrors(t, p)

	if len(programStatements.Statements) != 4 {
		t.Fatalf("Test Failed! Expected Statement Count Of 4. Got Count Of <%d>",
			len(programStatements.Statements))
	}

	expectedResult := []struct {
		expectedIdentifierName string
	}{
		{"x"},
		{"y"},
		{"z"},
		{"a"},
	}

	for i, v := range expectedResult {
		programStatement := programStatements.Statements[i]
		if !testStatement(t, programStatement, &v.expectedIdentifierName) {
			return
		}
	}
}

func TestPutStatementParsing2(t *testing.T) {
	input := `
	put + = 5;
	put z = 3;
	put - = 10;
	put a = 1;
	`

	l := lexer.CreateLexer(input)
	p := CreateParser(l)

	p.ParseCardBoard()

	// The Parser Should have 2 Errors in its error log,
	// therefore validate that this is the case!
	if len(p.getErrors()) != 2 {
		t.Fatalf("Test Failed! Expected 2 Errors! Got <%d>.", len(p.getErrors()))
	}
}

func TestUnboxStatementParsing(t *testing.T) {
	input := `
		unbox 5;
		unbox 3;
		unbox 1;

	`
	l := lexer.CreateLexer(input)
	p := CreateParser(l)
	program := p.ParseCardBoard()

	// Expect 3 Unbox Statements
	if len(program.Statements) != 3 {
		t.Fatalf("Test Failed! Expected 3 Unbox Statements. Got <%d>.", len(program.Statements))
	}

	for _, v := range p.ParseCardBoard().Statements {
		if !testStatement(t, v, nil) {
			return
		}
	}
	// TODO: Parse Expression
}

func testStatement(t *testing.T, parsedStatement ast.Statement, expectedName *string) bool {

	// Validating an 'unbox' statement, since no statement name check is involved
	if expectedName == nil {
		v, ok := parsedStatement.(*ast.UnboxStatement)
		if !ok {
			t.Fatalf("Test Failed! Expected Type <Unbox>. Got Type <%T>.", v)
			return false
		}

		if v.TokenLiteral() != "unbox" {
			t.Fatalf("Test Failed! Expected Token Literal <unbox>. Got Literal <%s>.", v.TokenLiteral())
			return false
		}

		return true
	} else
	// Validating an 'put' statement
	{
		v, ok := parsedStatement.(*ast.PutStatement)
		if !ok {
			t.Fatalf("Test Failed! Expected Type <Put>. Got Type <%T>.", v)
			return false
		}

		if v.TokenLiteral() != "put" {
			t.Fatalf("Test Failed! Expected Token Literal <put>. Got Literal <%s>.", v.TokenLiteral())
			return false
		}

		if v.NodeIdentifier.TokenLiteral() != *expectedName {
			t.Fatalf("Test Failed! Expected Identifier Literal <%s>. Got Literal <%s>.", *expectedName, v.TokenLiteral())
			return false
		}

		return true
	}
}

// This test will fail for now.
// func TestStringImplementations(t *testing.T) {
// 	input := "put number = 2002;"
// 	parser := CreateParser(lexer.CreateLexer(input))
// 	program := parser.ParseCardBoard()
// 	checkParserErrors(t, parser)

// 	if len(program.Statements) != 1 {
// 		t.Fatalf("Test Failed! Expected Program Length Of 1. Got Length <%d>", len(program.Statements))
// 	}

// 	programString := program.String()

// 	if programString != "put number = 2002;" {
// 		t.Fatalf("Test Failed! Expected Program String 'put number = 2002;'. Got String <%s>", programString)
// 	}

// }

func TestIdentifierExpression(t *testing.T) {
	input := "hello;"
	p := CreateParser(lexer.CreateLexer(input))
	program := p.ParseCardBoard()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Test Failed! Expected Program Length Of 1. Got Length <%d>", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Test Failed! Statement is not *ast.ExpressionStatement. Got <%T>", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("Test Failed! Statement is not *ast.Identifier. Got <%T>", stmt.Expression)

	}

	if ident.Value != "hello" {
		t.Fatalf("Test Failed! Identifier Value not %s. got=%s", "hello", ident.Value)
	}
	if ident.TokenLiteral() != "hello" {
		t.Fatalf("Test Failed! Identifier TokenLiteral not %s. got=%s", "hello",
			ident.TokenLiteral())
	}
}

func TestIntegerLiteral(t *testing.T) {
	input := "100;"
	p := CreateParser(lexer.CreateLexer(input))
	program := p.ParseCardBoard()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Test Failed! Expected Program Length Of 1. Got Length <%d>", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Test Failed! Statement is not *ast.ExpressionStatement. Got <%T>", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Test Failed! Statement is not *ast.IntegerLiteral. Got <%T>", stmt.Expression)

	}

	if ident.Value != 100 {
		t.Fatalf("Test Failed! Identifier Value not %d. got=%d", 100, ident.Value)
	}
}

func TestPrefixExpression(t *testing.T) {
	prefixExpression := []struct {
		input    string
		operator string
		value    int64
	}{
		{"-10;", "-", 10},
		{"-5;", "-", 5},
	}

	for _, tc := range prefixExpression {
		p := CreateParser(lexer.CreateLexer(tc.input))
		program := p.ParseCardBoard()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Log(program.Statements)
			t.Fatalf("Test Failed! Expected Program Length Of 1. Got Length <%d>", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Test Failed! Statement is not *ast.ExpressionStatement. Got <%T>", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("Test Failed! Statement is not *ast.PrefixExpression. Got <%T>", stmt.Expression)

		}

		if exp.Operator != tc.operator {
			t.Fatalf("Test Failed! Operator not equal to %s . Got <%s>", tc.operator, exp.Operator)
		}

		if !testIntegerLiterals(t, tc.value, exp.Right) {
			return
		}

	}
}

func testIntegerLiterals(t *testing.T, tcVal int64, exp ast.Expression) bool {
	intexp, ok := exp.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("Error. Didn't Get *ast.IntegerLiteral. Got <%T>\n", exp)
		return false
	}

	if intexp.Value != tcVal {
		t.Errorf("Value not %d. Got <%d>\n", tcVal, intexp.Value)
		return false
	}

	if intexp.TokenLiteral() != fmt.Sprintf("%d", tcVal) {
		t.Errorf("TokenLiteral not %d. Got <%s>\n", tcVal,
			intexp.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.getErrors()
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
		t.FailNow()
	}
}
