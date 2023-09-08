package linker

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"google.golang.org/protobuf/encoding/protojson"
)

// must not have duplicate symbols
// behaviour undefined if there are
func Concatenate(os []*Object) *Object {
	total := &Object{Fragments: map[string]*Fragment{}}
	for _, o := range os {
		for key, frag := range o.Fragments {
			total.Fragments[key] = frag
		}
	}
	return total
}

func WriteUint8(w io.Writer, u uint8) (int, error) {
	b := [1]byte{u}
	return w.Write(b[:])
}

func WriteUint16(w io.Writer, u uint16) (int, error) {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], u)
	return w.Write(b[:])
}

// returns a prg file
func LinkToPrg(o []*Object) ([]byte, error) {
	bigly := Concatenate(o)

	main, ok := bigly.Fragments["main"]
	if !ok {
		return nil, errors.New("youre missing a main fragment")
	}

	println(protojson.MarshalOptions{Indent: "  "}.Format(main))

	var b bytes.Buffer

	WriteUint16(&b, 0x0801) // memory location to load into
	WriteUint16(&b, 0x080C) // pointer to line of basic code
	WriteUint16(&b, 0x000A) // line number
	WriteUint8(&b, 0x9E)    // sys token
	WriteUint8(&b, 0x32)    // "2"
	WriteUint8(&b, 0x30)    // "0"
	WriteUint8(&b, 0x36)    // "6"
	WriteUint8(&b, 0x31)    // "1"
	WriteUint8(&b, 0x00)    // nul, line terminator
	WriteUint16(&b, 0x0000) // pointer to line of basic code (0x0000 == end of program)

	// we start execution at address 0x0810
	for _, expr := range main.Expressions {
		switch t := expr.Inner.(type) {
		case *Expression_Literal_:
			b.Write(t.Literal.Value)
		case *Expression_Symbol_:
			panic("not implemented")
		case *Expression_Unary_:
			panic("not implemented")
		case *Expression_Subsymbol_:
			panic("not implemented")
		default:
			panic("unhandled case")
		}
	}

	return b.Bytes(), nil
}
