package compiler

type Symbollike interface {
	Object
	isSymbollike()
}

type Symbol struct {
	Environment
}

func (*Symbol) isSymbollike() {}

func (s *Symbol) NewSubsymbol(name string) *Subsymbol {
	return &Subsymbol{
		MyName:       name,
		ParentSymbol: s,
	}
}

type Subsymbol struct {
	MyName string

	ParentSymbol *Symbol
}

func (*Subsymbol) isSymbollike() {}

func (s *Subsymbol) Name() string {
	return s.MyName
}

func (s *Subsymbol) Parent() (Object, bool) {
	return s.ParentSymbol, true
}

func (*Subsymbol) Property(string) (Object, bool) {
	return nil, false
}
