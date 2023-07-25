package eval

import (
	"cardboard/object"
	"cardboard/parser/ast"
)

var NULL = &object.Null{}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements, env)
	case *ast.UnboxStatement:
		return evalUnboxStatement(node, env)
	case *ast.PutStatement:
		return evalPutStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	// Expressions
	case *ast.BoxExpression:
		return &object.Box{Env: env,
			ParameterList: node.ParameterList,
			Body:          node.Body}
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		l := Eval(node.Left, env)
		r := Eval(node.Right, env)
		return evalInfixExpression(node.Operator, l, r)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}
	return NULL
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)

		if result.Type() == object.UNBOX_OBJ {
			return result.(*object.Unbox).Value
		}
	}
	return result
}

func evalPrefixExpression(expr *ast.PrefixExpression, env *object.Environment) object.Object {
	operand := Eval(expr.Right, env)

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

func evalUnboxStatement(unbox *ast.UnboxStatement, env *object.Environment) object.Object {
	val := Eval(unbox.NodeExpression, env)
	return &object.Unbox{Value: val}
}

func evalPutStatement(stmt *ast.PutStatement, env *object.Environment) object.Object {
	val := Eval(stmt.NodeExpression, env)
	return env.Set(stmt.NodeIdentifier.Value, val)
}

func evalIdentifier(ident *ast.Identifier, env *object.Environment) object.Object {
	obj, ok := env.Get(ident.Value)
	if !ok {
		return NULL
	}
	return obj
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.UNBOX_OBJ {
				return result
			}
		}
	}
	return result
}

func evalCallExpression(call *ast.CallExpression, env *object.Environment) object.Object {
	var arguments []object.Object

	box := Eval(call.Function, env)

	for _, arg := range call.Arguments {
		evaluated := Eval(arg, env)
		arguments = append(arguments, evaluated)
	}

	return applyBoxFunction(box, arguments)
}

func applyBoxFunction(box object.Object, args []object.Object) object.Object {
	fn, ok := box.(*object.Box)

	if !ok {
		return NULL
	}

	fn.Env = object.CreateEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.ParameterList {
		fn.Env.Set(param.Value, args[paramIdx])
	}

	evaluated := Eval(fn.Body, fn.Env)

	if unbox, ok := evaluated.(*object.Unbox); ok {
		return unbox.Value
	}
	return evaluated
}
