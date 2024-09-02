package gulp

import "fmt"

// LinearProgram The Linear Program
type LinearProgram struct {
	ObjectiveFunction Objective
	Constraints       []Constraint
	ShowSlackVariable bool

	// Solution
	Solution     []float64
	OptimalValue float64
}

func NewLinearProgram() LinearProgram {
	return LinearProgram{}
}

func (lp *LinearProgram) AddObjectiveFunction(objective Objective) {
	lp.ObjectiveFunction = objective
}

func (lp *LinearProgram) AddConstraint(constraint Constraint) {
	if lp.ShowSlackVariable {
		constraint.ShowSlackVariable = true
		constraint.Pairs[len(constraint.Pairs)-1].Variable.Name = fmt.Sprintf("s%d", len(lp.Constraints)+1)
	}
	lp.Constraints = append(lp.Constraints, constraint)
}

func (lp *LinearProgram) String() string {
	stringBuilder := fmt.Sprintf("%v\n", lp.ObjectiveFunction.String())
	for _, c := range lp.Constraints {
		stringBuilder += fmt.Sprintf("\t%v\n", c.String())
	}
	return stringBuilder
}

func (lp *LinearProgram) ShowSlack() {
	lp.ShowSlackVariable = true
	var slacks []Variable
	for i := range lp.Constraints {
		lp.Constraints[i].ShowSlackVariable = true

		// Set Slacks to unique names
		for j := range lp.Constraints[i].Pairs {
			if lp.Constraints[i].Pairs[j].Variable.IsSlack {
				lp.Constraints[i].Pairs[j].Variable.Name = fmt.Sprintf("s%d", i+1)
				slacks = append(slacks, lp.Constraints[i].Pairs[j].Variable)
			}
		}
	}
	// Add to objective function with 0 coefficient
	for _, v := range slacks {
		lp.ObjectiveFunction.Pairs = append(lp.ObjectiveFunction.Pairs, Pair{0, v})
	}
}

func (lp *LinearProgram) Solve() LpStatus {
	if !lp.ShowSlackVariable {
		lp.ShowSlack()
	}
	tableau := NewTableau(lp)
	fmt.Println(tableau.String())

	// FIXME: this cannot be specified this way as the solution could be unbounded
	for !tableau.IsOptimal() {
		tableau.Pivot()
	}
	fmt.Println(tableau.TableauValue)
	return LpStatusOptimal
}
