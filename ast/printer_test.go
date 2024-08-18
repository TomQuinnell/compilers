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

		expr := d.BinaryExpr{
			Left: d.UnaryExpr{
				Operator: d.NewToken(d.MINUS, "-", nil, 1),
				Right: d.LiteralExpr{
					Value: 123,
				},
			},
			Operator: d.NewToken(d.STAR, "*", nil, 1),
			Right: d.GroupingExpr{
				Expression: d.LiteralExpr{
					Value: 45.67,
				},
			},
		}
		ret := printer.Print(expr)

		expectedRet := "(* (- 123) (group 45.67))"
		assert.Equal(expectedRet, ret)
	})

}
