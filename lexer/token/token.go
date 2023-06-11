package token

type TokenType string

type Token struct {
	TokenType    TokenType
	TokenLiteral string
}

const (
	UNKNOWN TokenType = "UNKNOWN"
	EOF     TokenType = "EOF"

	// BRACES AND DELIMITERS
	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LCURLY TokenType = "{"
	RCURLY TokenType = "}"
	COMMA  TokenType = ","
	SCOLON TokenType = ";"

	// Arithmetic Operators
	ADD    TokenType = "+"
	SUB    TokenType = "-"
	ASSIGN TokenType = "="

	// User IDENTIFIERS
	IDENTIFIER TokenType = "IDENTIFIER"

	// Keywords
	BOX   TokenType = "BOX"
	PUT   TokenType = "PUT"
	UNBOX TokenType = "UNBOX"
	SHOW  TokenType = "SHOW"

	// Integers
	INT TokenType = "INT"
)

func NewToken(t_type TokenType, t_value string) Token {
	return Token{t_type, t_value}
}

func GetIdentifierType(identifier string) TokenType {
	switch identifier {
	case "box":
		return BOX
	case "put":
		return PUT
	case "unbox":
		return UNBOX
	case "show":
		return SHOW
	default:
		return IDENTIFIER
	}
}
