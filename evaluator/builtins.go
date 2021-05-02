package evaluator

import (
	"fmt"

	"github.com/dedisuryadi/bilang/object"
)

var builtins = map[string]*object.Builtin{
	"panjang": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			default:
				return newError("argument to `panjang` not supported, got %s", args[0].Type())
			}
		},
	},
	"awal": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newError("argument to `awal` must be ARRAY. got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"akhir": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newError("argument to `akhir` must be ARRAY. got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			Len := len(arr.Elements)
			if Len > 0 {
				return arr.Elements[Len-1]
			}
			return NULL
		},
	},
	"ekor": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newError("argument to `ekor` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			Len := len(arr.Elements)
			if Len > 0 {
				newElements := make([]object.Object, Len-1, Len-1)
				copy(newElements, arr.Elements[1:Len])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			Len := len(arr.Elements)
			newElements := make([]object.Object, Len+1, Len+1)
			copy(newElements, arr.Elements)
			newElements[Len] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"stdout": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
