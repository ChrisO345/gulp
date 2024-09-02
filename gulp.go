package gulp

import (
	"fmt"
)

// TODO:
// 		Variables should probably be pointers
// 		Anything the user of this package needs to use should be in this package.
//		Everything private can be in a different package
// 		Change the way that a linear program is created, as there is definitely a better way to do this

func Gulp() bool {
	// Objective: Maximise 7 * x1 + 6 * x2
	objectiveCoefficients := []float64{7, 6}
	pairs := make([]Pair, len(objectiveCoefficients))
	for i, v := range objectiveCoefficients {
		pairs[i] = Pair{v, NewVariable(fmt.Sprintf("x%d", i+1), nil)}
	}
	objective := NewObjective(LpMaximise, pairs)

	// Constraints: 2 * x1 + 4 * x2 <= 16
	// Constraints: 3 * x1 + 2 * x2 <= 12

	constraint1Pairs := []Pair{
		{2, NewVariable("x1", nil)},
		{4, NewVariable("x2", nil)},
	}

	constraint2Pairs := []Pair{
		{3, NewVariable("x1", nil)},
		{2, NewVariable("x2", nil)},
	}

	constraint1 := NewConstraint(LpConstraintLE, constraint1Pairs, 16)
	constraint2 := NewConstraint(LpConstraintLE, constraint2Pairs, 12)
	constraints := []Constraint{constraint1, constraint2}

	lp := NewLinearProgram()
	lp.AddObjectiveFunction(objective)
	for _, c := range constraints {
		lp.AddConstraint(c)
	}
	lp.ShowSlack()
	fmt.Println(lp.String())

	lp.Solve()

	TestingMinimise()

	return Runner()
}

// TestingMinimise TODO: In order for this to work, artificial variables need to be added to the constraints
func TestingMinimise() {
	// Objective: Minimise -6 * x1 + 7 * x2 + 3 * x3
	objectiveCoefficients := []float64{-6, 7, 3}
	pairs := make([]Pair, len(objectiveCoefficients))
	for i, v := range objectiveCoefficients {
		pairs[i] = Pair{v, NewVariable(fmt.Sprintf("x%d", i+1), nil)}
	}

	objective := NewObjective(LpMinimise, pairs)

	// Constraints: 2 * x1 + 5 * x2 - x3 <= 10
	// Constraints: 1 * x1 - 1 * x2 - 4 * x3 <= -14
	// Constraints: 3 * x1 + 2 * x2 + 2 * x3 = 26

	constraint1Pairs := []Pair{
		{2, NewVariable("x1", nil)},
		{5, NewVariable("x2", nil)},
		{-1, NewVariable("x3", nil)},
	}

	constraint2Pairs := []Pair{
		{1, NewVariable("x1", nil)},
		{-1, NewVariable("x2", nil)},
		{-4, NewVariable("x3", nil)},
	}

	constraint3Pairs := []Pair{
		{3, NewVariable("x1", nil)},
		{2, NewVariable("x2", nil)},
		{2, NewVariable("x3", nil)},
	}

	constraint1 := NewConstraint(LpConstraintLE, constraint1Pairs, 10)
	constraint2 := NewConstraint(LpConstraintLE, constraint2Pairs, -14)
	constraint3 := NewConstraint(LpConstraintEQ, constraint3Pairs, 26)
	constraints := []Constraint{constraint1, constraint2, constraint3}

	lp := NewLinearProgram()
	lp.AddObjectiveFunction(objective)
	for _, c := range constraints {
		lp.AddConstraint(c)
	}

	lp.ShowSlack()
	fmt.Println(lp.String())

	// Multiply the objective function by -1 to convert it to a maximisation problem
	for i := range lp.ObjectiveFunction.Pairs {
		lp.ObjectiveFunction.Pairs[i].Coefficient *= -1
	}

	lp.ObjectiveFunction.ObjectiveType = LpMaximise

	fmt.Println(lp.String())
	lp.Solve()
}
