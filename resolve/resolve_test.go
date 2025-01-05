package resolve

import (
	d "example/compilers/domain"
	"example/compilers/eval"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolve(t *testing.T) {
	type ResolveTestCase struct {
		stmts []d.Stmt
	}

	vToken := d.NewToken(d.VAR, "v", nil, 0)
	initToken := d.NewToken(d.IDENTIFIER, "init", nil, 0)
	radiusToken := d.NewToken(d.IDENTIFIER, "radius", nil, 0)
	returnToken := d.NewToken(d.RETURN, "return", nil, 0)

	testCases := []ResolveTestCase{
		// Inherit from iteself
		{[]d.Stmt{
			d.ClassStmt{
				Name: vToken,
				SuperClass: &d.VariableExpr{
					Name: vToken,
				},
				Methods: []d.FunctionStmt{},
			},
		}},
		// Super outside of class
		{[]d.Stmt{
			d.ExpressionStmt{
				Expression: d.SuperExpr{
					Keyword: vToken,
					Method:  radiusToken,
				},
			}},
		},
		// Super in a non-subclass
		{[]d.Stmt{
			d.ClassStmt{
				Name:       vToken,
				SuperClass: nil,
				Methods: []d.FunctionStmt{{
					Name:   radiusToken,
					Params: []*d.Token{},
					Body: []d.Stmt{
						d.ExpressionStmt{
							Expression: d.SuperExpr{
								Keyword: vToken,
								Method:  radiusToken,
							},
						}},
				}},
			},
		}},
		// Own local var in initializer
		{[]d.Stmt{
			d.VarStmt{
				Name:        vToken,
				Initializer: d.VariableExpr{Name: vToken},
			},
		}},
		// Multiple declarations
		{[]d.Stmt{
			d.BlockStmt{
				Stmts: []d.Stmt{
					d.VarStmt{
						Name:        vToken,
						Initializer: d.LiteralExpr{Value: 1},
					},
					d.VarStmt{
						Name:        vToken,
						Initializer: d.LiteralExpr{Value: 1},
					},
				},
			},
		}},
		// Return from top-level code
		{[]d.Stmt{
			d.ReturnStmt{Keyword: vToken, Value: d.LiteralExpr{Value: 1}},
		}},
		// This from top-level code
		{[]d.Stmt{
			d.ExpressionStmt{
				Expression: d.ThisExpr{Keyword: vToken},
			},
		}},
		// Return from init
		{[]d.Stmt{
			d.ClassStmt{
				Name: vToken,
				Methods: []d.FunctionStmt{
					{
						Name:   initToken,
						Params: []*d.Token{radiusToken},
						Body: []d.Stmt{d.ReturnStmt{
							Keyword: returnToken,
							Value:   d.LiteralExpr{Value: 1},
						}}},
				},
			},
		},
		},
	}

	for _, c := range testCases {
		t.Run("Does not resolve:", func(t *testing.T) {
			assert := assert.New(t)

			err := NewResolver(eval.NewInterpreter()).Resolve(c.stmts)
			assert.Error(err)
		})
	}
}
