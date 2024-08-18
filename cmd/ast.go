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
		"Unary    : Operator Token, Right Expr[R]",
		"Binary   : Left Expr[R], Operator Token, Right Expr[R]",
		"Literal  : Value fmt.Stringer",
		"Grouping : Expression Expr[R]",
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
	ret += "\nimport \"fmt\"\n"

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
type %s[R any] interface {
	Accept(visitor %sVisitor[R]) R
}
`, name, name)
}

func defineTypes(name string, types []string) (str string) {
	for _, t := range types {
		splitType := strings.Split(t, ":")
		fullTypeName := strings.Trim(splitType[0], " ") + name
		str += fmt.Sprintf("\ntype %s[R any] struct {\n", fullTypeName)

		fields := strings.Split(splitType[1], ", ")
		for _, field := range fields {
			str += fmt.Sprintf("\t%s\n", strings.Trim(field, " "))
		}

		str += "}\n"

		str += fmt.Sprintf(`
func (b %s[R]) Accept(visitor %sVisitor[R]) R {
	return visitor.Visit%s(b)
}
`, fullTypeName, name, fullTypeName)
	}
	return str
}

func defineVisitor(name string, types []string) (str string) {
	str += fmt.Sprintf("\ntype %sVisitor[R any] interface {\n", name)
	for _, t := range types {
		splitType := strings.Split(t, ":")
		fullTypeName := strings.Trim(splitType[0], " ") + name
		str += fmt.Sprintf("\tVisit%s(expr %s[R]) R\n", fullTypeName, fullTypeName)
	}
	str += "}\n"
	return str
}
