package lexer

import (
	"errors"
	"fmt"

	"github.com/dedisuryadi/bilang/token"
)

type Lexer struct {
	ch           byte
	position     int
	readPosition int
	col          int
	line         int
	input        string
	prev         token.Token
}

func New(input string) *Lexer {
	l := &Lexer{input: input, col: 1, line: 1}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() (token.Token, error) {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.FATARROW, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '%':
		tok = newToken(token.MOD, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '"':
		str, err := l.readString()
		if err != nil {
			return tok, err
		}
		tok.Literal = str
		tok.Type = token.STRING
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.ARROW, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.MINUS, l.ch)

		}

	case '/':
		if l.prev.Type == token.RPAREN || // (a+c) / b
			l.prev.Type == token.RBRACKET || // a[3] / b
			l.prev.Type == token.IDENT || // a / b
			l.prev.Type == token.INT { // 3 / b
			tok = newToken(token.SLASH, l.ch)
		} else {
			//regexp
			tok.Literal = l.readRegex('/')
			tok.Type = token.REGEX
		}
	case '~':
		tok.Literal = l.readRegex('~')
		tok.Type = token.REGEX

	case '.':
		tok = newToken(token.DOT, l.ch)

	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := string(l.ch)
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: ch + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := string(l.ch)
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: ch + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '_':
		tok = newToken(token.UNDERSCORE, l.ch)

	case '|':
		if l.peekChar() == '>' {
			tok = token.Token{Type: token.PIPE, Literal: string(l.ch) + string(l.peekChar())}
			l.readChar()
		} else if l.peekChar() == '|' {
			tok = token.Token{Type: token.OR, Literal: string(l.ch) + string(l.peekChar())}
			l.readChar()
		}
	//} else {
	//	tok = newToken(token.BITOR, l.ch)
	//}

	case '&':
		if l.peekChar() != '&' {
			panic("unsupported infix operator")
		}
		tok = token.Token{Type: token.AND, Literal: string(l.ch) + string(l.peekChar())}
		l.readChar()

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			tok.Literal = ident
			tok.Type = token.LookupIdent(ident)
			l.prev = tok
			return tok, nil
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			l.prev = tok
			return tok, nil
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
			return tok, fmt.Errorf("illegal syntax on %v token %s", l.position, tok.Literal)
		}
	}

	tok.Col = l.col
	tok.Line = l.line

	l.readChar()

	l.prev = tok
	return tok, nil
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII for NUL
	} else {
		l.ch = l.input[l.readPosition]
		if l.ch == '\n' {
			l.col = 1
			l.line += 1
		} else {
			l.col += 1
		}
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() (string, error) {
	var ret []byte
	for {
		l.readChar()
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'a':
				ret = append(ret, '\a')
				continue
			case 'b':
				ret = append(ret, '\b')
				continue
			case 'f':
				ret = append(ret, '\f')
				continue
			case 'r':
				ret = append(ret, '\r')
				continue
			case 'n':
				ret = append(ret, '\n')
				continue
			case 't':
				ret = append(ret, '\t')
				continue
			case 'v':
				ret = append(ret, '\v')
				continue
			}
			ret = append(ret, '\\', l.ch)
			continue
		} else if l.ch == '"' {
			break
		} else if len(ret) > len(l.input) {
			return "", errors.New("syntax error")
		} else {
			ret = append(ret, l.ch)
		}
	}
	return string(ret), nil
}

func (l *Lexer) readRegex(delim byte) (lit string) {
	start := l.position + 1
	for {
		l.readChar()
		if l.ch == '\\' {
			// skip escape sequence
			l.readChar()
		} else if l.ch == delim {
			lit = l.input[start:l.position]
			l.readChar() // skip the closing delim
			return
		}
	}
}

func newToken(Type token.Type, ch byte) token.Token {
	return token.Token{Type: Type, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
