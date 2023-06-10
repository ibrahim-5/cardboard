package lexer

import (
	"cardboard/lexer/token"
	"testing"
)

func TestLexer1(t *testing.T) {
	input := `=+{}-()`
	expectedResult := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{expectedType: token.EQ, expectedLiteral: "="},
		{expectedType: token.ADD, expectedLiteral: "+"},
		{expectedType: token.LCURLY, expectedLiteral: "{"},
		{expectedType: token.RCURLY, expectedLiteral: "}"},
		{expectedType: token.SUB, expectedLiteral: "-"},
		{expectedType: token.LBRAC, expectedLiteral: "("},
		{expectedType: token.RBRAC, expectedLiteral: ")"},
		{expectedType: token.EOF, expectedLiteral: ""},
	}

	l := CreateLexer(input)

	for _, testToken := range expectedResult {
		lexerToken := l.NextToken()

		// Checking Type

		if (lexerToken.TokenType != testToken.expectedType) ||
			(lexerToken.TokenLiteral != testToken.expectedLiteral) {
			t.Fatalf("Test Failed! Expected Token: <Type: %s, Literal: %s> but Got Token: <Type: %s, Literal: %s>\n",
				testToken.expectedType,
				testToken.expectedLiteral,
				lexerToken.TokenType,
				lexerToken.TokenLiteral)
		}
	}

}

func TestLexer2(t *testing.T) {
	input := `
	
	box add(a, b){
		put y = a + b + 5;
		unbox y;
	}
	
	`
	expectedResult := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{expectedType: token.BOX, expectedLiteral: "box"},
		{expectedType: token.IDENTIFIER, expectedLiteral: "add"},
		{expectedType: token.LBRAC, expectedLiteral: "("},
		{expectedType: token.IDENTIFIER, expectedLiteral: "a"},
		{expectedType: token.COMMA, expectedLiteral: ","},
		{expectedType: token.IDENTIFIER, expectedLiteral: "b"},
		{expectedType: token.RBRAC, expectedLiteral: ")"},
		{expectedType: token.LCURLY, expectedLiteral: "{"},
		{expectedType: token.PUT, expectedLiteral: "put"},
		{expectedType: token.IDENTIFIER, expectedLiteral: "y"},
		{expectedType: token.EQ, expectedLiteral: "="},
		{expectedType: token.IDENTIFIER, expectedLiteral: "a"},
		{expectedType: token.ADD, expectedLiteral: "+"},
		{expectedType: token.IDENTIFIER, expectedLiteral: "b"},
		{expectedType: token.ADD, expectedLiteral: "+"},
		{expectedType: token.INT, expectedLiteral: "5"},
		{expectedType: token.SCOLON, expectedLiteral: ";"},
		{expectedType: token.UNBOX, expectedLiteral: "unbox"},
		{expectedType: token.IDENTIFIER, expectedLiteral: "y"},
		{expectedType: token.SCOLON, expectedLiteral: ";"},
		{expectedType: token.RCURLY, expectedLiteral: "}"},
	}

	l := CreateLexer(input)

	for _, testToken := range expectedResult {
		lexerToken := l.NextToken()

		// Checking Type And Literal
		if (lexerToken.TokenType != testToken.expectedType) ||
			(lexerToken.TokenLiteral != testToken.expectedLiteral) {
			t.Fatalf("Test Failed! Expected Token: <Type: %s, Literal: %s> but Got Token: <Type: %s, Literal: %s>\n",
				testToken.expectedType,
				testToken.expectedLiteral,
				lexerToken.TokenType,
				lexerToken.TokenLiteral)
		}
	}
}
