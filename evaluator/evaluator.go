package evaluator

import (
	"fmt"
	"log"

	"github.com/dedisuryadi/bilang/ast"
	"github.com/dedisuryadi/bilang/object"
	"github.com/dedisuryadi/bilang/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

type Script struct {
	konst map[string]struct{}
}

func NewScript() *Script {
	return &Script{konst: make(map[string]struct{})}
}

func (s *Script) Free() {
	s.konst = nil
}

func (s *Script) Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.VarStatement:
		val := s.Eval(node.Value, env)
		if isError(val) {
			return val
		}
		name := node.Name.Value
		if _, ok := s.konst[name]; ok {
			return &object.Error{Message: fmt.Sprintf("konstanta %s tidak bisa ditugaskan kembali", name)}
		}
		if v, ok := env.Get(name); ok {
			from, to := v.Type(), val.Type()
			if from != to {
				return &object.Error{Message: fmt.Sprintf("perubahan tipe variabel %s dari %s menjadi %s tidak diizinkan", name, from, to)}
			}
		}
		env.Set(name, val)

	case *ast.KonstStatement:
		val := s.Eval(node.Value, env)
		if isError(val) {
			return val
		}
		konst := node.Name.Value
		if _, ok := s.konst[konst]; ok {
			return &object.Error{Message: fmt.Sprintf("konstanta %s tidak bisa ditugaskan kembali", konst)}
		}
		env.Set(konst, val)
		s.konst[konst] = struct{}{}

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		fn := s.Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := s.evalExpression(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return s.applyFunction(fn, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.PilahExpression:
		return s.evalPilahExpression(node, env)

	case *ast.PilihStatement:
		val := s.Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.Program:
		return s.evalProgram(node, env)

	case *ast.ExpressionStatement:
		return s.Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.PrefixExpression:
		right := s.Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := s.Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := s.Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.BlockStatement:
		return s.evalBlockStatement(node, env)

	case *ast.JikaExpression:
		return s.evalJikaExpression(node, env)

	case *ast.ArrayLiteral:
		elements := s.evalExpression(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.Pipe:
		return s.evalPipeExpression(node, env)

	case *ast.IndexExpression:
		left := s.Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := s.Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.HashLiteral:
		return s.evalHashLiteral(node, env)
	}

	return nil
}

func (s *Script) evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for k, v := range node.Pairs {
		key := s.Eval(k, env)
		if isError(key) {
			return key
		}

		pk, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := s.Eval(v, env)
		if isError(value) {
			return value
		}

		pairs[pk.HashKey()] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObj.Elements[idx]
}

func nativeBoolToBooleanObject(value bool) object.Object {
	if value {
		return TRUE
	}
	return FALSE

}

func (s *Script) evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = s.Eval(statement, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func (s *Script) evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = s.Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN || rt == object.ERROR {
				return result
			}
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
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

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left, right)

	case left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN:
		lVal := left.(*object.Boolean).Value
		rVal := right.(*object.Boolean).Value
		switch operator {
		case "&&":
			return nativeBoolToBooleanObject(lVal && rVal)
		case "||":
			return nativeBoolToBooleanObject(lVal || rVal)
		case "==":
			return nativeBoolToBooleanObject(left == right)
		case "!=":
			return nativeBoolToBooleanObject(left != right)
		default:
			return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
		}

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (s *Script) evalJikaExpression(je *ast.JikaExpression, env *object.Environment) object.Object {
	cond := s.Eval(je.Condition, env)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return s.Eval(je.Consequence, env)
	} else if je.Alternative != nil {
		return s.Eval(je.Alternative, env)
	} else {
		return NULL
	}
}

func (s *Script) evalPilahExpression(ps *ast.PilahExpression, env *object.Environment) object.Object {
	target := s.Eval(ps.Target, env)
	if isError(target) {
		return target
	}

	for i, v := range ps.Conditions {
		if v.Token.Type == token.UNDERSCORE {
			return s.Eval(ps.Values[i], env)
		}
		cond := s.Eval(v, env)
		val := evalInfixExpression("==", target, cond)
		if isTruthy(val) {
			return s.Eval(ps.Values[i], env)
		}
	}

	return NULL
}

func (s *Script) evalPipeExpression(p *ast.Pipe, env *object.Environment) object.Object {
	left := s.Eval(p.Left, env)

	// Convert the type object back to an expression
	// so it can be passed to the FunctionCall arguments.
	argument := obj2Expression(left)
	if argument == nil {
		return NULL
	}

	// The right side operator should be a function.
	switch rightFunc := p.Right.(type) {
	case *ast.MethodCallExpression:
		// Prepend the left-hand interpreted value
		// to the function arguments.
		switch rightFunc.Call.(type) {
		case *ast.Identifier:
			//e.g.
			//x = ["hello", "world"] |> strings.upper    : rightFunc.Call.(type) == *ast.Identifier
			//x = ["hello", "world"] |> strings.upper()  : rightFunc.Call.(type) == *ast.CallExpression
			//so here we convert *ast.Identifier to * ast.CallExpression
			rightFunc.Call = &ast.CallExpression{Token: p.Token, Function: rightFunc.Call}
		}
		rightFunc.Call.(*ast.CallExpression).Arguments = append([]ast.Expression{argument}, rightFunc.Call.(*ast.CallExpression).Arguments...)
		return s.Eval(rightFunc, env)

	case *ast.CallExpression:
		rightFunc.Arguments = append([]ast.Expression{argument}, rightFunc.Arguments...)
		return s.Eval(rightFunc, env)

	case *ast.Identifier:
		fn := s.Eval(rightFunc, env)
		if isError(fn) {
			return fn
		}
		args := s.evalExpression([]ast.Expression{argument}, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return s.applyFunction(fn, args)
	default:
		// TODO: handle lambda
		fmt.Printf("unhandled pipe %T %s\n", rightFunc, rightFunc)
	}

	return NULL
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if b, ok := builtins[node.Value]; ok {
		return b
	}
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}
	return val
}

func (s *Script) evalExpression(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := s.Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func (s *Script) applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := s.Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	if paramLen, argsLen := len(fn.Parameters), len(args); paramLen != argsLen {
		// TODO: handle runtime error
		log.Fatalf("invalid length between function parameter=%d & args=%d", paramLen, argsLen)
	}
	env := object.NewEnclosedEnvironment(fn.Env)
	for index, param := range fn.Parameters {
		env.Set(param.Value, args[index])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if retVal, ok := obj.(*object.ReturnValue); ok {
		return retVal.Value
	}
	return obj
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case TRUE:
		return true
	case NULL:
		return false
	case FALSE:
		return false
	default:
		switch obj := obj.(type) {
		case *object.Integer:
			if obj.Value > 0 {
				return true
			}
		}
		return false
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}

// Convert an Object to an ast.Expression.
func obj2Expression(obj object.Object) ast.Expression {
	switch value := obj.(type) {
	case *object.Boolean:
		return &ast.Boolean{Value: value.Value}
	case *object.Integer:
		return &ast.IntegerLiteral{Value: value.Value}
	case *object.String:
		return &ast.StringLiteral{Value: value.Value}
	case *object.Null:
		return &ast.NihilLiteral{}
	case *object.Array:
		array := &ast.ArrayLiteral{}
		for _, v := range value.Elements {
			result := obj2Expression(v)
			if result == nil {
				return nil
			}
			array.Elements = append(array.Elements, result)
		}
		return array
	case *object.Hash:
		hash := &ast.HashLiteral{}
		hash.Pairs = make(map[ast.Expression]ast.Expression)
		for hk := range value.Pairs { //hk:hash key
			v, _ := value.Pairs[hk]
			key := &ast.StringLiteral{Value: v.Key.Inspect()}
			result := obj2Expression(v.Value)
			if result == nil {
				return nil
			}
			hash.Pairs[key] = result
		}
		return hash
	}
	return nil
}
