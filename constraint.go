package gulp

import (
	"fmt"
	"math"
)

// Constraint A constraint in the linear program
type Constraint struct {
	ConstraintType    LpConstraint
	Pairs             []Pair
	RightHandSide     float64
	ShowSlackVariable bool
}

func NewConstraint(constraintType LpConstraint, pairs []Pair, rightHandSide float64) Constraint {
	constraint := Constraint{constraintType, pairs, rightHandSide, false}
	constraint.AddSlackVariable()
	return constraint
}

// Print the constraint in a human-readable format, with an optional variable showSlack
func (c *Constraint) String() string {
	stringBuilder := ""
	for i, v := range c.Pairs {
		if v.Variable.IsSlack && !c.ShowSlackVariable {
			// Slack should always be the last variable
			stringBuilder += ""
			switch v.Coefficient > 0 {
			case true:
				stringBuilder += "<= "
			case false:
				stringBuilder += ">= "
			}
			return stringBuilder + fmt.Sprintf("%v", c.RightHandSide)
		}

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
		c.Pairs = append(c.Pairs, Pair{1, NewSlackVariable("s1", nil)})
		c.ConstraintType = LpConstraintEQ
	}
	if c.ConstraintType == LpConstraintGE {
		c.Pairs = append(c.Pairs, Pair{-1, NewSlackVariable("s1", nil)})
		c.ConstraintType = LpConstraintEQ
	}
}

func (c *Constraint) ShowSlack() {
	c.ShowSlackVariable = true
}
