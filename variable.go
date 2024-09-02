package gulp

// Variable A variable in the linear program
type Variable struct {
	Name    string
	Value   *float64 // Value can be a float or nil, TODO: is there a better way to do this?
	IsSlack bool
}

func NewVariable(name string, value *float64) Variable {
	return Variable{name, value, false}
}

func NewSlackVariable(name string, value *float64) Variable {
	return Variable{name, value, true}
}
