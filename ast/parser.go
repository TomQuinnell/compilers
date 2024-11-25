package ast

import (
	"errors"
	d "example/compilers/domain"
	"fmt"
)

const (
	maxArgsSize = 255
)

type ErrParse struct {
	message string
	token   *d.Token
}

func (e ErrParse) Error() string {
	return fmt.Sprintf("%s %s", e.message, e.token)
}

type Parser struct {
	tokens  []*d.Token
	current int
}

func NewParser(tokens []*d.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ([]d.Stmt, error) {
	statements := make([]d.Stmt, 0)
	for !p.isAtEnd() {
		st, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, st)
	}

	return statements, nil
}

func (p *Parser) parseDeclaration() (d.Stmt, error) {
	pFunc := p.parseStatement
	if p.match(d.CLASS) {
		return p.parseClassDeclaration()
	}
	if p.match(d.FUN) {
		pFunc = func() (d.Stmt, error) {
			return p.parseFunction("function")()
		}
	}
	if p.match(d.VAR) {
		pFunc = p.parseVarDeclaration
	}

	s, err := pFunc()
	if err == nil {
		return s, nil
	}

	if !errors.Is(err, ErrParse{}) {
		return s, err
	}

	p.sync()
	return nil, nil
}

func (p *Parser) parseClassDeclaration() (d.Stmt, error) {
	name, err := p.consume(d.IDENTIFIER, "Expect class name")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(d.LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	methods := make([]d.FunctionStmt, 0)
	for !p.check(d.RIGHT_BRACE) && !p.isAtEnd() {
		parsedFn, err := p.parseFunction("method")()
		if err != nil {
			return nil, err
		}

		methods = append(methods, parsedFn)
	}

	_, err = p.consume(d.RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}

	return d.ClassStmt{
		Name:    name,
		Methods: methods,
	}, nil
}

func (p *Parser) parseFunction(kind string) func() (d.FunctionStmt, error) {
	nilFn := d.FunctionStmt{}
	return func() (d.FunctionStmt, error) {
		name, err := p.consume(d.IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
		if err != nil {
			return nilFn, err
		}

		_, err = p.consume(d.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))
		if err != nil {
			return nilFn, err
		}

		params := make([]*d.Token, 0)
		if !p.check(d.RIGHT_PAREN) {
			paramV, err := p.consume(d.IDENTIFIER, "Expect parameter name")
			if err != nil {
				return nilFn, err
			}
			params = append(params, paramV)

			for p.match(d.COMMA) {
				if len(params) >= maxArgsSize {
					return nilFn, ErrParse{
						message: fmt.Sprintf("Can't have more than %d params.", maxArgsSize),
						token:   p.peek(),
					}
				}

				paramV, err := p.consume(d.IDENTIFIER, "Expect parameter name")
				if err != nil {
					return nilFn, err
				}
				params = append(params, paramV)
			}
		}

		_, err = p.consume(d.RIGHT_PAREN, "Expect ')' after paramters.")
		if err != nil {
			return nilFn, err
		}

		_, err = p.consume(d.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
		if err != nil {
			return nilFn, err
		}

		body, err := p.parseBlock()
		if err != nil {
			return nilFn, err
		}

		return d.FunctionStmt{
			Name:   name,
			Params: params,
			Body:   body,
		}, nil
	}
}

func (p *Parser) parseVarDeclaration() (d.Stmt, error) {
	name, err := p.consume(d.IDENTIFIER, "Expect var name")
	if err != nil {
		return nil, err
	}

	var init d.Expr
	if p.match(d.EQUAL) {
		init, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(d.SEMICOLON, "Expect ';' after var declaration")
	if err != nil {
		return nil, err
	}

	return d.VarStmt{
		Name:        name,
		Initializer: init,
	}, nil
}

func (p *Parser) parseStatement() (d.Stmt, error) {
	if p.match(d.FOR) {
		return p.parseForStmt()
	}
	if p.match(d.IF) {
		return p.parseIfStmt()
	}
	if p.match(d.PRINT) {
		return p.parsePrintStatement()
	}
	if p.match(d.RETURN) {
		return p.parseReturnStatement()
	}
	if p.match(d.WHILE) {
		return p.parseWhileStatement()
	}
	if p.match(d.LEFT_BRACE) {
		stmts, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		return d.BlockStmt{
			Stmts: stmts,
		}, nil
	}

	return p.parseExpressionStmt()
}

func (p *Parser) parseForStmt() (d.Stmt, error) {
	// Consume tokens
	_, err := p.consume(d.LEFT_PAREN, "Expect '(' after 'for'")
	if err != nil {
		return nil, err
	}

	var initializer d.Stmt
	if p.match(d.SEMICOLON) {
		initializer = nil
	} else if p.match(d.VAR) {
		initializer, err = p.parseVarDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.parseExpressionStmt()
		if err != nil {
			return nil, err
		}
	}

	var condition d.Expr
	if !p.check(d.SEMICOLON) {
		condition, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	} else {
		condition = d.LiteralExpr{Value: true}
	}
	_, err = p.consume(d.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment d.Expr
	if !p.check(d.RIGHT_PAREN) {
		increment, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(d.RIGHT_PAREN, "Expect ')' after for clause.")
	if err != nil {
		return nil, err
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	// Sugarfy
	if increment != nil {
		body = d.BlockStmt{
			Stmts: []d.Stmt{body, d.ExpressionStmt{Expression: increment}},
		}
	}
	body = d.WhileStmt{
		Condition: condition,
		Body:      body,
	}
	if initializer != nil {
		body = d.BlockStmt{
			Stmts: []d.Stmt{initializer, body},
		}
	}

	return body, nil
}

func (p *Parser) parseIfStmt() (d.Stmt, error) {
	_, err := p.consume(d.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(d.RIGHT_PAREN, "Expect ')' after 'if' condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	var elseBranch d.Stmt
	if p.match(d.ELSE) {
		elseBranch, err = p.parseStatement()
		if err != nil {
			return nil, err
		}
	}

	return d.IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) parseWhileStatement() (d.Stmt, error) {
	_, err := p.consume(d.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(d.RIGHT_PAREN, "Expect ')' after 'while' condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return d.WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) parsePrintStatement() (d.Stmt, error) {
	ex, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(d.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return d.PrintStmt{
		Expression: ex,
	}, nil
}

func (p *Parser) parseReturnStatement() (d.Stmt, error) {
	keyword := p.previous()

	var value d.Expr
	if !p.check(d.SEMICOLON) {
		var err error
		value, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	_, err := p.consume(d.SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return d.ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *Parser) parseBlock() ([]d.Stmt, error) {
	stmts := make([]d.Stmt, 0)

	for !p.check(d.RIGHT_BRACE) && !p.isAtEnd() {
		s, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, s)
	}

	_, err := p.consume(d.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return stmts, nil
}

func (p *Parser) parseExpressionStmt() (d.Stmt, error) {
	ex, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(d.SEMICOLON, "Expect ';' after statement.")
	if err != nil {
		return nil, err
	}

	return d.ExpressionStmt{
		Expression: ex,
	}, nil
}

func (p *Parser) parseExpression() (d.Expr, error) {
	return p.parseAssignment()
}

func (p *Parser) parseAssignment() (d.Expr, error) {
	eqExpr, err := p.parseOr()
	if err != nil {
		return nil, err
	}

	if p.match(d.EQUAL) {
		eqToken := p.previous()
		value, err := p.parseAssignment()
		if err != nil {
			return nil, err
		}

		switch eqExprRaw := eqExpr.(type) {
		case d.VariableExpr:
			return d.AssignExpr{
				Name:  eqExprRaw.Name,
				Value: value,
			}, nil
		case d.GetExpr:
			return d.SetExpr{
				Object: eqExprRaw.Object,
				Name:   eqExprRaw.Name,
				Value:  value,
			}, nil
		}

		return nil, ErrParse{message: "Invalid assingment target.", token: eqToken}
	}

	return eqExpr, nil
}

func (p *Parser) parseOr() (d.Expr, error) {
	expr, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.match(d.OR) {
		operator := p.previous()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}

		expr = d.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) parseAnd() (d.Expr, error) {
	expr, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	for p.match(d.AND) {
		operator := p.previous()
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}

		expr = d.LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) parseEquality() (d.Expr, error) {
	expr, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.match(d.BANG_EQUAL, d.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		expr = d.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) parseComparison() (d.Expr, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for p.match(d.GREATER, d.GREATER_EQUAL, d.LESS, d.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		expr = d.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) parseTerm() (d.Expr, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.match(d.MINUS, d.PLUS) {
		operator := p.previous()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		expr = d.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) parseFactor() (d.Expr, error) {
	expr, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.match(d.SLASH, d.STAR) {
		operator := p.previous()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		expr = d.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) parseUnary() (d.Expr, error) {
	if p.match(d.BANG, d.MINUS) {
		operator := p.previous()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return d.UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.parseCall()
}

func (p *Parser) parseCall() (d.Expr, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(d.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(d.DOT) {
			name, err := p.consume(d.IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = d.GetExpr{
				Object: expr,
				Name:   name,
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee d.Expr) (d.Expr, error) {
	args := make([]d.Expr, 0)

	if !p.check(d.RIGHT_PAREN) {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		for p.match(d.COMMA) {
			if len(args) >= maxArgsSize {
				return nil, ErrParse{
					message: fmt.Sprintf("Can't have more than %d args.", maxArgsSize),
					token:   p.peek(),
				}
			}

			expr, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			args = append(args, expr)
		}
	}

	parenToken, err := p.consume(d.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return d.CallExpr{
		Callee: callee,
		Paren:  parenToken,
		Args:   args,
	}, nil
}

func (p *Parser) parsePrimary() (d.Expr, error) {
	if p.match(d.FALSE) {
		return d.LiteralExpr{
			Value: false,
		}, nil
	}
	if p.match(d.TRUE) {
		return d.LiteralExpr{
			Value: true,
		}, nil
	}
	if p.match(d.NIL) {
		return d.LiteralExpr{
			Value: nil,
		}, nil
	}

	if p.match(d.NUMBER, d.STRING) {
		return d.LiteralExpr{
			Value: p.previous().Literal,
		}, nil
	}

	if p.match(d.THIS) {
		return d.ThisExpr{
			Keyword: p.previous(),
		}, nil
	}

	if p.match(d.IDENTIFIER) {
		return d.VariableExpr{
			Name: p.previous(),
		}, nil
	}

	if p.match(d.LEFT_PAREN) {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(d.RIGHT_PAREN, "Expect closing ')' after expression.")
		if err != nil {
			return nil, err
		}
		return d.GroupingExpr{
			Expression: expr,
		}, nil
	}

	return nil, ErrParse{message: "Expected expression.", token: p.peek()}
}

func (p *Parser) consume(t d.TokenType, message string) (*d.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return nil, ErrParse{message: message, token: p.peek()}
}

func (p *Parser) sync() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Kind == d.SEMICOLON {
			return
		}

		switch p.peek().Kind {
		case d.CLASS:
		case d.FUN:
		case d.VAR:
		case d.FOR:
		case d.IF:
		case d.WHILE:
		case d.PRINT:
		case d.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) match(types ...d.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t d.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Kind == t
}

func (p *Parser) advance() *d.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	next := p.peek()
	return next == nil || next.Kind == d.EOF
}

func (p *Parser) peek() *d.Token {
	if p.current >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.current]
}

func (p *Parser) previous() *d.Token {
	return p.tokens[p.current-1]
}
