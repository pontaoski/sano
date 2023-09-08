package compiler

import (
	"Sano/cpu"
	"Sano/linker"
	"Sano/parser"
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

type Compiler struct {
	Instructions cpu.OpcodeSet
}

type CompilationError struct {
	Message  string
	Location lexer.Position
}

func (e *CompilationError) String() string {
	return fmt.Sprintf("%s: %s", e.Location, e.Message)
}

func (c *Compiler) Compile(f *parser.File) (*linker.Object, []CompilationError) {
	errors := []CompilationError{}
	fragments := map[string]*linker.Fragment{}
	env := NewRootEnvironment()

	for _, it := range f.Fragment {
		fragmentEnv := env.NewSymbol(it.Name)
		if !env.Bind(it.Name, fragmentEnv) {
			errors = append(errors, CompilationError{fmt.Sprintf("Duplicate symbol '%s'", it.Name), it.Pos})
		}
		for _, s := range it.Statements {
			switch s := s.(type) {
			case parser.OpcodeInvocation:
			case parser.SymbolDeclaration:
				if !fragmentEnv.Bind(s.Name, fragmentEnv.NewSubsymbol(s.Name)) {
					errors = append(errors, CompilationError{fmt.Sprintf("Duplicate symbol '%s'", s.Name), it.Pos})
				}
			default:
				panic("unhandled case")
			}
		}
	}

	for _, it := range f.Fragment {
		expressions := []*linker.Expression{}
		fragmentEnvObj, _ := env.Lookup(it.Name)
		fragmentEnv := fragmentEnvObj.(*Symbol)

		for _, s := range it.Statements {
			switch s := s.(type) {
			case parser.OpcodeInvocation:
				opcode, ok := cpu.OpcodeNames[strings.ToLower(s.Opcode)]
				if !ok {
					errors = append(errors, CompilationError{fmt.Sprintf("Invalid opcode '%s'", s.Opcode), s.Pos})
					continue
				}
				opcodes := c.Instructions.Find(opcode)
				if len(opcodes) == 0 {
					errors = append(errors, CompilationError{fmt.Sprintf("Opcode '%s' does not exist on the architecture", s.Opcode), s.Pos})
					continue
				}
				resolved, ok := opcodes.FindOne(opcode, s.Address.AddressingMode())
				if !ok {
					errors = append(errors, CompilationError{fmt.Sprintf("Opcode '%s' cannot be used with %s addressing", s.Opcode, s.Address.AddressingMode()), s.Pos})
				}
				_ = resolved

				expressions = append(expressions, &linker.Expression{
					Inner: &linker.Expression_Literal_{
						Literal: &linker.Expression_Literal{
							Value: []byte{resolved.Hex},
						},
					},
				})

				var expr parser.Expression
				var size linker.SymbolSize

				switch addr := s.Address.(type) {
				case parser.AddrImmediate:
					expr = addr.Value
					size = linker.SymbolSize_BYTE
				case parser.AddrAbsoluteIndirect:
					expr = addr.Address
					size = linker.SymbolSize_WORD
				case parser.AddrAbsoluteIndexedX:
					expr = addr.Address
					size = linker.SymbolSize_WORD
				case parser.AddrAbsoluteIndexedY:
					expr = addr.Address
					size = linker.SymbolSize_WORD
				case parser.AddrAbsoluteAddress:
					expr = addr.Address
					size = linker.SymbolSize_WORD
				case parser.AddrZeroPageXIndexedIndirect:
					expr = addr.Address
					size = linker.SymbolSize_BYTE
				case parser.AddrZeroPageIndirectYIndex:
					expr = addr.Address
					size = linker.SymbolSize_BYTE
				case parser.AddrZeroPageXIndexed:
					expr = addr.Address
					size = linker.SymbolSize_BYTE
				case parser.AddrZeroPageYIndexed:
					expr = addr.Address
					size = linker.SymbolSize_BYTE
				case parser.AddrZeroPage:
					expr = addr.Address
					size = linker.SymbolSize_BYTE
				case parser.AddrRelative:
					expr = addr.Address
					size = linker.SymbolSize_RELATIVE
				case parser.AddrAccumulator:
					expr = nil
				case parser.AddrImplied:
					expr = nil
				}

				if expr == nil {
					continue
				}

				// fmt.Printf("%#v\n", expr)
				switch e := expr.(type) {
				case parser.NumericLiteral:
					// TODO: check that it fits into the size
					// errors = append(errors, CompilationError{"Unhandled numeric literal", e.Pos})
					var numericBytes []byte
					switch size {
					case linker.SymbolSize_WORD:
						numericBytes = []byte{byte(e.Number), byte(e.Number >> 8)}
					case linker.SymbolSize_BYTE, linker.SymbolSize_RELATIVE:
						numericBytes = []byte{byte(e.Number)}
					}
					expressions = append(expressions, &linker.Expression{
						Inner: &linker.Expression_Literal_{
							Literal: &linker.Expression_Literal{
								Value: numericBytes,
							},
						},
					})
				case parser.Symbol:
					o, ok := fragmentEnv.Lookup(e.Name)
					if !ok {
						errors = append(errors, CompilationError{fmt.Sprintf("Symbol not found: '%s'", e.Name), e.Pos})
						continue
					}
					if _, ok := o.(Symbollike); !ok {
						errors = append(errors, CompilationError{fmt.Sprintf("'%s' is not a symbol", e.Name), e.Pos})
						continue
					}
					symbol, _ := o.(Symbollike)
					expressions = append(expressions, &linker.Expression{
						Inner: &linker.Expression_Symbol_{
							Symbol: &linker.Expression_Symbol{
								Name: GlobalName(symbol),
								Size: size,
							},
						},
					})
				}
			case parser.SymbolDeclaration:
				sym, _ := fragmentEnv.Lookup(s.Name)
				subsymbol := sym.(*Subsymbol)
				expressions = append(expressions, &linker.Expression{
					Inner: &linker.Expression_Subsymbol_{
						Subsymbol: &linker.Expression_Subsymbol{
							Name: GlobalName(subsymbol),
						},
					},
				})
			default:
				panic("unhandled case")
			}
		}

		fragments[it.Name] = &linker.Fragment{
			Expressions: expressions,
		}
	}

	total := &linker.Object{}
	total.Fragments = fragments
	return total, errors
}
