package token

const (
	PROGRAM    = "PROGRAM"
	ILLEGAL    = "ILLEGAL"
	EOF        = "EOF"
	IDENT      = "IDENT"
	INT        = "INT"
	DOT        = "."
	ASSIGN     = "="
	PLUS       = "+"
	MOD        = "%"
	COMMA      = ","
	SEMICOLON  = ";"
	FUNCTION   = "FUNCTION"
	VAR        = "VAR"
	KONST      = "KONST"
	LPAREN     = "("
	RPAREN     = ")"
	LBRACE     = "{"
	RBRACE     = "}"
	BANG       = "!"
	MINUS      = "-"
	UNDERSCORE = "_"
	SLASH      = "/"
	ASTERISK   = "*"
	LT         = "<"
	LTE        = "<="
	GT         = ">"
	GTE        = ">="
	EQ         = "=="
	NEQ        = "!="
	PIPE       = "|>"
	OR         = "||"
	AND        = "&&"
	JIKA       = "JIKA"
	ATAU       = "ATAU"
	PILAH      = "PILAH"
	PILIH      = "PILIH"
	BENAR      = "BENAR"
	SALAH      = "SALAH"
	NIHIL      = "NIHIL"
	STRING     = "STRING"
	REGEX      = "REGEX"
	LBRACKET   = "["
	RBRACKET   = "]"
	COLON      = ":"
	FATARROW   = "=>"
	ARROW      = "->"
	TILDE      = "~"
)

var (
	keywords = map[string]Type{
		"fn":    FUNCTION,
		"var":   VAR,
		"konst": KONST,
		"jika":  JIKA,
		"atau":  ATAU,
		"pilah": PILAH,
		"pilih": PILIH,
		"benar": BENAR,
		"salah": SALAH,
		"nihil": NIHIL,
	}
)

type Type string

type Token struct {
	Type    Type
	Literal string
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
