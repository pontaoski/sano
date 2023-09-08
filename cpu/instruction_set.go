package cpu

type Opcode byte

const (
	ADC Opcode = iota
	AND
	ASL
	BRA
	BCC // AKA BLT
	BCS // AKA BGT
	BEQ
	BIT
	BMI
	BNE
	BPL
	BRK
	BVC
	BVS
	CLC
	CLD
	CLI
	CLV
	CMP
	CPX
	CPY
	DEC
	DEX
	DEY
	EOR
	INC
	INX
	INY
	JMP
	JSR
	LDA
	LDX
	LDY
	LSR
	NOP
	ORA
	PHA
	PHP
	PLA
	PLP
	ROL
	ROR
	RTI
	RTS
	SBC
	SEC
	SED
	SEI
	STA
	STX
	STY
	STZ
	TAX
	TAY
	TSX
	TXA
	TXS
	TYA
	TRB
	TSB
	PHX
	PLX
	PHY
	PLY
)

var OpcodeNames = map[string]Opcode{
	"adc": ADC,
	"and": AND,
	"asl": ASL,
	"bra": BRA,
	"bcc": BCC,
	"bcs": BCS,
	"beq": BEQ,
	"bit": BIT,
	"bmi": BMI,
	"bne": BNE,
	"bpl": BPL,
	"brk": BRK,
	"bvc": BVC,
	"bvs": BVS,
	"clc": CLC,
	"cld": CLD,
	"cli": CLI,
	"clv": CLV,
	"cmp": CMP,
	"cpx": CPX,
	"cpy": CPY,
	"dec": DEC,
	"dex": DEX,
	"dey": DEY,
	"eor": EOR,
	"inc": INC,
	"inx": INX,
	"iny": INY,
	"jmp": JMP,
	"jsr": JSR,
	"lda": LDA,
	"ldx": LDX,
	"ldy": LDY,
	"lsr": LSR,
	"nop": NOP,
	"ora": ORA,
	"pha": PHA,
	"php": PHP,
	"pla": PLA,
	"plp": PLP,
	"rol": ROL,
	"ror": ROR,
	"rti": RTI,
	"rts": RTS,
	"sbc": SBC,
	"sec": SEC,
	"sed": SED,
	"sei": SEI,
	"sta": STA,
	"stx": STX,
	"sty": STY,
	"stz": STZ,
	"tax": TAX,
	"tay": TAY,
	"tsx": TSX,
	"txa": TXA,
	"txs": TXS,
	"tya": TYA,
	"trb": TRB,
	"tsb": TSB,
	"phx": PHX,
	"plx": PLX,
	"phy": PHY,
	"ply": PLY,
}

type Mode byte

const (
	Immediate Mode = iota
	Implied
	Relative
	Accumulator

	ZeroPage
	ZeroPageIndexedX
	ZeroPageIndexedY

	Absolute
	AbsoluteIndexedX
	AbsoluteIndexedY

	Indirect
	XIndexedIndirect
	IndirectYIndexed
)

func (s Mode) String() string {
	switch s {
	case Immediate:
		return "immediate"
	case Implied:
		return "implied"
	case Relative:
		return "relative"
	case Accumulator:
		return "accumulator"
	case ZeroPage:
		return "zero page"
	case ZeroPageIndexedX:
		return "x-indexed zero page"
	case ZeroPageIndexedY:
		return "y-indexed zero page"
	case Absolute:
		return "absolute"
	case AbsoluteIndexedX:
		return "x-indexed absolute"
	case AbsoluteIndexedY:
		return "y-indexed absolute"
	case Indirect:
		return "indirect"
	case XIndexedIndirect:
		return "x-indexed indirect zero page"
	case IndirectYIndexed:
		return "indirect y-indexed zero page"
	default:
		panic("invalid addressing mode")
	}
}

type OpcodeData struct {
	Operation Opcode
	Mode      Mode

	Hex byte
}

func (o OpcodeSet) And(other OpcodeSet) OpcodeSet {
	return append(o, other...)
}

func (o OpcodeSet) Find(of Opcode) OpcodeSet {
	ret := OpcodeSet{}
	for _, op := range o {
		if op.Operation == of {
			ret = append(ret, op)
		}
	}
	return ret
}

func (o OpcodeSet) FindOne(of Opcode, mode Mode) (OpcodeData, bool) {
	for _, op := range o {
		if op.Operation == of && op.Mode == mode {
			return op, true
		}
	}
	return OpcodeData{}, false
}

type OpcodeSet []OpcodeData

// The 6502 opcodes
var Base6502Opcodes = OpcodeSet{
	{LDA, Immediate, 0xA9},
	{LDA, ZeroPage, 0xA9},
	{LDA, ZeroPageIndexedX, 0xB5},
	{LDA, Absolute, 0xAD},
	{LDA, AbsoluteIndexedX, 0xBD},
	{LDA, AbsoluteIndexedY, 0xB9},
	{LDA, XIndexedIndirect, 0xA1},
	{LDA, IndirectYIndexed, 0xB1},

	{LDX, Immediate, 0xA2},
	{LDX, ZeroPage, 0xA6},
	{LDX, ZeroPageIndexedY, 0xB6},
	{LDX, Absolute, 0xAE},
	{LDX, AbsoluteIndexedY, 0xBE},

	{LDY, Immediate, 0xa0},
	{LDY, ZeroPage, 0xa4},
	{LDY, ZeroPageIndexedX, 0xb4},
	{LDY, Absolute, 0xac},
	{LDY, AbsoluteIndexedX, 0xbc},

	{STA, ZeroPage, 0x85},
	{STA, ZeroPageIndexedX, 0x95},
	{STA, Absolute, 0x8d},
	{STA, AbsoluteIndexedX, 0x9d},
	{STA, AbsoluteIndexedY, 0x99},
	{STA, XIndexedIndirect, 0x81},
	{STA, IndirectYIndexed, 0x91},
	{STA, Indirect, 0x92},

	{STX, ZeroPage, 0x86},
	{STX, ZeroPageIndexedY, 0x96},
	{STX, Absolute, 0x8e},

	{STY, ZeroPage, 0x84},
	{STY, ZeroPageIndexedX, 0x94},
	{STY, Absolute, 0x8c},

	{ADC, Immediate, 0x69},
	{ADC, ZeroPage, 0x65},
	{ADC, ZeroPageIndexedX, 0x75},
	{ADC, Absolute, 0x6d},
	{ADC, AbsoluteIndexedX, 0x7d},
	{ADC, AbsoluteIndexedY, 0x79},
	{ADC, XIndexedIndirect, 0x61},
	{ADC, IndirectYIndexed, 0x71},

	{SBC, Immediate, 0xe9},
	{SBC, ZeroPage, 0xe5},
	{SBC, ZeroPageIndexedX, 0xf5},
	{SBC, Absolute, 0xed},
	{SBC, AbsoluteIndexedX, 0xfd},
	{SBC, AbsoluteIndexedY, 0xf9},
	{SBC, XIndexedIndirect, 0xe1},
	{SBC, IndirectYIndexed, 0xf1},

	{CMP, Immediate, 0xc9},
	{CMP, ZeroPage, 0xc5},
	{CMP, ZeroPageIndexedX, 0xd5},
	{CMP, Absolute, 0xcd},
	{CMP, AbsoluteIndexedX, 0xdd},
	{CMP, AbsoluteIndexedY, 0xd9},
	{CMP, XIndexedIndirect, 0xc1},
	{CMP, IndirectYIndexed, 0xd1},

	{CPX, Immediate, 0xe0},
	{CPX, ZeroPage, 0xe4},
	{CPX, Absolute, 0xec},

	{CPY, Immediate, 0xc0},
	{CPY, ZeroPage, 0xc4},
	{CPY, Absolute, 0xcc},

	{BIT, ZeroPage, 0x24},
	{BIT, Absolute, 0x2c},

	{CLC, Implied, 0x18},
	{SEC, Implied, 0x38},
	{CLI, Implied, 0x58},
	{SEI, Implied, 0x78},
	{CLD, Implied, 0xd8},
	{SED, Implied, 0xf8},
	{CLV, Implied, 0xb8},

	{BCC, Relative, 0x90},
	{BCS, Relative, 0xb0},
	{BEQ, Relative, 0xf0},
	{BNE, Relative, 0xd0},
	{BMI, Relative, 0x30},
	{BPL, Relative, 0x10},
	{BVC, Relative, 0x50},
	{BVS, Relative, 0x70},

	{BRK, Implied, 0x00},

	{AND, Immediate, 0x29},
	{AND, ZeroPage, 0x25},
	{AND, ZeroPageIndexedX, 0x35},
	{AND, Absolute, 0x2d},
	{AND, AbsoluteIndexedX, 0x3d},
	{AND, AbsoluteIndexedY, 0x39},
	{AND, XIndexedIndirect, 0x21},
	{AND, IndirectYIndexed, 0x31},

	{ORA, Immediate, 0x09},
	{ORA, ZeroPage, 0x05},
	{ORA, ZeroPageIndexedX, 0x15},
	{ORA, Absolute, 0x0d},
	{ORA, AbsoluteIndexedX, 0x1d},
	{ORA, AbsoluteIndexedY, 0x19},
	{ORA, XIndexedIndirect, 0x01},
	{ORA, IndirectYIndexed, 0x11},

	{EOR, Immediate, 0x49},
	{EOR, ZeroPage, 0x45},
	{EOR, ZeroPageIndexedX, 0x55},
	{EOR, Absolute, 0x4d},
	{EOR, AbsoluteIndexedX, 0x5d},
	{EOR, AbsoluteIndexedY, 0x59},
	{EOR, XIndexedIndirect, 0x41},
	{EOR, IndirectYIndexed, 0x51},

	{INC, ZeroPage, 0xe6},
	{INC, ZeroPageIndexedX, 0xf6},
	{INC, Absolute, 0xee},
	{INC, AbsoluteIndexedX, 0xfe},

	{DEC, ZeroPage, 0xc6},
	{DEC, ZeroPageIndexedX, 0xd6},
	{DEC, Absolute, 0xce},
	{DEC, AbsoluteIndexedX, 0xde},

	{INX, Implied, 0xe8},
	{INY, Implied, 0xc8},

	{DEX, Implied, 0xca},
	{DEY, Implied, 0x88},

	{JMP, Absolute, 0x4c},
	{JMP, Indirect, 0x6c},

	{JSR, Absolute, 0x20},
	{RTS, Implied, 0x60},

	{RTI, Implied, 0x40},

	{NOP, Implied, 0xea},

	{TAX, Implied, 0xaa},
	{TXA, Implied, 0x8a},
	{TAY, Implied, 0xa8},
	{TYA, Implied, 0x98},
	{TXS, Implied, 0x9a},
	{TSX, Implied, 0xba},

	{PHA, Implied, 0x48},
	{PLA, Implied, 0x68},
	{PHP, Implied, 0x08},
	{PLP, Implied, 0x28},

	{ASL, Accumulator, 0x0a},
	{ASL, ZeroPage, 0x06},
	{ASL, ZeroPageIndexedX, 0x16},
	{ASL, Absolute, 0x0e},
	{ASL, AbsoluteIndexedX, 0x1e},

	{LSR, Accumulator, 0x4a},
	{LSR, ZeroPage, 0x46},
	{LSR, ZeroPageIndexedX, 0x56},
	{LSR, Absolute, 0x4e},
	{LSR, AbsoluteIndexedX, 0x5e},

	{ROL, Accumulator, 0x2a},
	{ROL, ZeroPage, 0x26},
	{ROL, ZeroPageIndexedX, 0x36},
	{ROL, Absolute, 0x2e},
	{ROL, AbsoluteIndexedX, 0x3e},

	{ROR, Accumulator, 0x6a},
	{ROR, ZeroPage, 0x66},
	{ROR, ZeroPageIndexedX, 0x76},
	{ROR, Absolute, 0x6e},
	{ROR, AbsoluteIndexedX, 0x7e},
}

// The opcodes found on the WDC 65C02
var WDC65C02ExtensionOpcodes = OpcodeSet{
	{LDA, Indirect, 0xB2},

	{STZ, ZeroPage, 0x64},
	{STZ, ZeroPageIndexedX, 0x74},
	{STZ, Absolute, 0x9c},
	{STZ, AbsoluteIndexedX, 0x9e},

	{ADC, Indirect, 0x72},

	{SBC, Indirect, 0xf2},

	{CMP, Indirect, 0xd2},

	{BIT, Immediate, 0x89},
	{BIT, ZeroPageIndexedX, 0x34},
	{BIT, AbsoluteIndexedX, 0x3c},

	{BRA, Relative, 0x80},

	{AND, Indirect, 0x32},

	{ORA, Indirect, 0x12},

	{EOR, Indirect, 0x52},

	{INC, Accumulator, 0x1a},

	{DEC, Accumulator, 0x3a},

	{JMP, AbsoluteIndexedX, 0x7c},

	{TRB, ZeroPage, 0x14},
	{TRB, Absolute, 0x1c},

	{TSB, ZeroPage, 0x04},
	{TSB, Absolute, 0x0c},

	{PHX, Implied, 0xda},
	{PLX, Implied, 0xfa},
	{PHY, Implied, 0x5a},
	{PLY, Implied, 0x7a},
}
