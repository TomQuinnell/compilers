package ast

import (
	d "example/compilers/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrinter(t *testing.T) {

	t.Run("Prints expr", func(t *testing.T) {
		assert := assert.New(t)

		printer := NewAstPrinter()

		expr := d.BinaryExpr[string]{
			Left: d.UnaryExpr[string]{
				Operator: *d.NewToken(d.MINUS, "-", nil, 1),
				Right: d.LiteralExpr[string]{
					Value: d.IntStringer(123),
				},
			},
			Operator: *d.NewToken(d.STAR, "*", nil, 1),
			Right: d.GroupingExpr[string]{
				Expression: d.LiteralExpr[string]{
					Value: d.FloatStringer(45.67),
				},
			},
		}
		ret := printer.Print(expr)

		expectedRet := "(* (- 123) (group 45.670000))"
		assert.Equal(expectedRet, ret)
	})

}
