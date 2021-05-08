package evaluator

import (
	"fmt"
)

var builtins = map[string]*Builtin{
	"panjang": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *String:
				return &Float{Value: float64(len(arg.Value))}

			case *Array:
				return &Float{Value: float64(len(arg.Elements))}

			default:
				return newError("argument to `panjang` not supported, got %s", args[0].Type())
			}
		},
	},
	"awal": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY {
				return newError("argument to `awal` must be ARRAY. got %s", args[0].Type())
			}
			arr := args[0].(*Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return _NULL
		},
	},
	"akhir": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY {
				return newError("argument to `akhir` must be ARRAY. got %s", args[0].Type())
			}
			arr := args[0].(*Array)
			Len := len(arr.Elements)
			if Len > 0 {
				return arr.Elements[Len-1]
			}
			return _NULL
		},
	},
	"ekor": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY {
				return newError("argument to `ekor` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*Array)
			Len := len(arr.Elements)
			if Len > 0 {
				newElements := make([]Object, Len-1, Len-1)
				copy(newElements, arr.Elements[1:Len])
				return &Array{Elements: newElements}
			}
			return _NULL
		},
	},
	"push": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*Array)
			Len := len(arr.Elements)
			newElements := make([]Object, Len+1, Len+1)
			copy(newElements, arr.Elements)
			newElements[Len] = args[1]
			return &Array{Elements: newElements}
		},
	},
	"stdout": {
		Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return _NULL
		},
	},
	"println": {
		Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return _NULL
		},
	},
}
