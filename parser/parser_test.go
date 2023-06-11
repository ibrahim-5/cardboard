package parser

import (
	"cardboard/lexer"
	"cardboard/parser/ast"
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

	checkParserErrors(t, &p)

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

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.getErrors()
	if len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
		t.FailNow()
	}
}
