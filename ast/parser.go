package ast

import (
	d "example/compilers/domain"
	"fmt"
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

func (p *Parser) Parse() (d.Expr, error) {
	return p.parseExpression()
}

func (p *Parser) parseExpression() (d.Expr, error) {
	return p.parseEquality()
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

	return p.parsePrimary()
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
