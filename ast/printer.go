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
	return expr.Accept(p)
}

func (p *AstPrinter) VisitBinaryExpr(expr d.BinaryExpr) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr d.GroupingExpr) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr d.LiteralExpr) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	return util.ToString(expr.Value)
}

func (p *AstPrinter) VisitUnaryExpr(expr d.UnaryExpr) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) parenthesize(name string, exprs ...d.Expr) string {
	ret := ""

	ret += "("
	ret += name
	for _, expr := range exprs {
		ret += " "
		ret += util.ToString(expr.Accept(p))
	}
	ret += ")"

	return ret
}
