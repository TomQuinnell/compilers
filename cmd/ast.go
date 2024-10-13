package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// Each string generates a struct for an AST node
// Format: "[struct-name] : [field-name] [field-type], ..."

func main() {
	writeAst("Expr", []string{
		"Unary    : Operator *Token, Right Expr",
		"Assign   : Name *Token, Value Expr",
		"Binary   : Left Expr, Operator *Token, Right Expr",
		"Call     : Callee Expr, Paren *Token, Args []Expr",
		"Literal  : Value interface{}",
		"Logical  : Left Expr, Operator *Token, Right Expr",
		"Grouping : Expression Expr",
		"Variable : Name *Token",
	}, true)

	writeAst("Stmt", []string{
		"Block      : Stmts []Stmt",
		"Expression : Expression Expr",
		"Function   : Name *Token, Params []*Token, Body []Stmt",
		"If         : Condition Expr, ThenBranch Stmt, ElseBranch Stmt",
		"Print      : Expression Expr",
		"Return     : Keyword *Token, Value Expr",
		"Var        : Name *Token, Initializer Expr",
		"While      : Condition Expr, Body Stmt",
	}, false)
}

// writeAst("Stmt", []string{
// 	"Block      : Statements []Stmt",
// 	"Expression : Expr Expr",
// 	"Var        : Name Token, Initializer Expr",
// })

func writeAst(baseName string, types []string, hasReturnValue bool) {
	ret := ""

	ret += "package domain\n"

	ret += defineInterface(baseName, hasReturnValue)
	ret += defineTypes(baseName, types, hasReturnValue)
	ret += defineVisitor(baseName, types, hasReturnValue)

	filename := fmt.Sprintf("./domain/%s.go", strings.ToLower(baseName))
	err := os.WriteFile(filename, []byte(ret), 0655)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to write AST")
	}
}

func getReturnType(hasReturnValue bool) string {
	ret := "error"
	if hasReturnValue {
		ret = "(interface{}, error)"
	}
	return ret
}

func defineInterface(name string, hasReturnValue bool) string {
	return fmt.Sprintf(`
type %s interface {
	Accept(visitor %sVisitor) %s
}
`, name, name, getReturnType(hasReturnValue))
}

func defineTypes(name string, types []string, hasReturnValue bool) (str string) {
	for _, t := range types {
		splitType := strings.Split(t, ":")
		fullTypeName := strings.Trim(splitType[0], " ") + name
		str += fmt.Sprintf("\ntype %s struct {\n", fullTypeName)

		fields := strings.Split(splitType[1], ", ")
		for _, field := range fields {
			str += fmt.Sprintf("\t%s\n", strings.Trim(field, " "))
		}

		str += "}\n"

		str += fmt.Sprintf(`
func (b %s) Accept(visitor %sVisitor) %s {
	return visitor.Visit%s(b)
}
`, fullTypeName, name, getReturnType(hasReturnValue), fullTypeName)
	}
	return str
}

func defineVisitor(name string, types []string, hasReturnValue bool) (str string) {
	str += fmt.Sprintf("\ntype %sVisitor interface {\n", name)
	for _, t := range types {
		splitType := strings.Split(t, ":")
		fullTypeName := strings.Trim(splitType[0], " ") + name
		str += fmt.Sprintf("\tVisit%s(expr %s) %s\n",
			fullTypeName, fullTypeName, getReturnType(hasReturnValue))
	}
	str += "}\n"
	return str
}
