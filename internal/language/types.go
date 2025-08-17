package language

type TokenType int

const (
	Identifier TokenType = iota
	String
	Number

	ASSERT
	MUNIS
)
