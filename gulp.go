package gulp

import (
	"fmt"
)

// TODO:
// 		Variables should probably be pointers
// 		Anything the user of this package needs to use should be in this package.
//		Everything private can be in a different package
// 		Change Linear Program struct and file to have special types for objective, constraints, etc.

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

	lp := NewLinearProgram(objective, constraints)
	lp.ShowSlack()
	fmt.Println(lp.String())

	lp.Solve()
	return Runner()
}
