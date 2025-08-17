package tokenizer

import (
	"errors"
	"io"
	"strconv"
	"unicode"

	"github.com/tyrkinn/stdtest/internal/language"
)

type runeReader interface {
	ReadRune() (rune, int, error)
	UnreadRune() error
}

type Tokenizer struct {
	position     uint
	tokens       []language.Token
	sourceReader runeReader
}

func New(reader runeReader) Tokenizer {
	return Tokenizer{
		sourceReader: reader,
		position:     0,
		tokens:       make([]language.Token, 0, 32),
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

func (t *Tokenizer) ScanTokens() ([]language.Token, error) {
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
			t.tokens = append(t.tokens, language.Token{Type: language.NEWLINE, Lexeme: "\n", Literal: nil, Position: pos})
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
			t.tokens = append(t.tokens, language.Token{Type: language.Number, Lexeme: numberString, Literal: number, Position: pos})
		}

		if unicode.IsLetter(r) || r == '_' {
			identifier, err := t.readIdentifier(r)
			if err != nil {
				return nil, err
			}
			t.tokens = append(t.tokens, language.Token{Type: language.Identifier, Lexeme: identifier, Literal: nil, Position: pos})
		}

		if r == '-' {
			next, err := t.readNext()
			if err != nil {
				return nil, err
			}
			if unicode.IsSpace(next) {
				t.tokens = append(t.tokens, language.Token{Type: language.MUNIS, Lexeme: string(r), Literal: nil, Position: pos})
			} else if next == '>' {
				t.tokens = append(t.tokens, language.Token{Type: language.ASSERT, Lexeme: "->", Literal: nil, Position: pos})
			}
		}

	}
}
