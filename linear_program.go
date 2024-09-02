package gulp

import "fmt"

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

func (lp *LinearProgram) ShowSlack() {
	for i := range lp.Constraints {
		lp.Constraints[i].ShowSlackVariable = true

		// Set Slacks to unique names
		for j := range lp.Constraints[i].Pairs {
			if lp.Constraints[i].Pairs[j].Variable.IsSlack {
				lp.Constraints[i].Pairs[j].Variable.Name = fmt.Sprintf("s%d", i+1)
			}
		}
	}
}

func (lp *LinearProgram) Solve() LpStatus {
	lp.BasisVariable()
	return LpStatusNotImplemented
}

func (lp *LinearProgram) BasisVariable() {
	stringBuilder := ""
	for _, v := range lp.Constraints {
		// One Slack Variable per Constraint
		for _, p := range v.Pairs {
			if p.Variable.IsSlack {
				stringBuilder += fmt.Sprintf(p.Variable.Name)
			}
		}
		stringBuilder += ", "
	}
	stringBuilder = stringBuilder[:len(stringBuilder)-2] + " >= 0"
	fmt.Println(stringBuilder)
}
