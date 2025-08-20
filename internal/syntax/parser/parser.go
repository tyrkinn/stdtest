package parser

import (
	"errors"
	"fmt"
	"slices"

	"github.com/tyrkinn/stdtest/internal/syntax"
)

var EmptyStack = errors.New("Expression stack is empty")

func expectedExpr(exprType ExprType) error {
	return fmt.Errorf("Expected expression of type %s\n", exprType.String())
}

func unexpectedToken(tt syntax.TokenType) error {
	return fmt.Errorf("Unexpected expression of type %s\n", tt.String())
}

type ExprType int

const (
	CommandCall ExprType = iota
	Assert
	String
	Number
)

func (e ExprType) String() string {
	switch e {
	case Assert:
		return "Assert"
	case CommandCall:
		return "CommandCall"
	case String:
		return "String"
	case Number:
		return "Number"
	default:
		panic(fmt.Sprintf("unexpected parser.ExprType: %#v", e))
	}
}

type Expr struct {
	Type    ExprType
	Payload any
}

type CommandCallExpr struct {
	Command string
	Args    []syntax.Token
}

func (cc CommandCallExpr) String() string {
	return fmt.Sprintf("(%s %s)", cc.Command, cc.Args)
}

type TestCase struct {
	Cmd      CommandCallExpr
	Expected Expr
}

func (tc TestCase) String() string {
	return fmt.Sprintf("%s -> %v", tc.Cmd.String(), tc.Expected.Payload)
}

type Parser struct {
	position  uint
	tokens    []syntax.Token
	testCases []TestCase
	stack     []Expr
}

func New(tokens []syntax.Token) Parser {
	return Parser{
		position:  0,
		tokens:    tokens,
		testCases: make([]TestCase, 0),
		stack:     make([]Expr, 0),
	}
}

func (p *Parser) Parse() ([]TestCase, error) {
	for tok := p.peek(); tok.Type != syntax.EOF; tok = p.advance() {
		switch tok.Type {
		case syntax.NEWLINE:
		case syntax.Identifier:
			p.pushExpr(p.commandCall(tok.Lexeme))
		case syntax.ASSERT:
			testCase, err := p.assert()
			if err != nil {
				return nil, err
			}
			p.testCases = append(p.testCases, testCase)
		case syntax.Number:
			fallthrough
		case syntax.String:
			fallthrough
		default:
			return nil, unexpectedToken(tok.Type)
		}
	}
	return p.testCases, nil
}

func (p *Parser) commandCall(cmd string) Expr {
	args := make([]syntax.Token, 0)
	for p.match(syntax.Number, syntax.String) {
		args = append(args, p.previous())
	}
	return Expr{
		Type:    CommandCall,
		Payload: CommandCallExpr{Command: cmd, Args: args}}
}

func (p *Parser) assert() (TestCase, error) {
	call, err := p.tryPopExpr(CommandCall)
	if err != nil {
		return TestCase{}, err
	}
	expr := call.Payload.(CommandCallExpr)
	if p.match(syntax.Identifier) {
		rhs := p.commandCall(p.previous().Lexeme)
		return TestCase{Cmd: expr, Expected: rhs}, nil
	}
	if p.match(syntax.String, syntax.Number) {
		return TestCase{Cmd: expr, Expected: p.primitive(p.previous())}, nil
	}
	return TestCase{}, fmt.Errorf(`Unexpected tok "%s" in ASSERT rhs`, p.peek().String())
}

func (p *Parser) pushExpr(e Expr) {
	p.stack = append(p.stack, e)
}

func (p *Parser) primitive(tok syntax.Token) Expr {
	switch tok.Type {
	case syntax.Number:
		return Expr{Type: Number, Payload: tok.Literal}
	case syntax.String:
		return Expr{Type: String, Payload: tok.Literal}
	default:
		panic("unexpected syntax.TokenType for primitive")
	}

}

func (p *Parser) tryPopExpr(ofType ExprType) (Expr, error) {
	top, err := p._popExpr()
	if err != nil {
		return Expr{}, err
	}
	if top.Type != ofType {
		return Expr{}, expectedExpr(ofType)
	}
	return top, nil
}

func (p *Parser) _popExpr() (Expr, error) {
	if len(p.stack) == 0 {
		return Expr{}, EmptyStack
	}
	top := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return top, nil
}

func (p *Parser) peek() syntax.Token {
	return p.tokens[p.position]
}

func (p *Parser) previous() syntax.Token {
	return p.tokens[p.position-1]
}

func (p *Parser) match(tts ...syntax.TokenType) bool {
	if slices.ContainsFunc(tts, p.check) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) check(t syntax.TokenType) bool {
	return !p.isEnd() && p.peek().Type == t
}

func (p *Parser) isEnd() bool {
	return p.peek().Type == syntax.EOF
}

func (p *Parser) advance() syntax.Token {
	if !p.isEnd() {
		p.position++
		return p.previous()
	}
	return syntax.Token{Type: syntax.EOF}
}
