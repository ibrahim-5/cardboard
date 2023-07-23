package eval

import (
	"cardboard/object"
	"cardboard/parser/ast"
)

var NULL = &object.Null{}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)

	// Expressions
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node)
	case *ast.InfixExpression:
		l := Eval(node.Left)
		r := Eval(node.Right)
		return evalInfixExpression(node.Operator, l, r)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return NULL
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}

func evalPrefixExpression(expr *ast.PrefixExpression) object.Object {
	operand := Eval(expr.Right)

	if operand.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := operand.(*object.Integer).Value

	switch expr.Operator {
	case "+":
		return &object.Integer{Value: value}
	case "-":
		return &object.Integer{Value: -value}
	}

	return NULL
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	// Infix Operations only defined for integers for now.
	if left.Type() != object.INTEGER_OBJ || right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	}

	return NULL
}
