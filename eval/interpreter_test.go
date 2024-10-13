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

	vToken := d.NewToken(d.IDENTIFIER, "v", nil, 0)
	v1Token := d.NewToken(d.IDENTIFIER, "v1", nil, 0)
	v2Token := d.NewToken(d.IDENTIFIER, "v2", nil, 0)
	orToken := d.NewToken(d.OR, "or", nil, 0)
	andToken := d.NewToken(d.AND, "and", nil, 0)
	eqeqToken := d.NewToken(d.EQUAL_EQUAL, "==", nil, 0)
	returnToken := d.NewToken(d.RETURN, "return", nil, 0)
	closeBracketToken := d.NewToken(d.RIGHT_PAREN, ")", nil, 0)

	testStmtCases := []InterpretStmtTestCase{
		{d.PrintStmt{Expression: one}},
		{d.VarStmt{Name: d.NewToken(d.IDENTIFIER, "v", nil, 0), Initializer: one}},
		{d.BlockStmt{Stmts: []d.Stmt{
			d.BlockStmt{Stmts: []d.Stmt{
				d.VarStmt{Name: vToken, Initializer: d.LiteralExpr{Value: 1}},
				d.ExpressionStmt{Expression: d.AssignExpr{Name: vToken, Value: d.LiteralExpr{Value: "a"}}},
				d.PrintStmt{Expression: d.VariableExpr{Name: vToken}},
			}},
			d.VarStmt{Name: vToken, Initializer: d.LiteralExpr{Value: 1}},
			d.ExpressionStmt{Expression: d.AssignExpr{Name: vToken, Value: d.LiteralExpr{Value: "a"}}},
			d.PrintStmt{Expression: d.VariableExpr{Name: vToken}},
		}}},
		{d.IfStmt{
			Condition: d.LogicalExpr{
				Left:     d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: eqeqToken, Right: d.LiteralExpr{Value: 1}},
				Operator: orToken,
				Right:    d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: eqeqToken, Right: d.LiteralExpr{Value: "a"}},
			},
			ThenBranch: d.BlockStmt{Stmts: []d.Stmt{d.PrintStmt{Expression: d.LiteralExpr{Value: 1}}}},
			ElseBranch: nil,
		}},
		{d.IfStmt{
			Condition: d.LogicalExpr{
				Left:     d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: eqeqToken, Right: d.LiteralExpr{Value: 1}},
				Operator: andToken,
				Right:    d.BinaryExpr{Left: d.LiteralExpr{Value: 2}, Operator: eqeqToken, Right: d.LiteralExpr{Value: 1}},
			},
			ThenBranch: d.BlockStmt{Stmts: []d.Stmt{d.PrintStmt{Expression: d.LiteralExpr{Value: 1}}}},
			ElseBranch: d.BlockStmt{Stmts: []d.Stmt{d.PrintStmt{Expression: d.LiteralExpr{Value: "a"}}}},
		}},
		{d.WhileStmt{
			Condition: d.LiteralExpr{Value: false},
			Body:      d.BlockStmt{Stmts: []d.Stmt{d.PrintStmt{Expression: d.LiteralExpr{Value: 1}}}},
		}},
		{d.BlockStmt{Stmts: []d.Stmt{
			d.VarStmt{Name: vToken, Initializer: d.LiteralExpr{Value: 0}},
			d.WhileStmt{
				Condition: d.BinaryExpr{Left: d.VariableExpr{Name: vToken}, Operator: eqeqToken, Right: d.LiteralExpr{Value: 0}},
				Body:      d.BlockStmt{Stmts: []d.Stmt{d.ExpressionStmt{Expression: d.AssignExpr{Name: vToken, Value: d.LiteralExpr{Value: 1}}}}},
			},
		}}},
		{d.BlockStmt{Stmts: []d.Stmt{
			d.FunctionStmt{
				Name:   vToken,
				Params: []*d.Token{v1Token, v2Token},
				Body:   []d.Stmt{d.ReturnStmt{Keyword: returnToken, Value: d.LiteralExpr{Value: nil}}},
			},
			d.ExpressionStmt{
				Expression: d.CallExpr{
					Callee: d.VariableExpr{Name: vToken},
					Paren:  closeBracketToken,
					Args:   []d.Expr{d.LiteralExpr{Value: 1}, d.LiteralExpr{Value: 1}},
				},
			},
		}}},
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
	invalidArity := d.CallExpr{
		Callee: d.VariableExpr{Name: vToken},
		Paren:  closeBracketToken,
		Args:   []d.Expr{d.LiteralExpr{Value: 1}, d.LiteralExpr{Value: 1}},
	}
	invalidCallable := d.CallExpr{
		Callee: d.LiteralExpr{Value: 1},
		Paren:  closeBracketToken,
		Args:   []d.Expr{d.LiteralExpr{Value: 1}, d.LiteralExpr{Value: 1}},
	}

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
		{invalidArity, nil},
		{invalidCallable, nil},
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
		{d.ExpressionStmt{Expression: minStr}},
		{d.PrintStmt{Expression: minStr}},
		{d.VarStmt{Name: d.NewToken(d.IDENTIFIER, "v", nil, 0), Initializer: minStr}},
		{d.BlockStmt{Stmts: []d.Stmt{d.ExpressionStmt{Expression: minStr}}}},
		{d.PrintStmt{Expression: d.VariableExpr{Name: vToken}}},
		{d.IfStmt{Condition: minStr}},
		{d.WhileStmt{Condition: minStr}},
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
