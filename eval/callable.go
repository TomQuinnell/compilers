package eval

import (
	d "example/compilers/domain"
	"example/compilers/env"
	"fmt"
	"time"
)

type Callable interface {
	Arity() int
	Call(in *Interpreter, args []interface{}) (interface{}, error)
}

type Func struct {
	declaration d.FunctionStmt
	closure     *env.Environment
}

func newFunc(declaration d.FunctionStmt, closure *env.Environment) Func {
	return Func{
		declaration: declaration,
		closure:     closure,
	}
}

var _ Callable = (*Func)(nil)

func (f Func) Arity() int {
	return len(f.declaration.Params)
}

func (f Func) Call(in *Interpreter, args []interface{}) (returnVal interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(ReturnVal); ok {
				returnVal = v.Value
				return
			}

			panic(err)
		}
	}()

	fnEnv := env.NewEnv(f.closure)

	for i := range f.Arity() {
		fnEnv.Define(f.declaration.Params[i].Lexeme, args[i])
	}

	err = in.executeBlock(f.declaration.Body, fnEnv)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (f Func) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

type ClockCallable struct{}

var _ Callable = (*ClockCallable)(nil)

func (cb ClockCallable) Arity() int {
	return 0
}

func (cb ClockCallable) Call(in *Interpreter, args []interface{}) (interface{}, error) {
	return time.Now().Unix(), nil
}

func (cb ClockCallable) String() string {
	return "<clock native fn>"
}
