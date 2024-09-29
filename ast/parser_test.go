package ast

import (
	d "example/compilers/domain"
	"example/compilers/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type ParseTestCase struct {
		rawTokens    []*d.Token
		expectedExpr d.Expr
	}

	openBracket := d.NewToken(d.LEFT_PAREN, "(", nil, 1)
	closeBracket := d.NewToken(d.RIGHT_PAREN, ")", nil, 1)
	one := d.NewToken(d.NUMBER, "1", 1, 1)
	a := d.NewToken(d.STRING, "\"a\"", "a", 1)
	eqeq := d.NewToken(d.EQUAL_EQUAL, "==", nil, 1)
	neqeq := d.NewToken(d.BANG_EQUAL, "!=", nil, 1)
	gt := d.NewToken(d.GREATER, ">", nil, 1)
	gteq := d.NewToken(d.GREATER_EQUAL, ">=", nil, 1)
	lt := d.NewToken(d.LESS, "<", nil, 1)
	lteq := d.NewToken(d.LESS_EQUAL, "<=", nil, 1)
	min := d.NewToken(d.MINUS, "-", nil, 1)
	plus := d.NewToken(d.PLUS, "+", nil, 1)
	mult := d.NewToken(d.STAR, "*", nil, 1)
	div := d.NewToken(d.SLASH, "/", nil, 1)
	bang := d.NewToken(d.BANG, "/", nil, 1)
	falso := d.NewToken(d.FALSE, "false", false, 1)
	trutho := d.NewToken(d.TRUE, "true", true, 1)
	nilo := d.NewToken(d.NIL, "nil", nil, 1)

	testCases := []ParseTestCase{
		{[]*d.Token{one}, d.LiteralExpr{Value: 1}},
		{[]*d.Token{a}, d.LiteralExpr{Value: "a"}},
		{[]*d.Token{falso}, d.LiteralExpr{Value: false}},
		{[]*d.Token{trutho}, d.LiteralExpr{Value: true}},
		{[]*d.Token{nilo}, d.LiteralExpr{Value: nil}},
		{[]*d.Token{one, eqeq, one}, d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: eqeq, Right: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{one, neqeq, one}, d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: neqeq, Right: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{one, gt, one}, d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: gt, Right: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{one, gteq, one}, d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: gteq, Right: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{one, lt, one}, d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: lt, Right: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{one, lteq, one}, d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: lteq, Right: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{one, min, one}, d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: min, Right: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{one, plus, one}, d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: plus, Right: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{one, mult, one}, d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: mult, Right: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{one, div, one}, d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: div, Right: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{min, one}, d.UnaryExpr{Operator: min, Right: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{bang, falso}, d.UnaryExpr{Operator: bang, Right: d.LiteralExpr{Value: false}}},
		{[]*d.Token{one, mult, one, plus, one, eqeq, bang, bang, falso}, d.BinaryExpr{
			Left: d.BinaryExpr{
				Left: d.BinaryExpr{
					Left: d.LiteralExpr{Value: 1}, Operator: mult, Right: d.LiteralExpr{Value: 1},
				},
				Operator: plus,
				Right:    d.LiteralExpr{Value: 1},
			},
			Operator: eqeq,
			Right: d.UnaryExpr{
				Operator: bang,
				Right: d.UnaryExpr{
					Operator: bang,
					Right:    d.LiteralExpr{Value: false},
				},
			},
		}},
		{[]*d.Token{openBracket, one, eqeq, one, closeBracket}, d.GroupingExpr{
			Expression: d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: eqeq, Right: d.LiteralExpr{Value: 1}}},
		},
	}

	for _, c := range testCases {
		t.Run(fmt.Sprintf("Parses raw expr tokens: %s", util.SprintTokens(c.rawTokens)), func(t *testing.T) {
			assert := assert.New(t)

			stmts, err := NewParser(c.rawTokens).Parse()
			assert.NoError(err)

			assert.Len(stmts, 1)
			st := stmts[0]

			var expr d.Expr
			switch s := st.(type) {
			case d.ExpressionStmt:
				expr = s.Expression
			default:
				assert.Fail("Expected expression stmt")
			}

			if !util.IsEqualExpr(c.expectedExpr, expr) {
				fmt.Println("Expected:")
				fmt.Println(NewAstPrinter().Print(c.expectedExpr))
				fmt.Println("Actual:")
				fmt.Println(NewAstPrinter().Print(expr))
			}
			assert.True(util.IsEqualExpr(c.expectedExpr, expr))
		})
	}

	semicolon := d.NewToken(d.SEMICOLON, ";", nil, 0)

	t.Run("Parses 2 raw exprs", func(t *testing.T) {
		assert := assert.New(t)

		rawTokens := []*d.Token{one, eqeq, one, semicolon, one, eqeq, one}

		stmts, err := NewParser(rawTokens).Parse()
		assert.NoError(err)

		assert.Len(stmts, 2)

		for _, st := range stmts {
			var expr d.Expr
			switch s := st.(type) {
			case d.ExpressionStmt:
				expr = s.Expression
			default:
				assert.Fail("Expected expression stmt")
			}

			expectedExpr := d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: eqeq, Right: d.LiteralExpr{Value: 1}}
			if !util.IsEqualExpr(expectedExpr, expr) {
				fmt.Println("Expected:")
				fmt.Println(NewAstPrinter().Print(expectedExpr))
				fmt.Println("Actual:")
				fmt.Println(NewAstPrinter().Print(expr))
			}
			assert.True(util.IsEqualExpr(expectedExpr, expr))
		}
	})

	type ParseStmtTestCase struct {
		rawTokens    []*d.Token
		expectedStmt d.Stmt
	}

	varToken := d.NewToken(d.VAR, "var", nil, 0)
	vToken := d.NewToken(d.IDENTIFIER, "v", nil, 0)
	eqToken := d.NewToken(d.EQUAL, "=", nil, 0)
	printToken := d.NewToken(d.PRINT, "print", nil, 0)

	stmtTestCases := []ParseStmtTestCase{
		{[]*d.Token{varToken, vToken, eqToken, one, semicolon}, d.VarStmt{Name: *d.NewToken(d.IDENTIFIER, "v", nil, 0), Initializer: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{printToken, one, semicolon}, d.PrintStmt{Expression: d.LiteralExpr{Value: 1}}},
	}

	for _, c := range stmtTestCases {
		t.Run(fmt.Sprintf("Parses raw stmt tokens: %s", util.SprintTokens(c.rawTokens)), func(t *testing.T) {
			assert := assert.New(t)

			stmts, err := NewParser(c.rawTokens).Parse()
			assert.NoError(err)

			assert.Len(stmts, 1)
			st := stmts[0]

			assert.True(util.IsEqualStmt(c.expectedStmt, st))
		})
	}
}
