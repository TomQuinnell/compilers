package ast

import (
	d "example/compilers/domain"
)

type AstPrinter struct {
	d.ExprVisitor[string]
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (p *AstPrinter) Print(expr d.Expr[string]) string {
	return expr.Accept(p)
}

func (p *AstPrinter) VisitBinaryExpr(expr d.BinaryExpr[string]) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr d.GroupingExpr[string]) string {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr d.LiteralExpr[string]) string {
	if expr.Value == nil {
		return "nil"
	}
	return expr.Value.String()
}

func (p *AstPrinter) VisitUnaryExpr(expr d.UnaryExpr[string]) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) parenthesize(name string, exprs ...d.Expr[string]) string {
	ret := ""

	ret += "("
	ret += name
	for _, expr := range exprs {
		ret += " "
		ret += expr.Accept(p)
	}
	ret += ")"

	return ret
}
