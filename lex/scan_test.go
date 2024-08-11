package lex

import (
	d "example/compilers/domain"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	type ScanTestCase struct {
		rawText       string
		expectedToken d.TokenType
	}

	testCases := []ScanTestCase{
		{"(", d.LEFT_PAREN},
		{")", d.RIGHT_PAREN},
		{"{", d.LEFT_BRACE},
		{"}", d.RIGHT_BRACE},
		{",", d.COMMA},
		{".", d.DOT},
		{"-", d.MINUS},
		{"+", d.PLUS},
		{";", d.SEMICOLON},
		{"*", d.STAR},
		{"!", d.BANG},
		{"!=", d.BANG_EQUAL},
		{"=", d.EQUAL},
		{"==", d.EQUAL_EQUAL},
		{"<", d.LESS},
		{"<=", d.LESS_EQUAL},
		{">", d.GREATER},
		{">=", d.GREATER_EQUAL},
		{"/", d.SLASH},
		{"i", d.IDENTIFIER},
		{"and", d.AND},
		{"class", d.CLASS},
		{"else", d.ELSE},
		{"false", d.FALSE},
		{"for", d.FOR},
		{"fun", d.FUN},
		{"if", d.IF},
		{"nil", d.NIL},
		{"or", d.OR},
		{"print", d.PRINT},
		{"return", d.RETURN},
		{"super", d.SUPER},
		{"this", d.THIS},
		{"true", d.TRUE},
		{"var", d.VAR},
		{"while", d.WHILE},
	}

	for _, c := range testCases {
		t.Run(fmt.Sprintf("Scans raw token: %s", c.rawText), func(t *testing.T) {
			assert := assert.New(t)

			scannedTokens, err := NewScanner(c.rawText).Scan()
			assert.NoError(err)

			assert.Len(scannedTokens, 2)
			assert.Equal(c.expectedToken, scannedTokens[0].Kind)
			assert.Equal(d.EOF, scannedTokens[1].Kind)
		})
	}

	wsTestCases := []ScanTestCase{
		{" ", d.EOF},
		{"\r", d.EOF},
		{"\t", d.EOF},
		{"\n", d.EOF},
		{"// comment", d.EOF},
	}

	for _, c := range wsTestCases {
		t.Run(fmt.Sprintf("Doesn't scan whitespace + comments token: %s", c.rawText), func(t *testing.T) {
			assert := assert.New(t)

			scannedTokens, err := NewScanner(c.rawText).Scan()
			assert.NoError(err)

			assert.Len(scannedTokens, 1)
			assert.Equal(c.expectedToken, scannedTokens[0].Kind)
		})
	}

	type ScanLiteralTestCase struct {
		rawText       string
		expectedToken d.Token
	}

	literalTestCases := []ScanLiteralTestCase{
		{"\"hello\"", d.Token{Kind: d.STRING, Literal: "hello"}},
		{"\"\"", d.Token{Kind: d.STRING, Literal: ""}},
		{"5", d.Token{Kind: d.NUMBER, Literal: 5.0}},
		{"5.01", d.Token{Kind: d.NUMBER, Literal: 5.01}},
		{"5.0005", d.Token{Kind: d.NUMBER, Literal: 5.0005}},
	}

	for _, c := range literalTestCases {
		t.Run(fmt.Sprintf("Scans token with literal: %s", c.rawText), func(t *testing.T) {
			assert := assert.New(t)

			scannedTokens, err := NewScanner(c.rawText).Scan()
			assert.NoError(err)

			assert.Len(scannedTokens, 2)
			assert.Equal(c.expectedToken.Kind, scannedTokens[0].Kind)
			assert.Equal(c.expectedToken.Literal, scannedTokens[0].Literal)
			assert.Equal(d.EOF, scannedTokens[1].Kind)
		})
	}
}
