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
	semicolon := d.NewToken(d.SEMICOLON, ";", nil, 0)

	fnToken := d.NewToken(d.FUN, "fn", nil, 0)
	returnToken := d.NewToken(d.RETURN, "return", nil, 0)
	commaToken := d.NewToken(d.COMMA, ",", nil, 0)
	varToken := d.NewToken(d.VAR, "var", nil, 0)
	vToken := d.NewToken(d.IDENTIFIER, "v", nil, 0)
	v1Token := d.NewToken(d.IDENTIFIER, "v1", nil, 0)
	v2Token := d.NewToken(d.IDENTIFIER, "v2", nil, 0)
	eqToken := d.NewToken(d.EQUAL, "=", nil, 0)
	printToken := d.NewToken(d.PRINT, "print", nil, 0)
	ifToken := d.NewToken(d.IF, "if", nil, 0)
	elseToken := d.NewToken(d.ELSE, "else", nil, 0)
	forToken := d.NewToken(d.FOR, "for", nil, 0)
	whileToken := d.NewToken(d.WHILE, "while", nil, 0)
	openBlockToken := d.NewToken(d.LEFT_BRACE, "{", nil, 0)
	closeBlockToken := d.NewToken(d.RIGHT_BRACE, "}", nil, 0)
	orToken := d.NewToken(d.OR, "or", nil, 0)
	andToken := d.NewToken(d.AND, "and", nil, 0)

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

			rawTokens := append(c.rawTokens, semicolon)
			stmts, err := NewParser(rawTokens).Parse()
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

	t.Run("Parses 2 raw exprs", func(t *testing.T) {
		assert := assert.New(t)

		rawTokens := []*d.Token{one, eqeq, one, semicolon, one, eqeq, one, semicolon}

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

	stmtTestCases := []ParseStmtTestCase{
		{[]*d.Token{varToken, vToken, eqToken, one, semicolon}, d.VarStmt{Name: d.NewToken(d.IDENTIFIER, "v", nil, 0), Initializer: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{printToken, one, semicolon}, d.PrintStmt{Expression: d.LiteralExpr{Value: 1}}},
		{[]*d.Token{ifToken, openBracket, one, eqeq, one, orToken, one, eqeq, a, closeBracket, openBlockToken, printToken, one, semicolon, closeBlockToken},
			d.IfStmt{
				Condition: d.LogicalExpr{
					Left:     d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: eqeq, Right: d.LiteralExpr{Value: 1}},
					Operator: orToken,
					Right:    d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: eqeq, Right: d.LiteralExpr{Value: "a"}},
				},
				ThenBranch: d.BlockStmt{Stmts: []d.Stmt{d.PrintStmt{Expression: d.LiteralExpr{Value: 1}}}},
				ElseBranch: nil,
			},
		},
		{[]*d.Token{ifToken, openBracket, one, eqeq, one, andToken, one, eqeq, a, closeBracket, openBlockToken, printToken, one, semicolon, closeBlockToken, elseToken, openBlockToken, printToken, a, semicolon, closeBlockToken},
			d.IfStmt{
				Condition: d.LogicalExpr{
					Left:     d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: eqeq, Right: d.LiteralExpr{Value: 1}},
					Operator: andToken,
					Right:    d.BinaryExpr{Left: d.LiteralExpr{Value: 1}, Operator: eqeq, Right: d.LiteralExpr{Value: "a"}},
				},
				ThenBranch: d.BlockStmt{Stmts: []d.Stmt{d.PrintStmt{Expression: d.LiteralExpr{Value: 1}}}},
				ElseBranch: d.BlockStmt{Stmts: []d.Stmt{d.PrintStmt{Expression: d.LiteralExpr{Value: "a"}}}},
			},
		},
		{[]*d.Token{whileToken, openBracket, trutho, closeBracket, openBlockToken, printToken, one, semicolon, closeBlockToken},
			d.WhileStmt{
				Condition: d.LiteralExpr{Value: true},
				Body:      d.BlockStmt{Stmts: []d.Stmt{d.PrintStmt{Expression: d.LiteralExpr{Value: 1}}}},
			},
		},
		{[]*d.Token{forToken, openBracket, semicolon, semicolon, closeBracket, openBlockToken, printToken, one, semicolon, closeBlockToken},
			d.WhileStmt{
				Condition: d.LiteralExpr{Value: true},
				Body:      d.BlockStmt{Stmts: []d.Stmt{d.PrintStmt{Expression: d.LiteralExpr{Value: 1}}}},
			},
		},
		// Increment
		{[]*d.Token{forToken, openBracket, semicolon, semicolon, one, closeBracket, openBlockToken, printToken, one, semicolon, closeBlockToken},
			d.WhileStmt{
				Condition: d.LiteralExpr{Value: true},
				Body: d.BlockStmt{
					Stmts: []d.Stmt{d.BlockStmt{Stmts: []d.Stmt{d.PrintStmt{Expression: d.LiteralExpr{Value: 1}}}}, d.ExpressionStmt{Expression: d.LiteralExpr{Value: 1}}},
				},
			},
		},
		// Initializer
		{[]*d.Token{forToken, openBracket, one, semicolon, semicolon, closeBracket, openBlockToken, printToken, one, semicolon, closeBlockToken},
			d.BlockStmt{
				Stmts: []d.Stmt{
					d.ExpressionStmt{Expression: d.LiteralExpr{Value: 1}},
					d.WhileStmt{
						Condition: d.LiteralExpr{Value: true},
						Body:      d.BlockStmt{Stmts: []d.Stmt{d.PrintStmt{Expression: d.LiteralExpr{Value: 1}}}},
					},
				},
			},
		},
		// Full for loop
		{[]*d.Token{forToken, openBracket, one, semicolon, a, semicolon, one, closeBracket, openBlockToken, printToken, one, semicolon, closeBlockToken},
			d.BlockStmt{
				Stmts: []d.Stmt{
					d.ExpressionStmt{Expression: d.LiteralExpr{Value: 1}},
					d.WhileStmt{
						Condition: d.LiteralExpr{Value: "a"},
						Body: d.BlockStmt{Stmts: []d.Stmt{
							d.BlockStmt{Stmts: []d.Stmt{
								d.PrintStmt{Expression: d.LiteralExpr{Value: 1}},
							}},
							d.ExpressionStmt{Expression: d.LiteralExpr{Value: 1}},
						}},
					},
				},
			},
		},
		// Empty function w/ return
		{[]*d.Token{fnToken, vToken, openBracket, closeBracket, openBlockToken, returnToken, nilo, semicolon, closeBlockToken}, d.FunctionStmt{
			Name:   vToken,
			Params: []*d.Token{},
			Body: []d.Stmt{d.ReturnStmt{
				Keyword: returnToken,
				Value:   d.LiteralExpr{Value: nil},
			}},
		}},
		// N-arity function w/ return
		{[]*d.Token{fnToken, vToken, openBracket, v1Token, commaToken, v2Token, closeBracket, openBlockToken, returnToken, nilo, semicolon, closeBlockToken}, d.FunctionStmt{
			Name:   vToken,
			Params: []*d.Token{v1Token, v2Token},
			Body: []d.Stmt{d.ReturnStmt{
				Keyword: returnToken,
				Value:   d.LiteralExpr{Value: nil},
			}},
		}},
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

	t.Run("Parses assignment block", func(t *testing.T) {
		assert := assert.New(t)

		rawTokens := []*d.Token{
			openBlockToken,
			varToken, vToken, eqToken, one, semicolon,
			vToken, eqToken, a, semicolon,
			closeBlockToken,
		}

		stmts, err := NewParser(rawTokens).Parse()
		assert.NoError(err)

		assert.Len(stmts, 1)
		st := stmts[0]

		expectedStmts := []d.Stmt{
			d.VarStmt{Name: vToken, Initializer: d.LiteralExpr{Value: 1}},
			d.ExpressionStmt{Expression: d.AssignExpr{Name: vToken, Value: d.LiteralExpr{Value: "a"}}},
		}
		assert.True(util.IsEqualStmt(d.BlockStmt{Stmts: expectedStmts}, st))
	})

	classToken := d.NewToken(d.CLASS, "class", nil, 0)
	initToken := d.NewToken(d.IDENTIFIER, "init", nil, 0)
	radiusToken := d.NewToken(d.IDENTIFIER, "radius", nil, 0)
	getRadiusToken := d.NewToken(d.IDENTIFIER, "getRadius", nil, 0)
	thisToken := d.NewToken(d.THIS, "this", nil, 0)
	dotToken := d.NewToken(d.DOT, ".", nil, 0)

	t.Run("Parses class block", func(t *testing.T) {
		assert := assert.New(t)

		rawTokens := []*d.Token{
			classToken, vToken, openBlockToken,
			initToken, openBracket, radiusToken, closeBracket, openBlockToken,
			thisToken, dotToken, radiusToken, eqToken, radiusToken, semicolon,
			closeBlockToken,

			getRadiusToken, openBracket, closeBracket, openBlockToken,
			returnToken, thisToken, dotToken, radiusToken, semicolon,
			closeBlockToken,

			closeBlockToken,

			vToken, openBracket, one, closeBracket, dotToken, getRadiusToken, openBracket, closeBracket, semicolon,
		}

		stmts, err := NewParser(rawTokens).Parse()
		assert.NoError(err)

		assert.Len(stmts, 2)
		st1 := stmts[0]

		expectedStmt1 := d.ClassStmt{
			Name: vToken,
			Methods: []d.FunctionStmt{
				{
					Name:   initToken,
					Params: []*d.Token{radiusToken},
					Body: []d.Stmt{d.ExpressionStmt{Expression: d.SetExpr{
						Object: d.ThisExpr{Keyword: thisToken},
						Name:   radiusToken,
						Value:  d.VariableExpr{Name: radiusToken},
					}}},
				},
				{
					Name:   getRadiusToken,
					Params: []*d.Token{},
					Body: []d.Stmt{d.ReturnStmt{
						Keyword: returnToken,
						Value: d.GetExpr{
							Object: d.ThisExpr{Keyword: thisToken},
							Name:   radiusToken,
						},
					}},
				},
			},
		}
		assert.True(util.IsEqualStmt(expectedStmt1, st1))

		st2 := stmts[1]

		expectedStmt2 := d.ExpressionStmt{Expression: d.CallExpr{
			Callee: d.GetExpr{
				Object: d.CallExpr{
					Callee: d.VariableExpr{Name: vToken},
					Paren:  closeBracket,
					Args:   []d.Expr{d.LiteralExpr{Value: 1}},
				},
				Name: getRadiusToken,
			},
			Paren: closeBracket,
			Args:  []d.Expr{},
		}}
		assert.True(util.IsEqualStmt(expectedStmt2, st2))
	})

	t.Run("Parses function call", func(t *testing.T) {
		assert := assert.New(t)

		rawTokens := []*d.Token{
			fnToken, vToken, openBracket, v1Token, commaToken, v2Token, closeBracket,
			openBlockToken, returnToken, nilo, semicolon, closeBlockToken,
			vToken, openBracket, one, commaToken, one, closeBracket, semicolon,
		}

		stmts, err := NewParser(rawTokens).Parse()
		assert.NoError(err)

		assert.Len(stmts, 2)
		st0 := stmts[0]
		expectedStmt0 := d.FunctionStmt{
			Name:   vToken,
			Params: []*d.Token{v1Token, v2Token},
			Body:   []d.Stmt{d.ReturnStmt{Keyword: returnToken, Value: d.LiteralExpr{Value: nil}}},
		}
		assert.True(util.IsEqualStmt(expectedStmt0, st0))
		st1 := stmts[1]
		expectedStmt1 := d.ExpressionStmt{
			Expression: d.CallExpr{
				Callee: d.VariableExpr{Name: vToken},
				Paren:  closeBracket,
				Args:   []d.Expr{d.LiteralExpr{Value: 1}, d.LiteralExpr{Value: 1}},
			},
		}
		assert.True(util.IsEqualStmt(expectedStmt1, st1))
	})

	openBracketToken := d.NewToken(d.LEFT_PAREN, "(", nil, 0)

	errTestCases := []ParseStmtTestCase{
		{[]*d.Token{varToken}, nil},
		{[]*d.Token{varToken, vToken, eqToken, one}, nil},
		{[]*d.Token{printToken, one}, nil},
		{[]*d.Token{openBlockToken, one}, nil},
		{[]*d.Token{openBracketToken, one}, nil},
		{[]*d.Token{forToken, one}, nil},
		{[]*d.Token{forToken, openBracket, one, closeBracket}, nil},
		{[]*d.Token{forToken, openBracket, one, semicolon}, nil},
		{[]*d.Token{ifToken, one}, nil},
		{[]*d.Token{ifToken, openBracket, one}, nil},
		{[]*d.Token{whileToken, one}, nil},
		{[]*d.Token{whileToken, openBracket, one}, nil},
		{[]*d.Token{fnToken}, nil},
		{[]*d.Token{fnToken, vToken}, nil},
		{[]*d.Token{fnToken, vToken, openBracketToken}, nil},
		{[]*d.Token{fnToken, vToken, openBracketToken, v1Token, commaToken}, nil},
		{[]*d.Token{fnToken, vToken, openBracketToken}, nil},
		{[]*d.Token{fnToken, vToken, openBracketToken, closeBracket, returnToken}, nil},
		{[]*d.Token{classToken}, nil},
		{[]*d.Token{classToken, vToken}, nil},
		{[]*d.Token{classToken, vToken, openBlockToken}, nil},
		{[]*d.Token{classToken, vToken, openBlockToken, closeBlockToken, vToken, dotToken}, nil},
	}

	for _, c := range errTestCases {
		t.Run(fmt.Sprintf("Errors raw stmt tokens: %s", util.SprintTokens(c.rawTokens)), func(t *testing.T) {
			assert := assert.New(t)

			_, err := NewParser(c.rawTokens).Parse()
			assert.Error(err)
		})
	}

	maxArgsSize := 255
	bigFnTokens := []*d.Token{fnToken, vToken, openBracketToken}
	for range maxArgsSize {
		bigFnTokens = append(bigFnTokens, v1Token, commaToken)
	}
	bigFnTokens = append(bigFnTokens, v2Token, closeBracket, openBlockToken, closeBlockToken)
	t.Run("Errors big function args", func(t *testing.T) {
		assert := assert.New(t)

		_, err := NewParser(bigFnTokens).Parse()
		assert.Error(err)
	})
}
