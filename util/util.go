package util

import (
	"errors"
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
	case d.AssignExpr:
		switch o.(type) {
		case d.AssignExpr:
			expected, other := e.(d.AssignExpr), o.(d.AssignExpr)
			return expected.Name.Lexeme == other.Name.Lexeme &&
				IsEqualExpr(expected.Value, other.Value)
		}
		return false
	}

	return false
}

func IsEqualStmt(s, o d.Stmt) bool {
	switch s.(type) {
	case d.ExpressionStmt:
		switch o.(type) {
		case d.ExpressionStmt:
			expected, other := s.(d.ExpressionStmt), o.(d.ExpressionStmt)
			return IsEqualExpr(expected.Expression, other.Expression)
		}
		return false
	case d.PrintStmt:
		switch o.(type) {
		case d.PrintStmt:
			expected, other := s.(d.PrintStmt), o.(d.PrintStmt)
			return IsEqualExpr(expected.Expression, other.Expression)
		}
		return false
	case d.VarStmt:
		switch o.(type) {
		case d.VarStmt:
			expected, other := s.(d.VarStmt), o.(d.VarStmt)
			return expected.Name.Lexeme == other.Name.Lexeme &&
				IsEqualExpr(expected.Initializer, other.Initializer)
		}
		return false
	case d.BlockStmt:
		switch o.(type) {
		case d.BlockStmt:
			expected, other := s.(d.BlockStmt), o.(d.BlockStmt)
			if len(expected.Stmts) != len(other.Stmts) {
				return false
			}
			for i, eStmt := range expected.Stmts {
				isEqst := IsEqualStmt(eStmt, other.Stmts[i])
				if !isEqst {
					return false
				}
			}
			return true
		}
		return false
	}

	return false
}

func ToDouble(v interface{}) (float64, error) {
	switch i := v.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case int8:
		return float64(i), nil
	case int16:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	}

	return 0.0, errors.New("expected floaty value")
}
