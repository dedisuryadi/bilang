package evaluator

import (
	"fmt"

	"github.com/dedisuryadi/bilang/ast"
	"github.com/dedisuryadi/bilang/token"
)

var (
	_NULL     = &Null{}
	_TRUE     = &Boolean{Value: true}
	_FALSE    = &Boolean{Value: false}
	_BREAK    = &Break{}
	_CONTINUE = &Continue{}
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

func (s *Script) Eval(node ast.Node, env *Environment) Object {
	switch node := node.(type) {
	case *ast.VarStatement:
		val := s.Eval(node.Value, env)
		if isError(val) {
			return val
		}
		name := node.Name.Value
		if _, ok := s.konst[name]; ok {
			return &Error{Message: fmt.Sprintf("konstanta %s tidak bisa ditugaskan kembali", name)}
		}
		if v, ok := env.Get(name); ok {
			from, to := v.Type(), val.Type()
			if from != to {
				return &Error{Message: fmt.Sprintf("perubahan tipe variabel %s dari %s menjadi %s tidak diizinkan", name, from, to)}
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
			return &Error{Message: fmt.Sprintf("konstanta %s tidak bisa ditugaskan kembali", konst)}
		}
		env.Set(konst, val)
		s.konst[konst] = struct{}{}

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &Function{Parameters: params, Env: env, Body: body}

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
		return &String{Value: node.Value}

	case *ast.LoopLiteral:
		return s.evalLoopExpression(node, env)
	case *ast.BreakExpression:
		return _BREAK
	case *ast.ContinueExpression:
		return _CONTINUE
	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.PilahExpression:
		return s.evalPilahExpression(node, env)

	case *ast.PilihStatement:
		val := s.Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}

	case *ast.Program:
		return s.evalProgram(node, env)

	case *ast.ExpressionStatement:
		return s.Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &Integer{Value: node.Value}

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
		return &Array{Elements: elements}

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

func (s *Script) evalHashLiteral(node *ast.HashLiteral, env *Environment) Object {
	pairs := make(map[HashKey]HashPair)
	for k, v := range node.Pairs {
		key := s.Eval(k, env)
		if isError(key) {
			return key
		}

		pk, ok := key.(Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := s.Eval(v, env)
		if isError(value) {
			return value
		}

		pairs[pk.HashKey()] = HashPair{Key: key, Value: value}
	}

	return &Hash{Pairs: pairs}
}

func evalIndexExpression(left, index Object) Object {
	switch {
	case left.Type() == ARRAY && index.Type() == INTEGER:
		return evalArrayIndexExpression(left, index)
	case left.Type() == HASH:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalHashIndexExpression(hash, index Object) Object {
	hashObject := hash.(*Hash)

	key, ok := index.(Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return _NULL
	}

	return pair.Value
}

func evalArrayIndexExpression(array, index Object) Object {
	arrayObj := array.(*Array)
	idx := index.(*Integer).Value
	max := int64(len(arrayObj.Elements) - 1)

	if idx < 0 || idx > max {
		return _NULL
	}

	return arrayObj.Elements[idx]
}

func nativeBoolToBooleanObject(value bool) Object {
	if value {
		return _TRUE
	}
	return _FALSE

}

func (s *Script) evalProgram(program *ast.Program, env *Environment) Object {
	var result Object
	for _, statement := range program.Statements {
		result = s.Eval(statement, env)
		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}
	return result
}

func (s *Script) evalBlockStatement(block *ast.BlockStatement, env *Environment) Object {
	var result Object
	for _, statement := range block.Statements {
		result = s.Eval(statement, env)
		if result != nil {
			rt := result.Type()
			if rt == RETURN || rt == ERROR {
				return result
			}
		}
	}
	return result
}

func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

func evalBangOperatorExpression(right Object) Object {
	switch right {
	case _TRUE:
		return _FALSE
	case _FALSE:
		return _TRUE
	case _NULL:
		return _TRUE
	default:
		return _FALSE
	}
}

func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == INTEGER && right.Type() == INTEGER:
		return evalIntegerInfixExpression(operator, left, right)

	case left.Type() == STRING && right.Type() == STRING:
		return evalStringInfixExpression(operator, left, right)

	case left.Type() == BOOLEAN && right.Type() == BOOLEAN:
		lVal := left.(*Boolean).Value
		rVal := right.(*Boolean).Value
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

func evalStringInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value
	switch operator {
	case "+":
		return &String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value
	switch operator {
	case "%":
		return &Integer{Value: leftVal % rightVal}
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		return &Integer{Value: leftVal / rightVal}
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

func (s *Script) evalJikaExpression(je *ast.JikaExpression, env *Environment) Object {
	cond := s.Eval(je.Condition, env)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return s.Eval(je.Consequence, env)
	} else if je.Alternative != nil {
		return s.Eval(je.Alternative, env)
	} else {
		return _NULL
	}
}

func (s *Script) evalPilahExpression(ps *ast.PilahExpression, env *Environment) Object {
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

	return _NULL
}

func (s *Script) evalPipeExpression(p *ast.Pipe, env *Environment) Object {
	left := s.Eval(p.Left, env)

	// Convert the type object back to an expression, so it can be passed to the FunctionCall arg
	argument := obj2Expression(left)
	if argument == nil {
		return _NULL
	}

	// The right side operator should be a function.
	switch rightFunc := p.Right.(type) {
	case *ast.MethodCallExpression:
		switch rightFunc.Call.(type) {
		case *ast.Identifier:
			//e.g.
			//x = ["hello", "world"] |> strings.upper    : rightFunc.Call.(type) == *ast.Identifier
			//x = ["hello", "world"] |> strings.upper()  : rightFunc.Call.(type) == *ast.CallExpression
			//so here we convert *ast.Identifier to *ast.CallExpression
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

	return _NULL
}

func evalIdentifier(node *ast.Identifier, env *Environment) Object {
	if b, ok := builtins[node.Value]; ok {
		return b
	}
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}
	return val
}

func (s *Script) evalExpression(exps []ast.Expression, env *Environment) []Object {
	var result []Object
	for _, e := range exps {
		evaluated := s.Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func (s *Script) applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		if paramLen, argsLen := len(fn.Parameters), len(args); paramLen != argsLen {
			return newError("invalid length between function parameter=%d & args=%d", paramLen, argsLen)
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := s.Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func (s *Script) evalLoopExpression(node *ast.LoopLiteral, env *Environment) Object {
	iter := evalIdentifier(node.Iter, env)
	if iterable, ok := iter.(Iterable); !ok || !iterable.Iter() {
		return newError("identifier %s is not iterable", node.Iter)
	}

	complete := len(node.KV) > 1
	switch iter := iter.(type) {
	case *String:
		for k, v := range []rune(iter.Value) {
			scope := NewEnclosedEnvironment(env)
			scope.Set(node.KV[0].Value, &Integer{Value: int64(k)})
			if complete {
				scope.Set(node.KV[1].Value, &String{Value: string(v)})
			}
			res := s.evalBlockStatement(node.Body, scope)
			if res.Type() == ERROR {
				return res
			}
			if _, ok := res.(*Continue); ok {
				continue
			}
			if _, ok := res.(*Break); ok {
				break
			}
		}

	case *Hash:
		for _, v := range iter.Pairs {
			scope := NewEnclosedEnvironment(env)
			scope.Set(node.KV[0].Value, v.Key)
			if complete {
				scope.Set(node.KV[1].Value, v.Value)
			}
			res := s.evalBlockStatement(node.Body, scope)
			if res.Type() == ERROR {
				return res
			}
			if _, ok := res.(*Continue); ok {
				continue
			}
			if _, ok := res.(*Break); ok {
				break
			}
		}

	case *Array:
		for k, v := range iter.Elements {
			scope := NewEnclosedEnvironment(env)
			scope.Set(node.KV[0].Value, &Integer{Value: int64(k)})
			if complete {
				scope.Set(node.KV[1].Value, v)
			}
			res := s.evalBlockStatement(node.Body, scope)
			if res.Type() == ERROR {
				return res
			}
			if _, ok := res.(*Continue); ok {
				continue
			}
			if _, ok := res.(*Break); ok {
				break
			}
		}

	default:
		return newError("type %s is not iterable", iter)
	}

	return _NULL
}

func extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)
	for index, param := range fn.Parameters {
		env.Set(param.Value, args[index])
	}
	return env
}

func unwrapReturnValue(obj Object) Object {
	if retVal, ok := obj.(*ReturnValue); ok {
		return retVal.Value
	}
	return obj
}

func isTruthy(obj Object) bool {
	switch obj {
	case _TRUE:
		return true
	case _NULL:
		return false
	case _FALSE:
		return false
	default:
		switch obj := obj.(type) {
		case *Integer:
			if obj.Value > 0 {
				return true
			}
		}
		return false
	}
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR
	}
	return false
}

// Convert an Object to an ast.Expression.
func obj2Expression(obj Object) ast.Expression {
	switch value := obj.(type) {
	case *Boolean:
		return &ast.Boolean{Value: value.Value}
	case *Integer:
		return &ast.IntegerLiteral{Value: value.Value}
	case *String:
		return &ast.StringLiteral{Value: value.Value}
	case *Null:
		return &ast.NihilLiteral{}
	case *Array:
		array := &ast.ArrayLiteral{}
		for _, v := range value.Elements {
			result := obj2Expression(v)
			if result == nil {
				return nil
			}
			array.Elements = append(array.Elements, result)
		}
		return array
	case *Hash:
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
