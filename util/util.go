package util

import (
	d "example/compilers/domain"
	"fmt"
)

func RuneAt(s string, i int) rune {
	runeSlice := []rune(s)
	return runeSlice[i]
}

func IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func IsAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func IsAlphaNumeric(c rune) bool {
	return IsAlpha(c) || IsDigit(c)
}

func SprintTokens(ts []*d.Token) string {
	ret := ""
	for _, t := range ts {
		ret += t.String()
	}
	return ret
}

func ToString(o interface{}) string {
	if stringer, ok := o.(interface{ String() }); ok {
		stringer.String()
	}
	return fmt.Sprintf("%v", o)
}

func IsEqualExpr(e, o d.Expr) bool {
	switch e.(type) {
	case d.BinaryExpr:
		switch o.(type) {
		case d.BinaryExpr:
			expected, other := e.(d.BinaryExpr), o.(d.BinaryExpr)
			return expected.Operator.Kind == other.Operator.Kind &&
				IsEqualExpr(expected.Left, other.Left) &&
				IsEqualExpr(expected.Right, other.Right)
		}
		return false
	case d.UnaryExpr:
		switch o.(type) {
		case d.UnaryExpr:
			expected, other := e.(d.UnaryExpr), o.(d.UnaryExpr)
			return expected.Operator.Kind == other.Operator.Kind &&
				IsEqualExpr(expected.Right, other.Right)
		}
		return false
	case d.LiteralExpr:
		switch o.(type) {
		case d.LiteralExpr:
			expected, other := e.(d.LiteralExpr), o.(d.LiteralExpr)
			return expected.Value == other.Value
		}
		return false
	case d.GroupingExpr:
		switch o.(type) {
		case d.GroupingExpr:
			expected, other := e.(d.GroupingExpr), o.(d.GroupingExpr)
			return IsEqualExpr(expected.Expression, other.Expression)
		}
		return false
	}

	return false
}
