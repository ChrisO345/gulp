package gulp

import (
	"fmt"
	"math"
)

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
