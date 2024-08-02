package lox

// program        → declaration* EOF ;
// declaration    → varDecl
//                | statement ;
// varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
// statement      → exprStmt
//                | printStmt
//                | block ;
// block          → "{" declaration* "}" ;
// exprStmt       → expression ";" ;
// printStmt      → "print" expression ";" ;
// expression     → assignment ;
// assignment     → IDENTIFIER "=" assignment
//                | equality ;
// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           → factor ( ( "-" | "+" ) factor )* ;
// factor         → unary ( ( "/" | "*" ) unary )* ;
// unary          → ( "!" | "-" ) unary
//                | primary ;
// primary        → "true" | "false" | "nil"
//                | NUMBER | STRING
//                | "(" expression ")"
//                | IDENTIFIER ;

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt

	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.synchronize()
			return nil, err
		}

		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after variable declaration."); err != nil {
		return nil, err
	}

	return NewVarDeclStmt(name, initializer), nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(LEFT_BRACE) {
		return p.block()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() (Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after value."); err != nil {
		return nil, err
	}

	return NewPrintStmt(value), nil
}

func (p *Parser) block() (Stmt, error) {
	var statements []Stmt

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	if _, err := p.consume(RIGHT_BRACE, "Expect '}' after block."); err != nil {
		return nil, err
	}

	return NewBlockStmt(statements), nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after value."); err != nil {
		return nil, err
	}

	return NewExprStmt(expr), nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if v, ok := expr.(*VariableExpr); ok {
			name := v.Name
			return NewAssignExpr(name, value), nil
		}

		return nil, NewParseError(equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return NewUnaryExpr(operator, right), nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(NUMBER, STRING) {
		return NewLiteralExpr(p.previous().Literal), nil
	}

	if p.match(FALSE) {
		return NewLiteralExpr(NewLiteral(false)), nil
	}

	if p.match(TRUE) {
		return NewLiteralExpr(NewLiteral(true)), nil
	}

	if p.match(NIL) {
		return NewLiteralExpr(NewLiteral(nil)), nil
	}

	if p.match(IDENTIFIER) {
		return NewVariableExpr(p.previous()), nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(RIGHT_PAREN, "Unmatched parentheses."); err != nil {
			return nil, err
		}

		return NewGroupingExpr(expr), nil
	}

	return nil, NewParseError(p.peek(), "Expect expression.")
}

func (p *Parser) match(tokens ...TokenType) bool {
	for _, token := range tokens {
		if p.check(token) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return Token{}, NewParseError(p.peek(), message)
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}
