package ast

import (
	d "example/compilers/domain"
	"example/compilers/util"
	"fmt"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

var _ d.ExprVisitor = (*AstPrinter)(nil)

func (p *AstPrinter) Print(expr d.Expr) interface{} {
	v, _ := expr.Accept(p)
	return v
}

func (p *AstPrinter) VisitSuperExpr(expr d.SuperExpr) (interface{}, error) {
	return "Super: " + expr.Keyword.Lexeme + expr.Method.Lexeme, nil
}

func (p *AstPrinter) VisitGetExpr(expr d.GetExpr) (interface{}, error) {
	return p.parenthesize("Get: "+expr.Name.Lexeme, expr.Object), nil
}

func (p *AstPrinter) VisitSetExpr(expr d.SetExpr) (interface{}, error) {
	return p.parenthesize("Set: "+expr.Name.Lexeme, expr.Object), nil
}

func (p *AstPrinter) VisitThisExpr(expr d.ThisExpr) (interface{}, error) {
	return "This: " + expr.Keyword.Lexeme, nil
}

func (p *AstPrinter) VisitAssignExpr(expr d.AssignExpr) (interface{}, error) {
	v, err := expr.Value.Accept(p)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("ASSIGN{%s}", v), nil
}

func (p *AstPrinter) VisitLogicalExpr(expr d.LogicalExpr) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
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

func (p *AstPrinter) VisitCallExpr(expr d.CallExpr) (interface{}, error) {
	callee, err := expr.Callee.Accept(p)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("CALL{%s, %s, %s}", callee, expr.Paren.Lexeme, p.parenthesize("", expr.Args...)), nil
}

func (p *AstPrinter) VisitVariableExpr(expr d.VariableExpr) (interface{}, error) {
	return fmt.Sprintf("VAR{%s}", expr.Name.Lexeme), nil
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
