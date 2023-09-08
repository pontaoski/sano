package compiler

import "strings"

type Environment struct {
	Bindings  map[string]Object
	MyName    string
	ParentEnv *Environment
}

func NewRootEnvironment() *Environment {
	return &Environment{
		Bindings:  map[string]Object{},
		MyName:    "world",
		ParentEnv: nil,
	}
}

func (e *Environment) NewSymbol(name string) *Symbol {
	return &Symbol{
		Environment: Environment{
			Bindings:  map[string]Object{},
			MyName:    name,
			ParentEnv: e,
		},
	}
}

func (e *Environment) Bind(s string, o Object) bool {
	if e.Contains(s) {
		return false
	}
	e.Bindings[s] = o
	return true
}

func (e *Environment) Contains(s string) bool {
	if _, ok := e.Property(s); ok {
		return true
	}
	if e.ParentEnv == nil {
		return false
	}
	return e.ParentEnv.Contains(s)
}

func (e *Environment) Lookup(s string) (Object, bool) {
	if v, ok := e.Property(s); ok {
		return v, true
	}
	if e.ParentEnv == nil {
		return nil, false
	}
	return e.ParentEnv.Lookup(s)
}

func (e *Environment) Property(s string) (Object, bool) {
	if v, ok := e.Bindings[s]; ok {
		return v, true
	}
	return nil, false
}

func (e *Environment) Parent() (Object, bool) {
	return e.ParentEnv, e.ParentEnv == nil
}

func (e *Environment) Name() string {
	return e.MyName
}

type Object interface {
	Name() string
	Parent() (Object, bool)
	Property(s string) (Object, bool)
}

func GlobalName(o Object) string {
	var s []string
	for o, ok := o, true; ok; o, ok = o.Parent() {
		s = append(s, o.Name())
	}
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return strings.Join(s, "/")
}
