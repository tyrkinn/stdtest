package language

import (
	"fmt"
	"strconv"
)

type TokenType int

const (
	Identifier TokenType = iota
	String
	Number

	ASSERT
	MUNIS
	NEWLINE
)

func TokenTypeToString(tt TokenType) string {
	switch tt {
	case ASSERT:
		return "ASSERT"
	case Identifier:
		return "IDENTIFIER"
	case MUNIS:
		return "MINUS"
	case Number:
		return "NUMBER"
	case String:
		return "STRING"
	case NEWLINE:
		return "NEWLINE"
	default:
		panic(fmt.Sprintf("unexpected language.TokenType: %#v", tt))
	}
}

type Token struct {
	Type     TokenType
	Lexeme   string
	Literal  any
	Position uint
}

func TokenToString(t Token) string {
	return fmt.Sprintf("{%s, %s :at %d}", TokenTypeToString(t.Type), strconv.Quote(t.Lexeme), t.Position)
}
