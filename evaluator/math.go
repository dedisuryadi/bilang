package evaluator

import (
	"math"
)

var mathBuiltin = map[string]*Builtin{
	"math.Abs": {
		Fn: func(args ...Object) Object {
			switch arg := args[0].(type) {
			case *Float:
				return &Float{Value: math.Abs(arg.Value)}
			default:
				return NewError("fungsi math.Abs hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
		},
	},
	//Acos(x float64) float64
	"math.Acos": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Acos parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Acos hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Acos(x.Value)}
		},
	},
	//Acosh(x float64) float64
	"math.Acosh": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Acosh parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Acosh hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Acosh(x.Value)}
		},
	},
	//Asin(x float64) float64
	"math.Asin": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Asin parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Asin hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Asin(x.Value)}
		},
	},
	//Asinh(x float64) float64
	"math.Asinh": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Asinh parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Asinh hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Asinh(x.Value)}
		},
	},
	//Atan(x float64) float64
	"math.Atan": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Atan parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Atan hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Atan(x.Value)}
		},
	},
	//Atan2(x, y float64) float64
	"math.Atan2": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return NewError("fungsi math.Atan2 parameter sebanyak 2, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Atan2 hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			y, ok := args[1].(*Float)
			if !ok {
				return NewError("fungsi math.Atan2 hanya bisa menerima ANGKA, didapat: %s", args[1].Type())
			}
			return &Float{Value: math.Atan2(x.Value, y.Value)}
		},
	},
	//Atanh(x float64) float64
	"math.Atanh": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Atanh parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Atanh hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Atanh(x.Value)}
		},
	},
	//Cbrt(x float64) float64
	"math.Cbrt": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Cbrt parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Cbrt hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Cbrt(x.Value)}
		},
	},
	//Ceil(x float64) float64
	"math.Ceil": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Ceil parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Ceil hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Ceil(x.Value)}
		},
	},
	//Copysign(x, y float64) float64
	"math.Copysign": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return NewError("fungsi math.Copysign parameter sebanyak 2, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Copysign hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			y, ok := args[1].(*Float)
			if !ok {
				return NewError("fungsi math.Copysign hanya bisa menerima ANGKA, didapat: %s", args[1].Type())
			}
			return &Float{Value: math.Copysign(x.Value, y.Value)}
		},
	},
	//Cos(x float64) float64
	"math.Cos": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Cos parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Cos hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Cos(x.Value)}
		},
	},
	//Cosh(x float64) float64
	"math.Cosh": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Cosh parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Cosh hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Cosh(x.Value)}
		},
	},
	//Dim(x, y float64) float64
	"math.Dim": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return NewError("fungsi math.Dim parameter sebanyak 2, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Dim hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			y, ok := args[1].(*Float)
			if !ok {
				return NewError("fungsi math.Dim hanya bisa menerima ANGKA, didapat: %s", args[1].Type())
			}
			return &Float{Value: math.Dim(x.Value, y.Value)}
		},
	},
	//Erf(x float64) float64
	"math.Erf": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Erf parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Erf hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Erf(x.Value)}
		},
	},
	//Erfc(x float64) float64
	"math.Erfc": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Erfc parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Erfc hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Erfc(x.Value)}
		},
	},
	//Erfcinv(x float64) float64
	"math.Erfcinv": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Erfcinv parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Erfcinv hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Erfcinv(x.Value)}
		},
	},
	//Erfinv(x float64) float64
	"math.Erfinv": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Erfinv parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Erfinv hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Erfinv(x.Value)}
		},
	},
	//Exp(x float64) float64
	"math.Exp": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Exp parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Exp hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Exp(x.Value)}
		},
	},
	//Exp2(x float64) float64
	"math.Exp2": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Exp2 parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Exp2 hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Exp2(x.Value)}
		},
	},
	//Expm1(x float64) float64
	"math.Expm1": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Expm1 parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Expm1 hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Expm1(x.Value)}
		},
	},
	//FMA(x, y, z float64) float64
	"math.FMA": {
		Fn: func(args ...Object) Object {
			if len(args) != 3 {
				return NewError("fungsi math.FMA parameter sebanyak 3, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.FMA hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			y, ok := args[1].(*Float)
			if !ok {
				return NewError("fungsi math.FMA hanya bisa menerima ANGKA, didapat: %s", args[1].Type())
			}
			z, ok := args[2].(*Float)
			if !ok {
				return NewError("fungsi math.FMA hanya bisa menerima ANGKA, didapat: %s", args[2].Type())
			}
			return &Float{Value: math.FMA(x.Value, y.Value, z.Value)}
		},
	},
	//Floor(x float64) float64
	"math.Floor": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Floor parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Floor hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Floor(x.Value)}
		},
	},
	//Gamma(x float64) float64
	"math.Gamma": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Gamma parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Gamma hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Gamma(x.Value)}
		},
	},
	//Hypot(p, q float64) float64
	"math.Hypot": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return NewError("fungsi math.Hypot parameter sebanyak 2, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Hypot hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			y, ok := args[1].(*Float)
			if !ok {
				return NewError("fungsi math.Hypot hanya bisa menerima ANGKA, didapat: %s", args[1].Type())
			}
			return &Float{Value: math.Hypot(x.Value, y.Value)}
		},
	},
	//J0(x float64) float64
	"math.J0": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.J0 parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.J0 hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.J0(x.Value)}
		},
	},
	//J1(x float64) float64
	"math.J1": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.J1 parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.J1 hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.J1(x.Value)}
		},
	},
	//Log(x float64) float64
	"math.Log": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Log parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Log hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Log(x.Value)}
		},
	},
	//Log10(x float64) float64
	"math.Log10": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Log10 parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Log10 hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Log10(x.Value)}
		},
	},
	//Log1p(x float64) float64
	"math.Log1p": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Log1p parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Log1p hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Log1p(x.Value)}
		},
	},
	//Log2(x float64) float64
	"math.Log2": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Log2 parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Log2 hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Log2(x.Value)}
		},
	},
	//Logb(x float64) float64
	"math.Logb": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Logb parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Logb hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Logb(x.Value)}
		},
	},
	//Max(x, y float64) float64
	"math.Max": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return NewError("fungsi math.Max parameter sebanyak 2, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Max hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			y, ok := args[1].(*Float)
			if !ok {
				return NewError("fungsi math.Max hanya bisa menerima ANGKA, didapat: %s", args[1].Type())
			}
			return &Float{Value: math.Max(x.Value, y.Value)}
		},
	},
	//Min(x, y float64) float64
	"math.Min": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return NewError("fungsi math.Min parameter sebanyak 2, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Min hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			y, ok := args[1].(*Float)
			if !ok {
				return NewError("fungsi math.Min hanya bisa menerima ANGKA, didapat: %s", args[1].Type())
			}
			return &Float{Value: math.Min(x.Value, y.Value)}
		},
	},
	//Mod(x, y float64) float64
	"math.Mod": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return NewError("fungsi math.Mod parameter sebanyak 2, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Mod hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			y, ok := args[1].(*Float)
			if !ok {
				return NewError("fungsi math.Mod hanya bisa menerima ANGKA, didapat: %s", args[1].Type())
			}
			return &Float{Value: math.Mod(x.Value, y.Value)}
		},
	},
	//Pow(x, y float64) float64
	"math.Pow": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return NewError("fungsi math.Pow parameter sebanyak 2, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Pow hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			y, ok := args[1].(*Float)
			if !ok {
				return NewError("fungsi math.Pow hanya bisa menerima ANGKA, didapat: %s", args[1].Type())
			}
			return &Float{Value: math.Pow(x.Value, y.Value)}
		},
	},
	//Remainder(x, y float64) float64
	"math.Remainder": {
		Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return NewError("fungsi math.Remainder parameter sebanyak 2, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Remainder hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			y, ok := args[1].(*Float)
			if !ok {
				return NewError("fungsi math.Remainder hanya bisa menerima ANGKA, didapat: %s", args[1].Type())
			}
			return &Float{Value: math.Remainder(x.Value, y.Value)}
		},
	},
	//Round(x float64) float64
	"math.Round": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Round parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Round hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Round(x.Value)}
		},
	},
	//RoundToEven(x float64) float64
	"math.RoundToEven": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.RoundToEven parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.RoundToEven hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.RoundToEven(x.Value)}
		},
	},
	//Sin(x float64) float64
	"math.Sin": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Sin parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Sin hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Sin(x.Value)}
		},
	},
	//Sinh(x float64) float64
	"math.Sinh": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Sinh parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Sinh hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Sinh(x.Value)}
		},
	},
	//Sqrt(x float64) float64
	"math.Sqrt": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Sqrt parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Sqrt hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Sqrt(x.Value)}
		},
	},
	//Tan(x float64) float64
	"math.Tan": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Tan parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Tan hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Tan(x.Value)}
		},
	},
	//Tanh(x float64) float64
	"math.Tanh": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Tanh parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Tanh hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Tanh(x.Value)}
		},
	},
	//Trunc(x float64) float64
	"math.Trunc": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Trunc parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Trunc hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Trunc(x.Value)}
		},
	},
	//Y0(x float64) float64
	"math.Y0": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Y0 parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Y0 hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Y0(x.Value)}
		},
	},
	//Y1(x float64) float64
	"math.Y1": {
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError("fungsi math.Y1 parameter sebanyak 1, didapat: %d", len(args))
			}
			x, ok := args[0].(*Float)
			if !ok {
				return NewError("fungsi math.Y1 hanya bisa menerima ANGKA, didapat: %s", args[0].Type())
			}
			return &Float{Value: math.Y1(x.Value)}
		},
	},
	////Float64frombits(b uint64) float64
	////Jn(n int, x float64) float64
	////Pow10(n int) float64
	////Yn(n int, x float64) float64
	////Float32bits(f float32) uint32
	////Float32frombits(b uint32) float32
	////Float64bits(f float64) uint64
	////Frexp(f float64) (frac float64, exp int)
	////Ilogb(x float64) int
	////IsInf(f float64, sign int) bool
	////IsNaN(f float64) (is bool)
	////Lgamma(x float64) (lgamma float64, sign int)
	////Modf(f float64) (int float64, frac float64)
	////Nextafter(x, y float64) (r float64)
	////Nextafter32(x, y float32) (r float32)
	////Signbit(x float64) bool
	////Sincos(x float64) (sin, cos float64)
	////Inf(sign int) float64
	////Ldexp(frac float64, exp int) float64
	////NaN() float64
}
