package tokenizer

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"

	"github.com/tyrkinn/stdtest/internal/syntax"
)

type runeReader interface {
	ReadRune() (rune, int, error)
	UnreadRune() error
}

type Tokenizer struct {
	position     uint
	tokens       []syntax.Token
	sourceReader runeReader
}

func New(reader runeReader) Tokenizer {
	return Tokenizer{
		sourceReader: reader,
		position:     0,
		tokens:       make([]syntax.Token, 0, 32),
	}
}

func (t *Tokenizer) readNext() (rune, error) {
	r, _, err := t.sourceReader.ReadRune()
	if err != nil {
		return 0, err
	}
	t.position++
	return r, nil
}

func (t *Tokenizer) ScanTokens() ([]syntax.Token, error) {
	for {
		pos := t.position
		r, err := t.readNext()
		if errors.Is(err, io.EOF) {
			t.addEOF()
			return t.tokens, nil
		}
		if err != nil {
			return nil, err
		}

		switch {

		case r == ' ' || r == '\t':

		case r == '\n':
			t.tokens = append(t.tokens, syntax.Token{Type: syntax.NEWLINE, Lexeme: "\n", Literal: nil, Position: pos})

		case r == '\'' || r == '"':
			str, err := t.readString(r)
			if err != nil {
				return nil, err
			}
			t.tokens = append(t.tokens, syntax.Token{Type: syntax.String, Lexeme: str, Literal: str, Position: pos})

		case r == '-':
			token, err := t.readAssert(pos)
			if err != nil {
				return nil, err
			}
			t.tokens = append(t.tokens, token)

		case unicode.IsDigit(r):
			numberString, err := t.readNumber(r)
			if err != nil {
				return nil, err
			}
			number, err := strconv.Atoi(numberString)
			if err != nil {
				return nil, err
			}
			t.tokens = append(t.tokens, syntax.Token{Type: syntax.Number, Lexeme: numberString, Literal: number, Position: pos})

		case unicode.IsLetter(r) || r == '_':
			identifier, err := t.readIdentifier(r)
			if err != nil {
				return nil, err
			}
			t.tokens = append(t.tokens, syntax.Token{Type: syntax.Identifier, Lexeme: identifier, Literal: nil, Position: pos})

		default:
			panic(fmt.Errorf("Unexpected lexem %c at %d", r, pos))
		}
	}
}

func (t *Tokenizer) readNumber(firstRune rune) (string, error) {
	acc := make([]rune, 1, 8)
	acc[0] = firstRune
	for {
		r, err := t.readNext()
		if errors.Is(err, io.EOF) {
			return string(acc), nil
		}
		if err != nil {
			return "", err
		}
		if !(unicode.IsDigit(r) || r == '.') {
			err = t.sourceReader.UnreadRune()
			if err != nil {
				return "", err
			}
			t.position--
			return string(acc), nil
		}
		acc = append(acc, r)
	}
}

func (t *Tokenizer) readString(quote rune) (string, error) {
	acc := make([]rune, 0, 8)
	for {
		r, err := t.readNext()
		if errors.Is(err, io.EOF) {
			return "", fmt.Errorf("Unexpected EOF while reading string")
		}
		if err != nil {
			return "", err
		}
		if r == quote {
			return string(acc), nil
		}
		acc = append(acc, r)
	}
}

func (t *Tokenizer) readIdentifier(firstRune rune) (string, error) {
	acc := make([]rune, 1, 8)
	acc[0] = firstRune
	for {
		r, err := t.readNext()
		if errors.Is(err, io.EOF) {
			return string(acc), nil
		}
		if err != nil {
			return "", err
		}
		if !(unicode.In(r, unicode.Digit, unicode.Letter) || r == '-' || r == '_') {
			err = t.sourceReader.UnreadRune()
			if err != nil {
				return "", err
			}
			t.position--
			return string(acc), nil
		}
		acc = append(acc, r)
	}
}

func (t *Tokenizer) readAssert(pos uint) (syntax.Token, error) {
	next, err := t.readNext()
	if err != nil {
		return syntax.Token{}, err
	}
	if next == '>' {
		return syntax.Token{Type: syntax.ASSERT, Lexeme: "->", Literal: nil, Position: pos}, nil
	}
	return syntax.Token{}, fmt.Errorf("Unexpected `-`")
}

func (t *Tokenizer) addEOF() {
	t.tokens = append(t.tokens, syntax.Token{Type: syntax.EOF, Position: t.position})
}
