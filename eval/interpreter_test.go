package eval

import (
	"example/compilers/ast"
	d "example/compilers/domain"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpret(t *testing.T) {
	type InterpretTestCase struct {
		expr        d.Expr
		expectedVal interface{}
	}

	one := d.LiteralExpr{Value: float64(1)}
	two := d.LiteralExpr{Value: float64(2)}
	ten := d.LiteralExpr{Value: float64(10)}
	a := d.LiteralExpr{Value: "a"}
	b := d.LiteralExpr{Value: "b"}
	falso := d.LiteralExpr{Value: false}
	trutho := d.LiteralExpr{Value: true}
	nilo := d.LiteralExpr{Value: nil}
	minOne := d.UnaryExpr{Operator: d.NewToken(d.MINUS, "min", nil, 0), Right: one}
	bangTrue := d.UnaryExpr{Operator: d.NewToken(d.BANG, "bang", nil, 0), Right: trutho}
	onePlusOne := d.BinaryExpr{
		Operator: d.NewToken(d.PLUS, "plus", nil, 0),
		Left:     one,
		Right:    one,
	}
	oneMinOne := d.BinaryExpr{
		Operator: d.NewToken(d.MINUS, "minus", nil, 0),
		Left:     one,
		Right:    one,
	}
	tenSlashTwo := d.BinaryExpr{
		Operator: d.NewToken(d.SLASH, "slash", nil, 0),
		Left:     ten,
		Right:    two,
	}
	tenGTTwo := d.BinaryExpr{
		Operator: d.NewToken(d.GREATER, "gt", nil, 0),
		Left:     ten,
		Right:    two,
	}
	tenGTETwo := d.BinaryExpr{
		Operator: d.NewToken(d.GREATER_EQUAL, "gte", nil, 0),
		Left:     ten,
		Right:    two,
	}
	tenLTTwo := d.BinaryExpr{
		Operator: d.NewToken(d.LESS, "lt", nil, 0),
		Left:     ten,
		Right:    two,
	}
	tenLTETwo := d.BinaryExpr{
		Operator: d.NewToken(d.LESS_EQUAL, "lte", nil, 0),
		Left:     ten,
		Right:    two,
	}
	tenBEqTwo := d.BinaryExpr{
		Operator: d.NewToken(d.BANG_EQUAL, "beq", nil, 0),
		Left:     ten,
		Right:    two,
	}
	tenEqEqTwo := d.BinaryExpr{
		Operator: d.NewToken(d.EQUAL_EQUAL, "eqeq", nil, 0),
		Left:     ten,
		Right:    two,
	}
	oneGTOne := d.BinaryExpr{
		Operator: d.NewToken(d.GREATER, "gt", nil, 0),
		Left:     one,
		Right:    one,
	}
	oneGTEOne := d.BinaryExpr{
		Operator: d.NewToken(d.GREATER_EQUAL, "gte", nil, 0),
		Left:     one,
		Right:    one,
	}
	oneLTOne := d.BinaryExpr{
		Operator: d.NewToken(d.LESS, "lt", nil, 0),
		Left:     one,
		Right:    one,
	}
	oneLTEOne := d.BinaryExpr{
		Operator: d.NewToken(d.LESS_EQUAL, "lte", nil, 0),
		Left:     one,
		Right:    one,
	}
	oneBeqOne := d.BinaryExpr{
		Operator: d.NewToken(d.BANG_EQUAL, "beq", nil, 0),
		Left:     one,
		Right:    one,
	}
	oneEqEqOne := d.BinaryExpr{
		Operator: d.NewToken(d.EQUAL_EQUAL, "eqeq", nil, 0),
		Left:     one,
		Right:    one,
	}
	aPlusB := d.BinaryExpr{
		Operator: d.NewToken(d.PLUS, "plus", nil, 0),
		Left:     a,
		Right:    b,
	}

	testCases := []InterpretTestCase{
		{one, float64(1)},
		{a, "a"},
		{falso, false},
		{trutho, true},
		{nilo, nil},
		{d.GroupingExpr{Expression: one}, float64(1)},
		{d.GroupingExpr{Expression: a}, "a"},
		{d.GroupingExpr{Expression: falso}, false},
		{d.GroupingExpr{Expression: trutho}, true},
		{d.GroupingExpr{Expression: nilo}, nil},
		{minOne, float64(-1)},
		{bangTrue, false},
		{onePlusOne, float64(2)},
		{aPlusB, "ab"},
		{oneMinOne, float64(0)},
		{tenSlashTwo, float64(5)},
		{tenGTTwo, true},
		{tenGTETwo, true},
		{tenLTTwo, false},
		{tenLTETwo, false},
		{tenBEqTwo, true},
		{tenEqEqTwo, false},
		{oneGTOne, false},
		{oneGTEOne, true},
		{oneLTOne, false},
		{oneLTEOne, true},
		{oneBeqOne, false},
		{oneEqEqOne, true},
	}

	for _, c := range testCases {
		s := ast.NewAstPrinter().Print(c.expr)
		t.Run(fmt.Sprintf("Interprets expr: %s", s), func(t *testing.T) {
			assert := assert.New(t)

			stmts := []d.Stmt{d.ExpressionStmt{Expression: c.expr}}
			err := NewInterpreter().Interpret(stmts)
			assert.NoError(err)
		})
	}

	type InterpretStmtTestCase struct {
		stmt d.Stmt
	}

	testStmtCases := []InterpretStmtTestCase{
		{d.PrintStmt{Expression: one}},
		{d.VarStmt{Name: *d.NewToken(d.IDENTIFIER, "v", nil, 0), Initializer: one}},
	}

	for _, c := range testStmtCases {
		t.Run("Interprets stmt: %s", func(t *testing.T) {
			assert := assert.New(t)

			stmts := []d.Stmt{c.stmt}
			err := NewInterpreter().Interpret(stmts)
			assert.NoError(err)
		})
	}

	minStr := d.UnaryExpr{Operator: d.NewToken(d.MINUS, "min", nil, 0), Right: a}
	minStrPlusOne := d.BinaryExpr{
		Operator: d.NewToken(d.PLUS, "plus", nil, 0),
		Left:     minStr,
		Right:    one,
	}
	onePlusMinStr := d.BinaryExpr{
		Operator: d.NewToken(d.PLUS, "plus", nil, 0),
		Left:     one,
		Right:    minStr,
	}
	onePlusStr := d.BinaryExpr{
		Operator: d.NewToken(d.PLUS, "plus", nil, 0),
		Left:     one,
		Right:    a,
	}
	strPlusOne := d.BinaryExpr{
		Operator: d.NewToken(d.PLUS, "plus", nil, 0),
		Left:     a,
		Right:    one,
	}
	strPlusBool := d.BinaryExpr{
		Operator: d.NewToken(d.PLUS, "plus", nil, 0),
		Left:     a,
		Right:    trutho,
	}
	strMinOne := d.BinaryExpr{
		Operator: d.NewToken(d.MINUS, "min", nil, 0),
		Left:     a,
		Right:    one,
	}
	oneMinStr := d.BinaryExpr{
		Operator: d.NewToken(d.MINUS, "min", nil, 0),
		Left:     one,
		Right:    a,
	}
	invalidBinary := d.BinaryExpr{
		Operator: d.NewToken(d.BANG, "bang", nil, 0),
		Left:     one,
		Right:    a,
	}
	invalidUnary := d.UnaryExpr{Operator: d.NewToken(d.PLUS, "plus", nil, 0), Right: a}

	errTestCases := []InterpretTestCase{
		{minStr, nil},
		{minStrPlusOne, nil},
		{onePlusMinStr, nil},
		{onePlusStr, nil},
		{strPlusOne, nil},
		{strPlusBool, nil},
		{strMinOne, nil},
		{oneMinStr, nil},
		{invalidBinary, nil},
		{invalidUnary, nil},
	}

	for _, c := range errTestCases {
		s := ast.NewAstPrinter().Print(c.expr)
		t.Run(fmt.Sprintf("Does not interpret expr: %s", s), func(t *testing.T) {
			assert := assert.New(t)

			stmts := []d.Stmt{d.ExpressionStmt{Expression: c.expr}}
			err := NewInterpreter().Interpret(stmts)
			assert.Error(err)
		})
	}

	errTestStmtCases := []InterpretStmtTestCase{
		{d.PrintStmt{Expression: minStr}},
		{d.VarStmt{Name: *d.NewToken(d.IDENTIFIER, "v", nil, 0), Initializer: minStr}},
	}

	for _, c := range errTestStmtCases {
		t.Run("Interprets stmt: %s", func(t *testing.T) {
			assert := assert.New(t)

			stmts := []d.Stmt{c.stmt}
			err := NewInterpreter().Interpret(stmts)
			assert.Error(err)
		})
	}
}
