package lexer

import "cardboard/token"

type Lexer struct {
	data    string
	curPos  int
	nextPos int
	char    byte
}

func CreateLexer(inputData string) *Lexer {
	lexer := Lexer{data: inputData}
	lexer.readChar()
	return &lexer
}

func (lex *Lexer) NextToken() token.Token {
	var curToken token.Token

	lex.eatWhiteSpace()

	switch lex.char {

	// BRACES AND DELIMITERS
	case '(':
		curToken = token.NewToken(token.LBRAC, "(")
	case ')':
		curToken = token.NewToken(token.RBRAC, ")")
	case '{':
		curToken = token.NewToken(token.LCURLY, "{")
	case '}':
		curToken = token.NewToken(token.RCURLY, "}")
	case ',':
		curToken = token.NewToken(token.COMMA, ",")
	case ';':
		curToken = token.NewToken(token.SCOLON, ";")

	// Arithmetic Operators
	case '+':
		curToken = token.NewToken(token.ADD, "+")
	case '-':
		curToken = token.NewToken(token.SUB, "-")
	case '=':
		curToken = token.NewToken(token.EQ, "=")

	// EOF
	case 0:
		curToken = token.NewToken(token.EOF, "")

	default:
		if isInteger(lex.char) {
			readInteger := lex.readInteger()
			curToken = token.NewToken(token.INT, readInteger)
		} else if isLetter(lex.char) {
			readIdentifier := lex.readIdentifier()
			curToken = token.NewToken(token.GetIdentifierType(readIdentifier), readIdentifier)
		}
		return curToken
	}
	lex.readChar()
	return curToken
}

func (lex *Lexer) readChar() {
	if lex.nextPos >= len(lex.data) {
		lex.char = 0
	} else {
		lex.char = lex.data[lex.nextPos]
	}
	lex.curPos = lex.nextPos
	lex.nextPos++
}

func (lex *Lexer) readIdentifier() string {
	startPos := lex.curPos
	for isLetter(lex.char) {
		lex.readChar()
	}
	return string(lex.data[startPos:lex.curPos])
}

func (lex *Lexer) readInteger() string {
	startPos := lex.curPos
	for isInteger(lex.char) {
		lex.readChar()
	}
	return string(lex.data[startPos:lex.curPos])
}

func isLetter(ch byte) bool {
	if (65 <= ch && ch <= 90) || (97 <= ch && ch <= 122) {
		return true
	}
	return false
}

func isInteger(ch byte) bool {
	if 48 <= ch && ch <= 57 {
		return true
	}
	return false
}

func (lex *Lexer) eatWhiteSpace() {
	for lex.char == '\n' || lex.char == ' ' || lex.char == '\t' {
		lex.readChar()
	}
}
