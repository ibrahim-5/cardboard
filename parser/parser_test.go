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
		if !testPutStatement(t, programStatement, v.expectedIdentifierName) {
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

func testPutStatement(t *testing.T, statement ast.Statement, expectedName string) bool {
	if statement.TokenLiteral() != "put" {
		t.Fatalf("Test Failed! Expected Token Literal <put>. Got <%s>",
			statement.TokenLiteral())
		return false
	}

	putStmt, ok := statement.(*ast.PutStatement)

	if !ok {
		// Fail, Got Wrong Type
		t.Errorf("Test Failed! Expected Type <PUT>. Got Type <%T>", putStmt)
		return false
	} else {
		putStmtName := putStmt.NodeIdentifier.NodeToken.TokenLiteral
		if putStmtName != expectedName {
			t.Errorf("Test Failed! Expected Identifier Name <%s>. Got <%s>", expectedName,
				putStmtName)
			return false
		}
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
