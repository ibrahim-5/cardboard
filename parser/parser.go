package parser

import (
	"cardboard/lexer"
	"cardboard/lexer/token"
	"cardboard/parser/ast"
	"fmt"
	"strconv"
)

type Parser struct {
	lexer       *lexer.Lexer
	curToken    token.Token
	peekToken   token.Token
	errors      []string
	prefixFuncs map[token.TokenType]prefixFunc
	infixFuncs  map[token.TokenType]infixFunc
}

// Precedence For Operators < Not all supported right now! >
const (
	_ int = iota
	LOWEST
	EQUALS  // == LESSGREATER // > or <
	SUM     //+
	PRODUCT //*
	PREFIX  //-X or !X
	CALL    // myFunction(X)
)

// Precedence mapping
var precedence = map[token.TokenType]int{
	token.ADD: SUM,
	token.SUB: SUM,
}

func CreateParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

	// Need to initialize both Tokens Pointers
	p.nextToken()
	p.nextToken()

	// Instantiate Mapping
	// Prefix
	p.prefixFuncs = make(map[token.TokenType]prefixFunc)
	p.setPrefixFunction(token.IDENTIFIER, p.parseIdentifier)
	p.setPrefixFunction(token.INT, p.parseIntegerLiteral)
	p.setPrefixFunction(token.SUB, p.parsePrefixExpression)

	// Infix
	p.infixFuncs = make(map[token.TokenType]infixFunc)
	p.setInfixFunction(token.ADD, p.parseInfixExpression)
	p.setInfixFunction(token.SUB, p.parseInfixExpression)

	return p
}

// Parses Cardboard Program
func (p *Parser) ParseCardBoard() ast.Program {
	program := ast.Program{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
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
		return p.parseExpressionStatement()
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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expStmt := &ast.ExpressionStatement{NodeToken: p.curToken}
	expStmt.Expression = p.parseExpression(LOWEST)

	// Optional Semi-colon
	if p.peekTokenIs(token.SCOLON) {
		p.nextToken()
	}
	return expStmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFuncs[p.curToken.TokenType]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("Couldn't find prefix function for %s", p.curToken.TokenLiteral))
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SCOLON) && precedence < p.peekPrecedence() {
		infix := p.infixFuncs[p.peekToken.TokenType]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		NodeToken: p.curToken,
		Value:     p.curToken.TokenLiteral,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.TokenLiteral, 10, 0)
	if err != nil {
		error := fmt.Sprintf("Integer Parse Error. Couldn't Parse Integer From String = %s", p.curToken.TokenLiteral)
		p.errors = append(p.errors, error)
		return nil
	}
	return &ast.IntegerLiteral{NodeToken: p.curToken, Value: val}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		NodeToken: p.curToken,
		Operator:  p.curToken.TokenLiteral,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		NodeToken: p.curToken,
		Left:      left,
		Operator:  p.curToken.TokenLiteral,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedence[p.peekToken.TokenType]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedence[p.curToken.TokenType]; ok {
		return precedence
	}
	return LOWEST
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
