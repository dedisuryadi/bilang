package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dedisuryadi/bilang/ast"
	"github.com/dedisuryadi/bilang/lexer"
	"github.com/dedisuryadi/bilang/token"
)

const (
	_ uint8 = iota
	LOWEST
	ASSIGN
	PIPE
	FATARROW
	OR
	AND
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // fn(X)
	INDEX       // array[index]
)

var err error
var precedences = map[token.Type]uint8{
	token.ASSIGN:   ASSIGN,
	token.PIPE:     PIPE,
	token.FATARROW: FATARROW,
	token.OR:       OR,
	token.AND:      AND,
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GT:       LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.MOD:      PRODUCT,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.DOT:      CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func (p *Parser) Error() string {
	return strings.Join(p.errors, "\n")
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BENAR, p.parseBoolean)
	p.registerPrefix(token.SALAH, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.JIKA, p.parseJikaExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.PILAH, p.parsePilahLiteral)
	p.registerPrefix(token.UNDERSCORE, p.parseWildcard)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.REGEX, p.parseRegExLiteralExpression)
	p.registerPrefix(token.USAI, p.parseBreakExpression)
	p.registerPrefix(token.LANJUT, p.parseContinueExpression)
	p.registerPrefix(token.TIAP, p.parseLoopExpression)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.FATARROW, p.parseFatArrowLiteral)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.DOT, p.parseMethodCallExpression)
	p.registerInfix(token.PIPE, p.parsePipeExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)

	// read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() (prog *ast.Program, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = p
		}
	}()
	prog = &ast.Program{}
	prog.Statements = []ast.Statement{}
	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
		p.nextToken()
		if len(p.errors) > 0 {
			return nil, p
		}
	}
	return
}
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.KONST:
		return p.parseKonstStatement()
	case token.PILIH:
		return p.parsePilihStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}
func (p *Parser) peekPrecedence() uint8 {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) curPrecedence() uint8 {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}
func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be:%s got:%s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken, err = p.l.NextToken()
	if err != nil {
		p.errors = append(p.errors, err.Error())
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence uint8) ast.Expression {
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)

	}
	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	if t == token.EOF {
		return
	}
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseContinueExpression() ast.Expression {
	return &ast.ContinueExpression{Token: p.curToken}
}
func (p *Parser) parseBreakExpression() ast.Expression {
	return &ast.BreakExpression{Token: p.curToken}
}
func (p *Parser) parseLoopExpression() ast.Expression {
	curToken := p.curToken
	kv := p.parseFunctionParameters(token.DI)
	if kv == nil {
		return nil
	}

	p.nextToken()

	iter := p.parseIdentifier()
	if iter == nil {
		return nil
	}

	p.nextToken()

	body := p.parseBlockStatement()
	if body == nil {
		return nil
	}

	return &ast.LoopLiteral{
		Token: curToken,
		KV:    kv[:],
		Iter:  iter.(*ast.Identifier),
		Body:  body,
	}
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parsePilihStatement() *ast.PilihStatement {
	stmt := &ast.PilihStatement{Token: p.curToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		return nil
	}

	stmt := &ast.VarStatement{
		Token: token.Token{Type: token.VAR, Literal: token.VAR},
		Name:  &ast.Identifier{Token: ident.Token, Value: ident.TokenLiteral()},
		Value: nil,
	}

	if !p.curTokenIs(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseKonstStatement() ast.Statement {
	stmt := &ast.KonstStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.BENAR)}
}

func (p *Parser) parseWildcard() ast.Expression {
	return &ast.Wildcard{Token: p.curToken}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseJikaExpression() ast.Expression {
	expression := &ast.JikaExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ATAU) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
		if len(p.errors) > 0 {
			panic(p)
		}
	}

	return block
}

func (p *Parser) parsePilahLiteral() ast.Expression {
	lit := &ast.PilahExpression{Token: p.curToken}
	lit.Target = p.parsePilahTarget()

	conditions := make([]*ast.ExpressionStatement, 0)
	expressions := make([]ast.Expression, 0)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		cond := &ast.ExpressionStatement{Token: p.curToken}
		cond.Expression = p.parseExpression(LOWEST)
		conditions = append(conditions, cond)

		if !p.expectPeek(token.ARROW) {
			panic("must small arrow")
		}

		p.nextToken()
		expressions = append(expressions, p.parseExpression(LOWEST))
	}

	p.nextToken()
	lit.Conditions = conditions
	lit.Values = expressions

	return lit
}

func (p *Parser) parsePilahTarget() *ast.ExpressionStatement {
	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		return nil
	}

	p.nextToken()

	ident := &ast.ExpressionStatement{Token: p.curToken, Expression: p.parseExpression(LOWEST)}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	return ident

}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters(token.RPAREN)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()
	return lit
}

func (p *Parser) parseFatArrowLiteral(param ast.Expression) ast.Expression {
	ident, ok := param.(*ast.Identifier)
	if !ok || ident == nil {
		return nil
	}

	lit := &ast.FunctionLiteral{
		Token:      p.curToken,
		Parameters: []*ast.Identifier{ident},
	}

	p.nextToken()

	stmt := p.parseStatement()
	lit.Body = &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Statement{stmt},
	}

	return lit
}

func (p *Parser) parseFunctionParameters(end token.Type) []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(end) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: fn}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseRegExLiteralExpression() ast.Expression {
	return &ast.RegExLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseMethodCallExpression(obj ast.Expression) ast.Expression {
	methodCall := &ast.MethodCallExpression{Token: p.curToken, Object: obj}
	p.nextToken()

	name := p.parseIdentifier()
	if !p.peekTokenIs(token.LPAREN) {
		methodCall.Call = p.parseExpression(CALL)
	} else {
		p.nextToken()
		methodCall.Call = p.parseCallExpressions(name)
	}

	return methodCall
}
func (p *Parser) parseCallExpressions(f ast.Expression) ast.Expression {
	call := &ast.CallExpression{Token: p.curToken, Function: f}
	call.Arguments = p.parseExpressionArray(call.Arguments, token.RPAREN)
	return call
}
func (p *Parser) parseExpressionArray(a []ast.Expression, closure token.Type) []ast.Expression {
	if p.peekTokenIs(closure) {
		p.nextToken()
		return a
	}
	p.nextToken()
	a = append(a, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		a = append(a, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(closure) {
		return nil
	}
	return a
}

func (p *Parser) parsePipeExpression(left ast.Expression) ast.Expression {
	expression := &ast.Pipe{
		Token: p.curToken,
		Left:  left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}
