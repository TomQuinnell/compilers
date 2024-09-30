package eval

import (
	d "example/compilers/domain"
	"example/compilers/env"
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
	env *env.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: env.NewEnv(nil),
	}
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
	return s.Accept(i)
}

func (i *Interpreter) VisitExpressionStmt(s d.ExpressionStmt) error {
	_, err := i.evaluate(s.Expression)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) VisitIfStmt(s d.IfStmt) error {
	condition, err := i.evaluate(s.Condition)
	if err != nil {
		return err
	}
	if i.isTruthy(condition) {
		return i.execute(s.ThenBranch)
	}

	if s.ElseBranch != nil {
		return i.execute(s.ElseBranch)
	}

	return nil
}

func (i *Interpreter) VisitWhileStmt(s d.WhileStmt) error {
	for {
		v, err := i.evaluate(s.Condition)
		if err != nil {
			return err
		}

		if !i.isTruthy(v) {
			return nil
		}

		err = i.execute(s.Body)
		if err != nil {
			return err
		}
	}
}

func (i *Interpreter) VisitPrintStmt(s d.PrintStmt) error {
	v, err := i.evaluate(s.Expression)
	if err != nil {
		return err
	}

	fmt.Println(util.ToString(v))
	return nil
}

func (i *Interpreter) VisitVarStmt(s d.VarStmt) error {
	var v interface{}
	if s.Initializer != nil {
		var err error
		v, err = i.evaluate(s.Initializer)
		if err != nil {
			return err
		}
	}

	i.env.Define(s.Name.Lexeme, v)
	return nil
}

func (i *Interpreter) VisitBlockStmt(s d.BlockStmt) error {
	return i.executeBlock(s.Stmts, env.NewEnv(i.env))
}

func (i *Interpreter) executeBlock(stmts []d.Stmt, environment *env.Environment) error {
	previousEnv := i.env
	defer func() {
		// Restore previuos env
		i.env = previousEnv
	}()

	i.env = environment
	for _, s := range stmts {
		err := i.execute(s)
		if err != nil {
			return err
		}
	}

	return nil
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

func (i *Interpreter) VisitLogicalExpr(e d.LogicalExpr) (interface{}, error) {
	left, err := i.evaluate(e.Left)
	if err != nil {
		return nil, err
	}

	if e.Operator.Kind == d.OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(e.Right)
}

func (i *Interpreter) VisitVariableExpr(e d.VariableExpr) (interface{}, error) {
	return i.env.Get(e.Name)
}

func (i *Interpreter) VisitAssignExpr(e d.AssignExpr) (interface{}, error) {
	v, err := i.evaluate(e.Value)
	if err != nil {
		return nil, err
	}

	i.env.Assign(e.Name, v)

	return v, nil
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
