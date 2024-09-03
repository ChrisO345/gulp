package gulp

import (
	"fmt"
	"math"
)

// LinearProgram The Linear Program
type LinearProgram struct {
	ObjectiveFunction LpExpression
	Sense             LpSense
	Constraints       []_constraint
	hiddenSense       LpSense

	// Solution
	Solution     map[string]float64
	OptimalValue float64
	Status       LpStatus
}

// NewLinearProgram Create a new Linear Program
func NewLinearProgram() LinearProgram {
	lp := LinearProgram{}
	lp.Status = LpStatusNotSolved
	return lp
}

// AddObjective Add an objective to the linear program
func (lp *LinearProgram) AddObjective(sense LpSense, objective LpExpression) *LinearProgram {
	lp.hiddenSense = sense
	if sense == LpMinimise {
		for i := range objective.Terms {
			objective.Terms[i].Coefficient *= -1
		}
		sense = LpMaximise
	}
	lp.Sense = sense
	lp.ObjectiveFunction = objective
	return lp
}

// AddConstraint Add a constraint to the linear program
func (lp *LinearProgram) AddConstraint(constraint LpExpression, constraintType LpConstraintType, rightHandSide float64) *LinearProgram {
	// Panic if objective function is not set
	if len(lp.ObjectiveFunction.Terms) == 0 {
		panic("Objective function not set")
	}

	if rightHandSide < 0 {
		// Multiply the constraint by -1, flip equality sign
		rightHandSide = math.Abs(rightHandSide)
		for i := range constraint.Terms {
			constraint.Terms[i].Coefficient *= -1
		}
		constraintType = -constraintType
	}

	// Add Artificial Variables
	if constraintType == LpConstraintEQ || constraintType == LpConstraintGE {
		variable := NewArtificialVariable(fmt.Sprintf("a%d", len(lp.Constraints)+1))
		constraint.Terms = append(constraint.Terms, NewTerm(1, variable))
		lp.ObjectiveFunction.Terms = append(lp.ObjectiveFunction.Terms, NewTerm(-1e20, variable))
	}

	// Add Slack Variables
	if constraintType == LpConstraintLE || constraintType == LpConstraintGE {
		variable := NewSlackVariable(fmt.Sprintf("s%d", len(lp.Constraints)+1))
		sign := 1.0
		if constraintType == LpConstraintGE {
			sign = -1.0
		}
		constraint.Terms = append(constraint.Terms, NewTerm(sign, variable))
		lp.ObjectiveFunction.Terms = append(lp.ObjectiveFunction.Terms, NewTerm(0, variable))
		constraintType = LpConstraintEQ
	}

	lp.Constraints = append(lp.Constraints, _constraint{constraintType, constraint.Terms, rightHandSide})

	return lp
}

func (lp *LinearProgram) Solve() *LinearProgram {
	tableau := NewTableau(lp)
	for !tableau.IsOptimal() {
		tableau.Pivot()
	}
	lp.OptimalValue = tableau.TableauValue * float64(lp.hiddenSense)
	lp.Status = LpStatusOptimal
	solution := tableau.GetSolution()
	lp.Solution = make(map[string]float64)
	for _, v := range lp.ObjectiveFunction.Terms {
		if v.Variable.IsSlack || v.Variable.IsArtificial {
			continue
		}
		lp.Solution[v.Variable.Name] = solution[v.Variable.Name]
	}

	return lp
}

func (lp *LinearProgram) PrintSolution() {
	fmt.Println(lp.Status.String())
	fmt.Println(lp.OptimalValue)
	for k, v := range lp.Solution {
		fmt.Printf("%v: %v\n", k, v)
	}
}

/* #####################################################################################################################
TO BE MOVED TO SEPARATE FILES, SOME OF THE STRUCTS ARE PRIVATE
##################################################################################################################### */

type LpExpression struct {
	Terms []LpTerm
}

func NewExpression(terms []LpTerm) LpExpression {
	return LpExpression{terms}
}

type LpTerm struct {
	Coefficient float64
	Variable    LpVariable // These get added to the variable list in the LinearProgram??
}

func NewTerm(coefficient float64, variable LpVariable) LpTerm {
	return LpTerm{coefficient, variable}
}

type LpVariable struct {
	Name         string
	Value        float64
	IsSlack      bool
	IsArtificial bool
}

func NewVariable(name string) LpVariable {
	return LpVariable{name, 0, false, false}
}

func NewSlackVariable(name string) LpVariable {
	return LpVariable{name, 0, true, false}
}

func NewArtificialVariable(name string) LpVariable {
	return LpVariable{name, 0, false, true}
}

type _constraint struct {
	ConstraintType LpConstraintType
	Terms          []LpTerm
	RightHandSide  float64
}

/* #####################################################################################################################
TO BE MOVED TO SEPARATE FILES
##################################################################################################################### */

// LpCategory Note that this is currently not used
type LpCategory string

const (
	LpContinuous = LpCategory("Continuous")
	LpInteger    = LpCategory("Integer")
	LpBinary     = LpCategory("Binary")
)

// TODO: Category Mapping ?? dictionary??

type LpSense int

const (
	LpMinimise = LpSense(-1)
	LpMaximise = LpSense(1)
)

// TODO: Sense Mapping ?? dictionary??

type LpStatus int

const (
	LpStatusNotSolved      = LpStatus(0)
	LpStatusOptimal        = LpStatus(1)
	LpStatusInfeasible     = LpStatus(2)
	LpStatusUnbounded      = LpStatus(3)
	LpStatusUndefined      = LpStatus(4)
	LpStatusNotImplemented = LpStatus(5)
)

var LpStatusMap = map[LpStatus]string{
	LpStatusNotSolved:      "Not Solved",
	LpStatusOptimal:        "Optimal",
	LpStatusInfeasible:     "Infeasible",
	LpStatusUnbounded:      "Unbounded",
	LpStatusUndefined:      "Undefined",
	LpStatusNotImplemented: "Not Implemented",
}

func (s *LpStatus) String() string {
	return LpStatusMap[*s]
}

type LpConstraintType int

const (
	LpConstraintLE = LpConstraintType(-1)
	LpConstraintEQ = LpConstraintType(0)
	LpConstraintGE = LpConstraintType(1)
)

// TODO: Constraint Mapping ?? dictionary??
