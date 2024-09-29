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
		"Binary   : Left Expr, Operator *Token, Right Expr",
		"Literal  : Value interface{}",
		"Grouping : Expression Expr",
		"Variable : Name Token",
	})

	writeAst("Stmt", []string{
		"Expression : Expression Expr",
		"Print      : Expression Expr",
		"Var        : Name Token, Initializer Expr",
	})
}

// writeAst("Stmt", []string{
// 	"Block      : Statements []Stmt",
// 	"Expression : Expr Expr",
// 	"Var        : Name Token, Initializer Expr",
// })

func writeAst(baseName string, types []string) {
	ret := ""

	ret += "package domain\n"

	ret += defineInterface(baseName)
	ret += defineTypes(baseName, types)
	ret += defineVisitor(baseName, types)

	filename := fmt.Sprintf("./domain/%s.go", strings.ToLower(baseName))
	err := os.WriteFile(filename, []byte(ret), 0655)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to write AST")
	}
}

func defineInterface(name string) string {
	return fmt.Sprintf(`
type %s interface {
	Accept(visitor %sVisitor) (interface{}, error)
}
`, name, name)
}

func defineTypes(name string, types []string) (str string) {
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
func (b %s) Accept(visitor %sVisitor) (interface{}, error) {
	return visitor.Visit%s(b)
}
`, fullTypeName, name, fullTypeName)
	}
	return str
}

func defineVisitor(name string, types []string) (str string) {
	str += fmt.Sprintf("\ntype %sVisitor interface {\n", name)
	for _, t := range types {
		splitType := strings.Split(t, ":")
		fullTypeName := strings.Trim(splitType[0], " ") + name
		str += fmt.Sprintf("\tVisit%s(expr %s) (interface{}, error)\n", fullTypeName, fullTypeName)
	}
	str += "}\n"
	return str
}
