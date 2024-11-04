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

	testCases := []ResolveTestCase{
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
	}

	for _, c := range testCases {
		t.Run("Does not resolve:", func(t *testing.T) {
			assert := assert.New(t)

			err := NewResolver(eval.NewInterpreter()).Resolve(c.stmts)
			assert.Error(err)
		})
	}
}
