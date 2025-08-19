package tokenizer

import (
	"errors"
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

func (t *Tokenizer) ScanTokens() ([]syntax.Token, error) {
	for {
		pos := t.position
		r, err := t.readNext()
		if errors.Is(err, io.EOF) {
			return t.tokens, nil
		}
		if err != nil {
			return nil, err
		}

		if r == ' ' || r == '\t' {
			continue
		}

		if r == '\n' {
			t.tokens = append(t.tokens, syntax.Token{Type: syntax.NEWLINE, Lexeme: "\n", Literal: nil, Position: pos})
		}

		if unicode.IsDigit(r) {
			numberString, err := t.readNumber(r)
			if err != nil {
				return nil, err
			}
			number, err := strconv.ParseFloat(numberString, 64)
			if err != nil {
				return nil, err
			}
			t.tokens = append(t.tokens, syntax.Token{Type: syntax.Number, Lexeme: numberString, Literal: number, Position: pos})
		}

		if unicode.IsLetter(r) || r == '_' {
			identifier, err := t.readIdentifier(r)
			if err != nil {
				return nil, err
			}
			t.tokens = append(t.tokens, syntax.Token{Type: syntax.Identifier, Lexeme: identifier, Literal: nil, Position: pos})
		}

		if r == '-' {
			next, err := t.readNext()
			if err != nil {
				return nil, err
			}
			if next == '>' {
				t.tokens = append(t.tokens, syntax.Token{Type: syntax.ASSERT, Lexeme: "->", Literal: nil, Position: pos})
			} else {
				t.tokens = append(t.tokens, syntax.Token{Type: syntax.MUNIS, Lexeme: string(r), Literal: nil, Position: pos})
				err = t.sourceReader.UnreadRune()
				if err != nil {
					return nil, err
				}
			}
		}

	}
}
