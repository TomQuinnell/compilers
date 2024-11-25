package eval

import (
	d "example/compilers/domain"
	"fmt"
)

type ErrClass struct {
	message string
	token   *d.Token
}

func (e ErrClass) Error() string {
	return fmt.Sprintf("%s %s", e.message, e.token)
}

func newErrClass(t *d.Token, msg string) ErrClass {
	return ErrClass{
		token:   t,
		message: msg,
	}
}

var _ Callable = (*Class)(nil)

type Class struct {
	name    string
	methods map[string]Func
}

func newClass(name string, methods map[string]Func) *Class {
	return &Class{
		name:    name,
		methods: methods,
	}
}

func (c *Class) String() string {
	return c.name
}

func (c *Class) Call(in *Interpreter, args []interface{}) (interface{}, error) {
	instance := NewInstance(c)

	initializer := c.FindMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(in, args)
	}

	return instance, nil
}

func (c *Class) Arity() int {
	initializer := c.FindMethod("init")
	if initializer != nil {
		return initializer.Arity()
	}
	return 0
}

func (c *Class) FindMethod(name string) *Func {
	if m, ok := c.methods[name]; ok {
		return &m
	}

	return nil
}

type Instance struct {
	Clazz *Class

	fields map[string]interface{}
}

func NewInstance(clazz *Class) Instance {
	return Instance{
		Clazz:  clazz,
		fields: make(map[string]interface{}),
	}
}

func (i *Instance) String() string {
	return i.Clazz.name + " instance"
}

func (i *Instance) Get(name *d.Token) (interface{}, error) {
	if v, ok := i.fields[name.Lexeme]; ok {
		return v, nil
	}

	method := i.Clazz.FindMethod(name.Lexeme)
	if method != nil {
		return method.Bind(*i), nil
	}

	return nil, newErrClass(name, fmt.Sprintf("Undefined property '%s'", name.Lexeme))
}

func (i *Instance) Set(name *d.Token, value interface{}) {
	i.fields[name.Lexeme] = value
}
