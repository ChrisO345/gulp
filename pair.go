package gulp

import "fmt"

// Pair A pair of a coefficient and a variable
//
//	TODO: Rename this because it's not clear what it is
type Pair struct {
	Coefficient float64
	Variable    Variable
}

func NewPair(coefficient float64, variable Variable) Pair {
	return Pair{coefficient, variable}
}

func (pair *Pair) String() string {
	return fmt.Sprintf("%v * %v", pair.Coefficient, pair.Variable.Name)
}

func (pair *Pair) Result() float64 {
	return pair.Coefficient * *pair.Variable.Value
}
