package env

import (
	d "example/compilers/domain"
	"fmt"
)

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewEnv(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name *d.Token) (interface{}, error) {
	if v, ok := e.values[name.Lexeme]; ok {
		return v, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, fmt.Errorf("env value '%s' not found", name.Lexeme) // TODO: runtime err?
}

func (e *Environment) GetAt(distance int, name string) (interface{}, error) {
	return e.ancestor(distance).values[name], nil
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for range distance {
		env = env.enclosing
	}
	return env
}

func (e *Environment) Assign(name *d.Token, v interface{}) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = v
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, v)
	}

	return fmt.Errorf("variable '%s' undefined", name.Lexeme)
}

func (e *Environment) AssignAt(distance int, name *d.Token, v interface{}) {
	e.ancestor(distance).values[name.Lexeme] = v
}

// TODO: tests? Ehh no
