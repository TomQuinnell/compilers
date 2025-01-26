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

	env     *env.Environment
	globals *env.Environment
	locals  map[d.Expr]int
}

func NewInterpreter() *Interpreter {
	globals := env.NewEnv(nil)
	globals.Define("clock", ClockCallable{})
	globals.Define("input", InputCallable{})

	return &Interpreter{
		env:     globals,
		globals: globals,
		locals:  make(map[d.Expr]int),
	}
}

func (i *Interpreter) Resolve(expr d.Expr, depth int) {
	i.locals[expr] = depth
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

func (i *Interpreter) VisitClassStmt(s d.ClassStmt) error {
	var superclass *Class
	if s.SuperClass != nil {
		sc, err := i.evaluate(s.SuperClass)
		if err != nil {
			return err
		}

		if clz, ok := sc.(*Class); ok {
			superclass = clz
		} else {
			return newErrInterpret(s.SuperClass.Name, "Superclass must be a class")
		}
	}

	i.env.Define(s.Name.Lexeme, nil)

	if s.SuperClass != nil {
		i.env = env.NewEnv(i.env)
		i.env.Define("super", superclass)
	}

	methods := make(map[string]Func)
	for _, method := range s.Methods {
		methods[method.Name.Lexeme] = newFunc(method, i.env, method.Name.Lexeme == "init")
	}

	klass := newClass(s.Name.Lexeme, superclass, methods)

	if s.SuperClass != nil {
		i.env = i.env.GetEnclosing()
	}

	i.env.Assign(s.Name, klass)
	return nil
}

func (i *Interpreter) VisitSuperExpr(expr d.SuperExpr) (interface{}, error) {
	distance := i.locals[expr]
	superclassRaw, err := i.env.GetAt(distance, "super")
	if err != nil {
		return nil, err
	}
	superclass := superclassRaw.(Class)

	instanceRaw, err := i.env.GetAt(distance-1, "this")
	if err != nil {
		return nil, err
	}
	instance := instanceRaw.(Instance)

	method := superclass.FindMethod(expr.Method.Lexeme)

	if method == nil {
		return nil, newErrInterpret(expr.Method, fmt.Sprintf("Undefined property '%s'.", expr.Method.Lexeme))
	}

	return method.Bind(instance), nil
}

func (i *Interpreter) VisitFunctionStmt(s d.FunctionStmt) error {
	fn := newFunc(s, i.env, false)
	i.env.Define(s.Name.Lexeme, fn)

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

type ReturnVal struct {
	Value interface{}
}

func (i *Interpreter) VisitReturnStmt(s d.ReturnStmt) error {
	var value interface{}
	if s.Value != nil {
		var err error
		value, err = i.evaluate(s.Value)
		if err != nil {
			return err
		}
	}

	// I don't like this but it's the easiest way to unwrap the
	// call stack to get back to the caller...
	panic(ReturnVal{Value: value})
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

func (i *Interpreter) VisitCallExpr(e d.CallExpr) (interface{}, error) {
	callee, err := i.evaluate(e.Callee)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, len(e.Args))
	for j, arg := range e.Args {
		argV, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}

		args[j] = argV
	}

	cb, ok := callee.(Callable)
	if !ok {
		return nil, newErrInterpret(e.Paren, "can only call function/class")
	}

	if len(args) != cb.Arity() {
		return nil, newErrInterpret(
			e.Paren,
			fmt.Sprintf("expected %d args but got %d instead.", len(args), cb.Arity()))
	}

	return cb.Call(i, args)
}

func (i *Interpreter) VisitGetExpr(e d.GetExpr) (interface{}, error) {
	obj, err := i.evaluate(e.Object)
	if err != nil {
		return nil, err
	}

	if instance, ok := obj.(Instance); ok {
		return instance.Get(e.Name)
	}

	return nil, newErrInterpret(e.Name, "Only instances have properties")
}

func (i *Interpreter) VisitSetExpr(e d.SetExpr) (interface{}, error) {
	obj, err := i.evaluate(e.Object)
	if err != nil {
		return nil, err
	}

	if instance, ok := obj.(Instance); ok {
		value, err := i.evaluate(e.Value)
		if err != nil {
			return nil, err
		}

		instance.Set(e.Name, value)
		return value, nil
	}

	return nil, newErrInterpret(e.Name, "Only instances have fields")
}

func (i *Interpreter) VisitThisExpr(e d.ThisExpr) (interface{}, error) {
	return i.lookUpVariable(e.Keyword, e)
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
	return i.lookUpVariable(e.Name, e)
}

func (i *Interpreter) lookUpVariable(name *d.Token, expr d.Expr) (interface{}, error) {
	distance, ok := i.locals[expr]
	if ok {
		return i.env.GetAt(distance, name.Lexeme)
	} else {
		return i.globals.Get(name)
	}
}

func (i *Interpreter) VisitAssignExpr(e d.AssignExpr) (interface{}, error) {
	v, err := i.evaluate(e.Value)
	if err != nil {
		return nil, err
	}

	distance, ok := i.locals[e]
	if ok {
		i.env.AssignAt(distance, e.Name, v)
	} else {
		i.globals.Assign(e.Name, v)
	}

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
