package evaluator

import (
	"fmt"
	"morty/ast"
	"morty/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Noder, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return ToBoolObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatements(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		funcLit_obj := evalFunctionLiteral(node, env) // Note: when evaling function decleration we only make an funcLit object
		if node.Name != nil {
			env.Set(node.Name.Value, funcLit_obj)
			return nil
		}
		return funcLit_obj

	case *ast.CallExpression: // when evaling func call we actually eval the funcLit obj
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}

	return nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {

	var result object.Object

	for _, statements := range stmts {
		result = Eval(statements, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func ToBoolObject(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(rigth object.Object) object.Object {
	if rigth.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", rigth.Type())
	}

	if rigth.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := rigth.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, leftExp object.Object, rightExp object.Object) object.Object {
	// if we compare two struct values directly, then go compares each fields one by one
	// if we compare two struct pointers, then go checks if they both point to the same memory

	switch {
	case leftExp.Type() == object.INTEGER_OBJ && rightExp.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, leftExp, rightExp)

	case leftExp.Type() == object.STRING_OBJ && rightExp.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, leftExp, rightExp)

	case operator == "==":
		return ToBoolObject(leftExp == rightExp) // left and right both point to either TRUE or FALSE structs

	case operator == "!=":
		return ToBoolObject(leftExp != rightExp)

	case leftExp.Type() != rightExp.Type():
		return newError("type mismatch: %s %s %s", leftExp.Type(), operator, rightExp.Type())

	default:
		return newError("unknown operator: %s %s %s", leftExp.Type(), operator, rightExp.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return ToBoolObject(leftVal < rightVal)
	case ">":
		return ToBoolObject(leftVal > rightVal)
	case "==":
		return ToBoolObject(leftVal == rightVal)
	case "!=":
		return ToBoolObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

func evalIfExpression(ife *ast.IfExpression, env *object.Environment) object.Object {

	condition := Eval(ife.Condition, env)

	if isTruthy(condition) {
		return Eval(ife.Concequence, env)
	} else if ife.Alternative != nil {
		return Eval(ife.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(conditionObj object.Object) bool {
	switch conditionObj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {

	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		} // this gives return object

	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv := setFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func setFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramidx, param := range fn.Parameters {
		env.Set(param.Value, args[paramidx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}

	case "!=":
		return ToBoolObject(leftVal != rightVal)

	case "==":
		return ToBoolObject(leftVal == rightVal)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

}

func evalFunctionLiteral(funcLit *ast.FunctionLiteral, env *object.Environment) object.Object {
	params := funcLit.Parameters
	body := funcLit.Body
	name := funcLit.Name
	return &object.Function{Parameters: params, Name: name, Body: body, Env: env}
}
