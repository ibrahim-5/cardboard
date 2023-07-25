package object

// Environment
type Environment struct {
	store map[string]Object
	outer *Environment
}

func CreateEnvironment() *Environment {
	m := make(map[string]Object)
	return &Environment{store: m, outer: nil}
}

func CreateEnclosedEnvironment(outer *Environment) *Environment {
	env := CreateEnvironment()
	env.outer = outer
	return env
}

func (env *Environment) Get(key string) (Object, bool) {
	obj, found := env.store[key]
	if !found && env.outer != nil {
		obj, found = env.outer.Get(key)
	}
	return obj, found
}

func (env *Environment) Set(key string, val Object) Object {
	env.store[key] = val
	return val
}
