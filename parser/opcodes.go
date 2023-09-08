package parser

import (
	"Sano/cpu"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type OpcodeInvocation struct {
	Pos lexer.Position

	Opcode  string  `@Ident`
	Address Address `@@ ";"`
}

func (OpcodeInvocation) isStatement() {

}

var Addresses = participle.Union[Address](
	AddrAccumulator{},
	AddrImmediate{},
	AddrAbsoluteIndirect{},
	AddrAbsoluteIndexedX{},
	AddrAbsoluteIndexedY{},
	AddrAbsoluteAddress{},
	AddrZeroPageXIndexedIndirect{},
	AddrZeroPageIndirectYIndex{},
	AddrZeroPageXIndexed{},
	AddrZeroPageYIndexed{},
	AddrZeroPage{},
	AddrRelative{},
	AddrImplied{},
)

type Address interface {
	AddressingMode() cpu.Mode
}

type AddrImplied struct {
	Dummy struct{} `"!"`
}

func (AddrImplied) AddressingMode() cpu.Mode { return cpu.Implied }

type AddrAccumulator struct {
	Dummy struct{} `"a"`
}

func (AddrAccumulator) AddressingMode() cpu.Mode { return cpu.Accumulator }

type AddrImmediate struct {
	Value Expression `"#" @@`
}

func (AddrImmediate) AddressingMode() cpu.Mode { return cpu.Immediate }

type AddrAbsoluteIndirect struct {
	Address Expression `"(" "=" @@ ")"`
}

func (AddrAbsoluteIndirect) AddressingMode() cpu.Mode { return cpu.Indirect }

type AddrAbsoluteAddress struct {
	Address Expression `"=" @@`
}

func (AddrAbsoluteAddress) AddressingMode() cpu.Mode { return cpu.Absolute }

type AddrAbsoluteIndexedX struct {
	Address Expression `"=" @@ "," "x"`
}

func (AddrAbsoluteIndexedX) AddressingMode() cpu.Mode { return cpu.AbsoluteIndexedX }

type AddrAbsoluteIndexedY struct {
	Address Expression `"=" @@ "," "y"`
}

func (AddrAbsoluteIndexedY) AddressingMode() cpu.Mode { return cpu.AbsoluteIndexedY }

type AddrZeroPage struct {
	Address Expression `":" @@`
}

func (AddrZeroPage) AddressingMode() cpu.Mode { return cpu.ZeroPage }

type AddrZeroPageXIndexed struct {
	Address Expression `":" @@ "," "x"`
}

func (AddrZeroPageXIndexed) AddressingMode() cpu.Mode { return cpu.ZeroPageIndexedX }

type AddrZeroPageYIndexed struct {
	Address Expression `":" @@ "," "y"`
}

func (AddrZeroPageYIndexed) AddressingMode() cpu.Mode { return cpu.ZeroPageIndexedY }

type AddrZeroPageXIndexedIndirect struct {
	Address Expression `"(" ":" @@ "," "x" ")"`
}

func (AddrZeroPageXIndexedIndirect) AddressingMode() cpu.Mode { return cpu.XIndexedIndirect }

type AddrZeroPageIndirectYIndex struct {
	Address Expression `"(" ":" @@ ")" "," "y"`
}

func (AddrZeroPageIndirectYIndex) AddressingMode() cpu.Mode { return cpu.IndirectYIndexed }

type AddrRelative struct {
	Sign    string     `"~"`
	Address Expression `@@`
}

func (AddrRelative) AddressingMode() cpu.Mode { return cpu.Relative }
