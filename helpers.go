package gulp

import (
	"fmt"
	"math"
)

func (lp *LinearProgram) String() string {
	stringBuilder := ""
	if lp.Sense == LpMinimise {
		stringBuilder += "Min: "
	} else {
		stringBuilder += "Max: "
	}

	for i, v := range lp.ObjectiveFunction.Terms {
		if i != 0 && lp.ObjectiveFunction.Terms[i].Coefficient >= 0 {
			stringBuilder += "+ "
		} else if lp.ObjectiveFunction.Terms[i].Coefficient < 0 {
			stringBuilder += "- "
		}

		stringBuilder += fmt.Sprintf("%v * %v", math.Abs(lp.ObjectiveFunction.Terms[i].Coefficient), v.Variable.Name)
		if i < len(lp.ObjectiveFunction.Terms)-1 {
			stringBuilder += " "
		}
	}

	for _, c := range lp.Constraints {
		stringBuilder += "\n\t"
		for i, v := range c.Terms {
			if i != 0 && c.Terms[i].Coefficient >= 0 {
				stringBuilder += "+ "
			} else if c.Terms[i].Coefficient < 0 {
				stringBuilder += "- "
			}

			stringBuilder += fmt.Sprintf("%v * %v", math.Abs(c.Terms[i].Coefficient), v.Variable.Name)
			if i < len(c.Terms)-1 {
				stringBuilder += " "
			}
		}
		switch c.ConstraintType {
		case LpConstraintLE:
			stringBuilder += " <= "
		case LpConstraintEQ:
			stringBuilder += " = "
		case LpConstraintGE:
			stringBuilder += " >= "
		}

		stringBuilder += fmt.Sprintf("%v", c.RightHandSide)
	}

	return stringBuilder
}
