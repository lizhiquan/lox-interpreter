package lox

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]any)}
}

func NewEnvironmentWithEnclosing(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) (any, error) {
	if val, ok := e.values[name.Lexeme]; ok {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}

func (e *Environment) Assign(name Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}

	return NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}
