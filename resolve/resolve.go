package resolve

import (
	d "example/compilers/domain"
	"example/compilers/eval"
	"fmt"
)

type ClassType = int32

const (
	ClassType_None = iota
	ClassType_Class
	ClassType_SubClass
)

type ErrResolve struct {
	message string
	token   *d.Token
}

func (e ErrResolve) Error() string {
	return fmt.Sprintf("%s %s", e.message, e.token)
}

func newErrResolve(t *d.Token, msg string) ErrResolve {
	return ErrResolve{
		token:   t,
		message: msg,
	}
}

var _ d.ExprVisitor = (*Resolver)(nil)
var _ d.StmtVisitor = (*Resolver)(nil)

type scope = map[string]bool

func newScope() scope {
	return make(map[string]bool)
}

type Resolver struct {
	interpreter  *eval.Interpreter
	scopes       []scope
	currentFunc  d.FunctionType
	currentClass ClassType
}

func NewResolver(interpreter *eval.Interpreter) *Resolver {
	return &Resolver{
		interpreter:  interpreter,
		scopes:       make([]scope, 0),
		currentFunc:  d.FUNCTION_TYPE_NONE,
		currentClass: ClassType_None,
	}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, newScope())
}

func (r *Resolver) peekScope() scope {
	return r.scopes[len(r.scopes)-1]
}

func (r *Resolver) getFromScope(name *d.Token) (exists bool, ok bool) {
	if len(r.scopes) == 0 {
		return false, false
	}
	s := r.peekScope()
	exists, ok = s[name.Lexeme]

	return exists, ok
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) VisitBlockStmt(s d.BlockStmt) error {
	r.beginScope()
	err := r.Resolve(s.Stmts)
	if err != nil {
		return err
	}

	r.endScope()

	return nil
}

func (r *Resolver) Resolve(stmts []d.Stmt) error {
	for _, stmt := range stmts {
		err := r.resolveStmt(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) resolveStmt(stmt d.Stmt) error {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr d.Expr) error {
	_, err := expr.Accept(r)
	return err
}

func (r *Resolver) VisitVarStmt(s d.VarStmt) error {
	err := r.declare(s.Name)
	if err != nil {
		return err
	}

	if s.Initializer != nil {
		err := r.resolveExpr(s.Initializer)
		if err != nil {
			return err
		}
	}

	r.define(s.Name)

	return nil
}

func (r *Resolver) declare(name *d.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}

	sc := r.peekScope()
	if _, ok := sc[name.Lexeme]; ok {
		return newErrResolve(name, "already a variable with this name in scope")
	}
	sc[name.Lexeme] = false

	return nil
}

func (r *Resolver) define(name *d.Token) {
	if len(r.scopes) == 0 {
		return
	}

	sc := r.peekScope()
	sc[name.Lexeme] = true
}

func (r *Resolver) VisitVariableExpr(expr d.VariableExpr) (interface{}, error) {
	if len(r.scopes) > 0 {
		if exists, ok := r.getFromScope(expr.Name); !ok || !exists {
			return nil, newErrResolve(expr.Name, "can't read local var in own initializer")
		}
	}

	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *Resolver) resolveLocal(expr d.Expr, name *d.Token) error {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
		}
	}

	return nil
}

func (r *Resolver) VisitAssignExpr(expr d.AssignExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Value)
	if err != nil {
		return nil, err
	}

	err = r.resolveLocal(expr, expr.Name)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt d.FunctionStmt) error {
	err := r.declare(stmt.Name)
	if err != nil {
		return err
	}
	r.define(stmt.Name)

	err = r.resolveFunction(stmt, d.FUNCTION_TYPE_FN)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resolver) VisitClassStmt(stmt d.ClassStmt) error {
	enclosingClass := r.currentClass
	r.currentClass = ClassType_Class

	err := r.declare(stmt.Name)
	if err != nil {
		return err
	}
	r.define(stmt.Name)

	if stmt.SuperClass != nil && stmt.Name.Lexeme == stmt.SuperClass.Name.Lexeme {
		return newErrResolve(stmt.SuperClass.Name, "A class can't inherit from itself")
	}

	if stmt.SuperClass != nil {
		r.currentClass = ClassType_SubClass
		err := r.resolveExpr(stmt.SuperClass)
		if err != nil {
			return err
		}
	}

	if stmt.SuperClass != nil {
		r.beginScope()
		r.peekScope()["super"] = true
	}

	r.beginScope()
	r.peekScope()["this"] = true

	for _, method := range stmt.Methods {
		declaration := d.FUNCTION_TYPE_METHOD
		if method.Name.Lexeme == "init" {
			declaration = d.FUNCTION_TYPE_INITIALIZER
		}

		err = r.resolveFunction(method, declaration)
		if err != nil {
			return err
		}
	}

	r.endScope()

	if stmt.SuperClass != nil {
		r.endScope()
	}

	r.currentClass = enclosingClass

	return nil
}

func (r *Resolver) VisitSuperExpr(expr d.SuperExpr) (interface{}, error) {
	if r.currentClass == ClassType_None {
		return nil, newErrResolve(expr.Keyword, "Can't use 'super' outside of a class.")
	}
	if r.currentClass != ClassType_SubClass {
		return nil, newErrResolve(expr.Keyword, "Can't use 'super' in a class with no superclass.")
	}

	err := r.resolveLocal(expr, expr.Keyword)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr d.ThisExpr) (interface{}, error) {
	if r.currentClass == ClassType_None {
		return nil, newErrResolve(expr.Keyword, "Can't use 'this' outside of a class.")
	}

	err := r.resolveLocal(expr, expr.Keyword)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) resolveFunction(f d.FunctionStmt, fnType d.FunctionType) error {
	enclosingFnType := r.currentFunc
	r.currentFunc = fnType

	r.beginScope()
	for _, param := range f.Params {
		err := r.declare(param)
		if err != nil {
			return err
		}
		r.define(param)
	}

	err := r.Resolve(f.Body)
	if err != nil {
		return err
	}

	r.endScope()
	r.currentFunc = enclosingFnType

	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt d.ExpressionStmt) error {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitIfStmt(stmt d.IfStmt) error {
	err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(stmt.ThenBranch)
	if err != nil {
		return err
	}
	if stmt.ElseBranch != nil {
		return r.resolveStmt(stmt.ElseBranch)
	}

	return nil
}

func (r *Resolver) VisitPrintStmt(stmt d.PrintStmt) error {
	return r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitReturnStmt(stmt d.ReturnStmt) error {
	if r.currentFunc == d.FUNCTION_TYPE_NONE {
		return newErrResolve(stmt.Keyword, "can't return from top-level code")
	}
	if stmt.Value != nil {
		if r.currentFunc == d.FUNCTION_TYPE_INITIALIZER {
			return newErrResolve(stmt.Keyword, "Can't return a value from an init.")
		}
		return r.resolveExpr(stmt.Value)
	}

	return nil
}

func (r *Resolver) VisitWhileStmt(stmt d.WhileStmt) error {
	err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(stmt.Body)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resolver) VisitBinaryExpr(expr d.BinaryExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}
	err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr d.CallExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Callee)
	if err != nil {
		return nil, err
	}
	for _, arg := range expr.Args {
		err = r.resolveExpr(arg)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr d.GetExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Object)
	return nil, err
}

func (r *Resolver) VisitSetExpr(expr d.SetExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Value)
	if err != nil {
		return nil, err
	}
	err = r.resolveExpr(expr.Object)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr d.GroupingExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr d.LiteralExpr) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr d.LogicalExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr d.UnaryExpr) (interface{}, error) {
	err := r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
