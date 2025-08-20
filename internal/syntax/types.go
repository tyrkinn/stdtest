package syntax

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
	NEWLINE
	EOF
)

func (tt TokenType) String() string {
	switch tt {
	case ASSERT:
		return "ASSERT"
	case Identifier:
		return "IDENTIFIER"
	case Number:
		return "NUMBER"
	case String:
		return "STRING"
	case NEWLINE:
		return "NEWLINE"
	case EOF:
		return "EOF"
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

func (t Token) String() string {
	return fmt.Sprintf("{%s, %s :at %d}", t.Type.String(), strconv.Quote(t.Lexeme), t.Position)
}
