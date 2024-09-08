package ast

import (
	d "example/compilers/domain"
	"example/compilers/util"
)

type AstPrinter struct {
	d.ExprVisitor
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (p *AstPrinter) Print(expr d.Expr) interface{} {
	v, _ := expr.Accept(p)
	return v
}

func (p *AstPrinter) VisitBinaryExpr(expr d.BinaryExpr) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (p *AstPrinter) VisitGroupingExpr(expr d.GroupingExpr) (interface{}, error) {
	return p.parenthesize("group", expr.Expression), nil
}

func (p *AstPrinter) VisitLiteralExpr(expr d.LiteralExpr) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return util.ToString(expr.Value), nil
}

func (p *AstPrinter) VisitUnaryExpr(expr d.UnaryExpr) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (p *AstPrinter) parenthesize(name string, exprs ...d.Expr) string {
	ret := ""

	ret += "("
	ret += name
	for _, expr := range exprs {
		ret += " "
		v, _ := expr.Accept(p)
		ret += util.ToString(v)
	}
	ret += ")"

	return ret
}
