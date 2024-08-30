package gulp

import (
	"fmt"
	"math"
)

// TODO:
// 		Variables should probably be pointers
// 		Anything the user of this package needs to use should be in this package.
//		Everything private can be in a different package

func Gulp() bool {
	// Objective: Minimise 3x1 + 4x2
	pairs := []Pair{{3, Variable{"x1", nil}}, {-4, Variable{"x2", nil}}}

	objective := NewObjective(LpMinimise, pairs)

	// Constraints: 1 * x1 + 2 * x2 + 1
	constraintPairs := []Pair{{1, Variable{"x1", nil}}, {-2, Variable{"x2", nil}}}
	constraint1 := Constraint{LpConstraintGE, constraintPairs, 3}
	constraint2 := Constraint{LpConstraintEQ, constraintPairs, 3}
	constraint3 := Constraint{LpConstraintLE, constraintPairs, 3}
	constraints := []Constraint{constraint1, constraint2, constraint3}

	lp := NewLinearProgram(objective, constraints)
	fmt.Println(lp.String())
	return Runner()
}

// LinearProgram The Linear Program
type LinearProgram struct {
	ObjectiveFunction Objective
	Constraints       []Constraint
}

func NewLinearProgram(objective Objective, constraints []Constraint) *LinearProgram {
	return &LinearProgram{objective, constraints}
}

func (lp *LinearProgram) String() string {
	stringBuilder := fmt.Sprintf("%v\n", lp.ObjectiveFunction.String())
	for _, c := range lp.Constraints {
		stringBuilder += fmt.Sprintf("\t%v\n", c.String())
	}
	return stringBuilder
}

func (lp *LinearProgram) Solve() LpStatus {
	fmt.Println("Not Implemented")
	return LpStatusNotImplemented
}

// Objective The objective function in the linear program
type Objective struct {
	ObjectiveType LpSense
	Pairs         []Pair
}

func NewObjective(objectiveType LpSense, pairs []Pair) Objective {
	return Objective{objectiveType, pairs}
}

func (o *Objective) String() string {
	stringBuilder := ""
	if o.ObjectiveType == LpMinimise {
		stringBuilder += "Min: "
	} else {
		stringBuilder += "Max: "
	}

	for i, v := range o.Pairs {
		if i != 0 && o.Pairs[i].Coefficient >= 0 {
			stringBuilder += "+ "
		} else if o.Pairs[i].Coefficient < 0 {
			stringBuilder += "- "
		}

		stringBuilder += fmt.Sprintf("%v * %v", math.Abs(o.Pairs[i].Coefficient), v.Variable.Name)
		if i < len(o.Pairs)-1 {
			stringBuilder += " "
		}
	}
	return stringBuilder
}

// Constraint A constraint in the linear program
type Constraint struct {
	ConstraintType LpConstraint
	Pairs          []Pair
	RightHandSide  float64
}

// Print the constraint in a human-readable format
func (c *Constraint) String() string {
	stringBuilder := ""
	for i, v := range c.Pairs {
		if i != 0 && c.Pairs[i].Coefficient >= 0 {
			stringBuilder += "+ "
		} else if c.Pairs[i].Coefficient < 0 {
			stringBuilder += "- "
		}

		stringBuilder += fmt.Sprintf("%v * %v", math.Abs(c.Pairs[i].Coefficient), v.Variable.Name)
		if i < len(c.Pairs)-1 {
			stringBuilder += " "
		}
	}
	// TODO: have a mapping instead of a switch
	switch c.ConstraintType {
	case LpConstraintLE:
		stringBuilder += " <= "
	case LpConstraintEQ:
		stringBuilder += " = "
	case LpConstraintGE:
		stringBuilder += " >= "
	}
	return stringBuilder + fmt.Sprintf("%v", c.RightHandSide)
}

// AddSlackVariable Add a slack variable to the constraint
func (c *Constraint) AddSlackVariable() {
	if c.ConstraintType == LpConstraintLE {
		c.Pairs = append(c.Pairs, Pair{1, Variable{"s1", nil}})
		c.ConstraintType = LpConstraintEQ
	}
	if c.ConstraintType == LpConstraintGE {
		c.Pairs = append(c.Pairs, Pair{1, Variable{"s1", nil}})
		c.ConstraintType = LpConstraintEQ
	}
}

// Variable A variable in the linear program
type Variable struct {
	Name string
	// Value can be a float or nil, is there a better way to do this?
	Value *float64
}

func NewVariable(name string, value *float64) Variable {
	return Variable{name, value}
}

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
