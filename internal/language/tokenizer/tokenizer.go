package tokenizer

import (
	"bufio"
	"strconv"
	"unicode"

	"github.com/tyrkinn/stdtest/internal/language"
)

type Token struct {
	_type    language.TokenType
	lexeme   string
	literal  any
	position uint
}

type Tokenizer struct {
	position     uint
	tokens       []Token
	sourceReader *bufio.Reader
}

func New(reader *bufio.Reader) Tokenizer {
	return Tokenizer{
		sourceReader: reader,
		position:     0,
		tokens:       make([]Token, 0, 32),
	}
}

func (t *Tokenizer) isEnd() bool {
	return t.position >= uint(t.sourceReader.Size())
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
		if t.isEnd() {
			return string(acc), nil
		}
		r, err := t.readNext()
		if err != nil {
			return "", err
		}
		if !(unicode.IsDigit(r) || r == '.') {
			err = t.sourceReader.UnreadRune()
			if err != nil {
				return "", err
			}
			return string(acc), nil
		}
		acc = append(acc, r)
	}
}

func (t *Tokenizer) readIdentifier(firstRune rune) (string, error) {
	acc := make([]rune, 1, 8)
	acc[0] = firstRune
	for {
		if t.isEnd() {
			return string(acc), nil
		}
		r, err := t.readNext()
		if err != nil {
			return "", err
		}
		if !(unicode.In(r, unicode.Digit, unicode.Letter) || r == '-' || r == '_') {
			err = t.sourceReader.UnreadRune()
			if err != nil {
				return "", err
			}
			return string(acc), nil
		}
		acc = append(acc, r)
	}
}

func (t *Tokenizer) ScanTokens() ([]Token, error) {
	for {
		if t.isEnd() {
			return t.tokens, nil
		}

		r, err := t.readNext()
		if err != nil {
			return nil, err
		}

		if unicode.IsSpace(r) {
			continue
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
			t.tokens = append(t.tokens, Token{_type: language.Number, lexeme: numberString, literal: number, position: t.position})
		}

		if unicode.IsLetter(r) || r == '_' {
			identifier, err := t.readIdentifier(r)
			if err != nil {
				return nil, err
			}
			t.tokens = append(t.tokens, Token{_type: language.Identifier, lexeme: identifier, literal: nil, position: t.position})
		}

		if r == '-' {
			next, err := t.readNext()
			if err != nil {
				return nil, err
			}
			if unicode.IsSpace(next) {
				t.tokens = append(t.tokens, Token{_type: language.MUNIS, lexeme: string(r), literal: nil, position: t.position})
			} else if next == '>' {
				t.tokens = append(t.tokens, Token{_type: language.ASSERT, lexeme: string(r), literal: nil, position: t.position})
			}
		}

		if err != nil {
			return nil, err
		}
	}
}
