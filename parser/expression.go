package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var Expressions = participle.Union[Expression](
	NumericLiteral{},
	Symbol{},
)

type Expression interface {
	Position() lexer.Position
	isExpression()
}

type NumericLiteral struct {
	Pos    lexer.Position
	Number int `@Int`
}

func (n NumericLiteral) Position() lexer.Position {
	return n.Pos
}

func (NumericLiteral) isExpression() {}

type Symbol struct {
	Pos  lexer.Position
	Name string `@Ident`
}

func (s Symbol) Position() lexer.Position {
	return s.Pos
}

func (Symbol) isExpression() {}
