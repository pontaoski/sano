package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var Parser = participle.MustBuild[File](
	Addresses,
	Statements,
	Expressions,
	participle.UseLookahead(3),
)

type File struct {
	Fragment []Fragment `@@*`
}

type Fragment struct {
	Pos lexer.Position

	Name       string      `"@" @Ident "{"`
	Statements []Statement `(@@)* "}"`
}
