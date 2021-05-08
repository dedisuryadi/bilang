package lexer

import (
	"testing"

	"github.com/dedisuryadi/bilang/token"
)

func TestNextToken(t *testing.T) {
	input :=
		`var five = 5;
		var ten;
		ten = 10;		

		konst two = 2;

		konst add = a, b => a + b;
		konst addTo = a => b => a + b;
		
		a |> b;
		a || b;
		a && b;
		a % b;
		nihil;

		/[a-z]+/
		~[0-9]+~

		function.call

		var a = pilah x,y {
			isEven, isOdd -> "x is even and y is odd"
			isOdd, isEven -> "y is even and x is odd"
			_, _ -> "both either even or odd" 
		}
		
		var add = fn(x, y) {
		  x + y;
		};
		
		var result = add(five, ten);
		!-*5;
		5 < 10 > 5;
		5 <= 10 >= 5;
		
		jika (5 < 10) {
			pilih benar;
		} atau {
			pilih salah;
		}
		
		10 == 10;
		10 != 9;
		
		"foobar"
		"foo bar"
		"foo\"bar"

		[1,2];
		{"foo": "bar"}
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "ten"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.KONST, "konst"},
		{token.IDENT, "two"},
		{token.ASSIGN, "="},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},

		{token.KONST, "konst"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.IDENT, "a"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.FATARROW, "=>"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},

		{token.KONST, "konst"},
		{token.IDENT, "addTo"},
		{token.ASSIGN, "="},
		{token.IDENT, "a"},
		{token.FATARROW, "=>"},
		{token.IDENT, "b"},
		{token.FATARROW, "=>"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a"},
		{token.PIPE, "|>"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a"},
		{token.OR, "||"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a"},
		{token.MOD, "%"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},

		{token.NIHIL, "nihil"},
		{token.SEMICOLON, ";"},

		{token.REGEX, `[a-z]+`},
		{token.REGEX, `[0-9]+`},

		{token.IDENT, "function"},
		{token.DOT, "."},
		{token.IDENT, "call"},

		{token.VAR, "var"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.PILAH, "pilah"},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.LBRACE, "{"},
		{token.IDENT, "isEven"},
		{token.COMMA, ","},
		{token.IDENT, "isOdd"},
		{token.ARROW, "->"},
		{token.STRING, "x is even and y is odd"},
		{token.IDENT, "isOdd"},
		{token.COMMA, ","},
		{token.IDENT, "isEven"},
		{token.ARROW, "->"},
		{token.STRING, "y is even and x is odd"},
		{token.UNDERSCORE, "_"},
		{token.COMMA, ","},
		{token.UNDERSCORE, "_"},
		{token.ARROW, "->"},
		{token.STRING, "both either even or odd"},
		{token.RBRACE, "}"},

		{token.VAR, "var"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LTE, "<="},
		{token.INT, "10"},
		{token.GTE, ">="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.JIKA, "jika"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.PILIH, "pilih"},
		{token.BENAR, "benar"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ATAU, "atau"},
		{token.LBRACE, "{"},
		{token.PILIH, "pilih"},
		{token.SALAH, "salah"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NEQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.STRING, `foo\"bar`},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatal(err)
		}
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
