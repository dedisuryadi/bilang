package evaluator

import (
	"reflect"
	"testing"

	"github.com/dedisuryadi/bilang/lexer"
	"github.com/dedisuryadi/bilang/object"
	"github.com/dedisuryadi/bilang/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"100 % 30 % 4 + 140 % 100", 42},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	script := NewScript()
	return script.Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"benar", true},
		{"salah", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 <= 2", true},
		{"1 >= 2", false},
		{"1 <= 1", true},
		{"1 >= 1", true},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"1 >= 1", true},
		{"1 >= 1", true},
		{"1 >= 2", false},
		{"1 >= 2", false},
		{"1 <= 1", true},
		{"1 <= 1", true},
		{"1 <= 2", true},
		{"1 <= 2", true},
		{"benar == benar", true},
		{"salah == salah", true},
		{"benar == salah", false},
		{"benar != salah", true},
		{"salah != benar", true},
		{"benar || benar", true},
		{"salah || salah", false},
		{"benar || salah", true},
		{"salah || benar", true},
		{"benar && benar", true},
		{"salah && salah", false},
		{"benar && salah", false},
		{"salah && benar", false},
		{"(1 < 2) == benar", true},
		{"(1 < 2) == salah", false},
		{"(1 > 2) == benar", false},
		{"(1 > 2) == salah", true},
		{"(1 <= 2) == benar", true},
		{"(1 <= 2) == salah", false},
		{"(1 >= 2) == benar", false},
		{"(1 >= 2) == salah", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected, tt.input)
	}
}
func testBooleanObject(t *testing.T, obj object.Object, expected bool, input string) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t input=%s", result.Value, expected, input)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!benar", false},
		{"!salah", true},
		{"!5", false},
		{"!!benar", true},
		{"!!salah", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`jika ("foo") { 10 }`, nil},
		{"jika (benar) { 10 }", 10},
		{"jika (salah) { 10 }", nil},
		{"jika (1) { 10 }", 10},
		{"jika (0) { 10 }", nil},
		{"jika (-1) { 10 }", nil},
		{"jika (1 < 2) { 10 }", 10},
		{"jika (1 > 2) { 10 }", nil},
		{"jika (1 > 2) { 10 } atau { 20 }", 20},
		{"jika (1 < 2) { 10 } atau { 20 }", 10},
		{"jika (1 <= 2) { 10 }", 10},
		{"jika (1 >= 2) { 10 }", nil},
		{"jika (1 >= 2) { 10 } atau { 20 }", 20},
		{"jika (1 <= 2) { 10 } atau { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"pilih 10;", 10},
		{"pilih 10; 9;", 10},
		{"pilih 2 * 5; 9;", 10},
		{"9; pilih 2 * 5; 9;", 10},
		{`
jika (10 > 1) {
	jika (10 > 1) {
		pilih 10;
	}
	pilih 1;
}`, 10,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			"5 + benar;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + benar; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-benar",
			"unknown operator: -BOOLEAN",
		},
		{
			"benar + salah;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; benar + salah; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"jika (10 > 1) { benar + salah; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
jika (10 > 1) {
jika (10 > 1) {
pilih benar + salah;
}
pilih 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		wantErr  bool
	}{
		{input: "var a = 5; a = \"foo\"; a;", wantErr: true},
		{input: "var a = 5; a = 10+5-5; a;", expected: 10},
		{input: "var a = 5; var a = 10; a;", expected: 10},
		{input: "var a = 5; a;", expected: 5},
		{input: "var a = 5 * 5; a;", expected: 25},
		{input: "var a = 5; var b = a; b;", expected: 5},
		{input: "var a = 5; var b = a; var c = a + b + 5; c;", expected: 15},
	}
	for _, tt := range tests {
		res := testEval(tt.input)
		if !tt.wantErr {
			testIntegerObject(t, res, tt.expected)
			continue
		}
		_, ok := res.(*object.Error)
		if !ok {
			t.Errorf("expected %v got %v", &object.Error{}, res)
		}
	}
}

func TestKonstStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		wantErr  bool
	}{
		{input: "konst a = 5; a = 10; a;", wantErr: true},
		{input: "konst a = 5; konst a = 10; a;", wantErr: true},
		{input: "konst a = 5; a;", expected: 5},
		{input: "konst a = 5 * 5; a;", expected: 25},
		{input: "konst a = 5; konst b = a; b;", expected: 5},
		{input: "konst a = 5; konst b = a; konst c = a + b + 5; c;", expected: 15},
	}
	for _, tt := range tests {
		res := testEval(tt.input)
		if !tt.wantErr {
			testIntegerObject(t, res, tt.expected)
			continue
		}
		_, ok := res.(*object.Error)
		if !ok {
			t.Errorf("expected %v got %v", &object.Error{}, res)
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var identity = fn(x) { x; }; identity(5);", 5},
		{"var identity = fn(x) { pilih x; }; identity(5);", 5},
		{"var double = fn(x) { x * 2; }; double(5);", 10},
		{"var add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"var add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},

		{"var identity = x => x; identity(5);", 5},
		{"var identity = x => pilih x; identity(5);", 5},
		{"var double = x => x * 2; double(5);", 10},
		{"(x => x)(5)", 5},
		{"var addTo = x => y => x+y; var addFive = addTo(5); addFive(0);", 5},
		{"var addTo = x => y => x+y; var addFive = addTo(5); addFive(10);", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestPilahExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			input: `
				var x = 6;
				pilah x {
					5 -> "lima"
					6 -> "enam"
					_ -> "lainnya"
				}
			`,
			expected: &object.String{Value: "enam"},
		},
		{
			input: `
				var x = 6;
				var y = pilah x {
					_ -> "lainnya"
					5 -> "lima"
					6 -> "enam"
				}
				y
			`,
			expected: &object.String{Value: "lainnya"},
		},
	}
	for _, tt := range tests {
		got := testEval(tt.input)
		if !reflect.DeepEqual(tt.expected, got) {
			t.Errorf("TestPilahExpression expected=%v got=%v", tt.expected, got)
		}
	}
}

func TestPipeExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			input: `
				var adder = x => y => x+y;
				var addOne = adder(1);
				var double = x => x*2;
				var grade = fn(x) {
					pilah x {
						0 -> "e"
						1 -> "c"
						2 -> "b"
						3 -> "a"
						_ -> "lainnya"
					}
				}
				
				var z = 0 
					|> addOne 
					|> double 
					|> grade
					;
				
				z;
			`,
			expected: &object.String{Value: "b"},
		},
	}
	for _, tt := range tests {
		got := testEval(tt.input)
		if !reflect.DeepEqual(tt.expected, got) {
			t.Errorf("TestPipeExpression expected=%v got=%v", tt.expected, got)
		}
	}
}

func TestClosures(t *testing.T) {
	input := `
var newAdder = fn(x) {
fn(y) { x + y };
};
var addTwo = newAdder(2);
addTwo(2);`
	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`panjang("")`, 0},
		{`panjang("four")`, 4},
		{`panjang("hello world")`, 11},
		{`panjang(1)`, "argument to `panjang` not supported, got INTEGER"},
		{`panjang("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"var i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"var myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"var myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"var myArray = [1, 2, 3]; var i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
		{
			`var reduce = fn(arr, init, f) {
						var iter = fn(arr, hasil) {
							jika (panjang(arr) == 0) { pilih hasil };
							iter(ekor(arr), f(hasil, awal(arr)))
						}
						iter(arr, init)
					}
					
					var sum = arr => reduce(arr, 0, fn(init, nilai){ init+nilai })	
					sum([1,2,3,4,5])
					`,
			15,
		},
		{
			`var reduce = fn(arr, init, f) {
						var iter = fn(arr, hasil) {
							jika (panjang(arr) == 0) { pilih hasil };
							iter(ekor(arr), f(hasil, arr |> awal))
						}
						iter(arr, init)
					}
					
					var sum = arr => reduce(arr, 0, fn(init, nilai){ init+nilai })
					var map = fn(arr, f) {
						var iter = fn(arr, akum) {
							jika (panjang(arr) == 0) { pilih akum }
							var hasil = push(akum, arr |> awal |> f)
							iter(arr |> ekor, hasil)
						}
						iter(arr, [])
					}
					
					var a = [1,2,3,4,5]
					var ganda = x => x*2
					map(a, ganda) |> sum
					`,
			30,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `var two = "two";
{
"one": 10 - 9,
two: 1 + 1,
"thr" + "ee": 6 / 2,
4: 4,
benar: 5,
salah: 6
}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`var key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{benar: 5}[benar]`,
			5,
		},
		{
			`{salah: 5}[salah]`,
			5,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
