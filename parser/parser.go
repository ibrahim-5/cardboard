package parser

import (
	"cardboard/lexer"
	"cardboard/lexer/token"
	"cardboard/parser/ast"
	"fmt"
)

type Parser struct {
	lexer       *lexer.Lexer
	curToken    token.Token
	peekToken   token.Token
	errors      []string
	prefixFuncs map[token.TokenType]prefixFunc
	infixFuncs  map[token.TokenType]infixFunc
}

func CreateParser(l *lexer.Lexer) Parser {
	p := Parser{lexer: l}

	// Need to initialize both Tokens Pointers
	p.nextToken()
	p.nextToken()

	return p
}

// Parses Cardboard Program
func (p *Parser) ParseCardBoard() ast.Program {
	program := ast.Program{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.TokenType {
	case token.PUT:
		return p.parsePutStatement()
	case token.UNBOX:
		return p.parseUnboxStatement()
	default:
		return nil
	}
}

func (p *Parser) parsePutStatement() *ast.PutStatement {
	putStmt := &ast.PutStatement{}

	// Parse Put Token
	putStmt.NodeToken = p.curToken

	// Parse Identifier
	if !p.expectPeek(token.IDENTIFIER) {
		p.typeError(token.IDENTIFIER, p.peekToken.TokenType)
		return nil
	}

	putStmt.NodeIdentifier = ast.Identifier{NodeToken: p.curToken, Value: p.curToken.TokenLiteral}

	// Ensure Next Token is Assign
	if !p.expectPeek(token.ASSIGN) {
		p.typeError(token.ASSIGN, p.peekToken.TokenType)
		return nil
	}
	// TODO: Parse Expression
	for !p.curTokenIs(token.SCOLON) {
		p.nextToken()
	}

	return putStmt
}

func (p *Parser) parseUnboxStatement() *ast.UnboxStatement {
	unboxStmt := &ast.UnboxStatement{}

	// Parse Unbox Token
	unboxStmt.NodeToken = p.curToken

	// TODO: Parse Expression
	for !p.curTokenIs(token.SCOLON) {
		p.nextToken()
	}

	return unboxStmt
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.TokenType == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.TokenType == t
}

func (p *Parser) getErrors() []string {
	return p.errors
}

// In the case where the statement is invalid, we'll
// need to skip it!
func (p *Parser) skipStatement() {
	for !(p.curTokenIs(token.SCOLON)) {
		p.nextToken()
	}
}

func (p *Parser) typeError(expectedType token.TokenType, gotType token.TokenType) {
	p.errors = append(p.errors,
		fmt.Sprintf("Error. Expected Token Type <%s>. Got Token Type <%s>.\n", expectedType, gotType))
	p.skipStatement()
}

// Parser Functions -> Certain token types are associated to infix and prefix operations.
type (
	infixFunc  func(ast.Expression) ast.Expression
	prefixFunc func() ast.Expression
)

func (p *Parser) setInfixFunction(token token.TokenType, function infixFunc) {
	p.infixFuncs[token] = function
}

func (p *Parser) setPrefixFunction(token token.TokenType, function prefixFunc) {
	p.prefixFuncs[token] = function
}
