package eval

import (
	d "example/compilers/domain"
	"example/compilers/util"
	"fmt"
	"reflect"
)

type ErrInterpret struct {
	message string
	token   *d.Token
}

func (e ErrInterpret) Error() string {
	return fmt.Sprintf("%s %s", e.message, e.token)
}

func newErrInterpret(t *d.Token, msg string) ErrInterpret {
	return ErrInterpret{
		token:   t,
		message: msg,
	}
}

type Interpreter struct {
	d.ExprVisitor
	d.StmtVisitor
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(stmts []d.Stmt) error {
	for _, s := range stmts {
		err := i.execute(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) execute(s d.Stmt) error {
	_, err := s.Accept(i)
	return err
}

func (i *Interpreter) VisitExpressionStmt(s d.ExpressionStmt) (interface{}, error) {
	_, err := i.evaluate(s.Expression)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Interpreter) VisitPrintStmt(s d.PrintStmt) (interface{}, error) {
	v, err := i.evaluate(s.Expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(util.ToString(v))
	return nil, nil
}

func (i *Interpreter) VisitLiteralExpr(e d.LiteralExpr) (interface{}, error) {
	return e.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(e d.GroupingExpr) (interface{}, error) {
	return i.evaluate(e.Expression)
}

func (i *Interpreter) VisitUnaryExpr(e d.UnaryExpr) (interface{}, error) {
	right, err := i.evaluate(e.Right)
	if err != nil {
		return nil, err
	}

	switch e.Operator.Kind {
	case d.MINUS:
		rVal, err := util.ToDouble(right)
		if err != nil {
			return nil, newErrInterpret(e.Operator, "expected floaty literal")
		}
		return -float64(rVal), nil
	case d.BANG:
		return !i.isTruthy(right), nil
	}

	return nil, newErrInterpret(e.Operator, "invalid unary operator")
}

func (i *Interpreter) VisitBinaryExpr(e d.BinaryExpr) (interface{}, error) {
	left, err := i.evaluate(e.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(e.Right)
	if err != nil {
		return nil, err
	}

	if e.Operator.Kind == d.PLUS {
		switch l := left.(type) {
		case float64:
			switch r := right.(type) {
			case float64:
				return l + r, nil
			default:
				return nil, newErrInterpret(e.Operator, "expected floaty literal")
			}
		case string:
			switch r := right.(type) {
			case string:
				return l + r, nil
			default:
				return nil, newErrInterpret(e.Operator, "expected stringy literal")
			}
		default:
			return nil, newErrInterpret(e.Operator, "unsupported literal type")
		}

	}

	leftVal, err := util.ToDouble(left)
	if err != nil {
		return nil, newErrInterpret(e.Operator, "expected floaty literal")
	}
	rightVal, err := util.ToDouble(right)
	if err != nil {
		return nil, newErrInterpret(e.Operator, "expected floaty literal")
	}

	switch e.Operator.Kind {
	case d.MINUS:
		return leftVal - rightVal, nil
	case d.SLASH:
		return leftVal / rightVal, nil
	case d.STAR:
		return leftVal * rightVal, nil
	case d.GREATER:
		return leftVal > rightVal, nil
	case d.GREATER_EQUAL:
		return leftVal >= rightVal, nil
	case d.LESS:
		return leftVal < rightVal, nil
	case d.LESS_EQUAL:
		return leftVal <= rightVal, nil
	case d.BANG_EQUAL:
		return !i.isEqual(left, right), nil
	case d.EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	}

	return nil, newErrInterpret(e.Operator, "invalid binary operator")
}

func (i *Interpreter) evaluate(e d.Expr) (interface{}, error) {
	return e.Accept(i)
}

func (i *Interpreter) isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	switch b := v.(type) {
	case bool:
		return b
	}

	return true
}

func (i *Interpreter) isEqual(a interface{}, b interface{}) bool {
	if a == nil {
		return b == nil
	}

	return reflect.DeepEqual(a, b)
}
