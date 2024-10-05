package object

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outerEnv = outer
	return env
}

type Environment struct {
	store    map[string]Object
	outerEnv *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outerEnv != nil {
		obj, ok = e.outerEnv.Get(name)
	}
	return obj, ok
}
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
