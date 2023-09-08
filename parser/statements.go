package parser

import "github.com/alecthomas/participle/v2"

var Statements = participle.Union[Statement](
	SymbolDeclaration{},
	OpcodeInvocation{},
)

type Statement interface {
	isStatement()
}

type SymbolDeclaration struct {
	Name string `"&" @Ident ":"`
}

func (SymbolDeclaration) isStatement() {}
