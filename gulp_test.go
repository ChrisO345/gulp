package gulp

import (
	"fmt"
	"testing"
)

/***********************************************************************************************************************
***********************************************************************************************************************/

func TestConstraint_AddSlackVariableEQ(t *testing.T) {
	expected := "1 * x1 + 2 * x2 = 3"

	pairs := []Pair{{1, Variable{"x1", nil}}, {2, Variable{"x2", nil}}}
	constraint := Constraint{LpConstraintEQ, pairs, 3}
	constraint.AddSlackVariable()

	result := constraint.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestConstraint_AddSlackVariableGE(t *testing.T) {
	expected := "1 * x1 + 2 * x2 + 1 * s1 = 3"

	pairs := []Pair{{1, Variable{"x1", nil}}, {2, Variable{"x2", nil}}}
	constraint := Constraint{LpConstraintGE, pairs, 3}
	constraint.AddSlackVariable()

	result := constraint.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestConstraint_AddSlackVariableLE(t *testing.T) {
	expected := "1 * x1 + 2 * x2 + 1 * s1 = 3"

	pairs := []Pair{{1, Variable{"x1", nil}}, {2, Variable{"x2", nil}}}
	constraint := Constraint{LpConstraintLE, pairs, 3}
	constraint.AddSlackVariable()

	result := constraint.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

/***********************************************************************************************************************
***********************************************************************************************************************/

func TestNewObjective(t *testing.T) {
	expected := "Min: 3 * x1 + 4 * x2"

	pairs := []Pair{{3, Variable{"x1", nil}}, {4, Variable{"x2", nil}}}
	objective := NewObjective(LpMinimise, pairs)

	result := objective.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestNewLinearProgram(t *testing.T) {
	expected := "Min: 3 * x1 + 4 * x2\n\t1 * x1 + 2 * x2 >= 3\n\t1 * x1 + 2 * x2 = 3\n\t1 * x1 + 2 * x2 <= 3\n"

	objective := Objective{LpMinimise, []Pair{{3, Variable{"x1", nil}}, {4, Variable{"x2", nil}}}}
	pairs := []Pair{{1, Variable{"x1", nil}}, {2, Variable{"x2", nil}}}
	constraint1 := Constraint{LpConstraintGE, pairs, 3}
	constraint2 := Constraint{LpConstraintEQ, pairs, 3}
	constraint3 := Constraint{LpConstraintLE, pairs, 3}
	constraints := []Constraint{constraint1, constraint2, constraint3}
	lp := NewLinearProgram(objective, constraints)

	result := lp.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

/***********************************************************************************************************************
***********************************************************************************************************************/

func TestNewPair(t *testing.T) {
	expected := "3 * x1"

	pair := NewPair(3, Variable{"x1", nil})

	result := pair.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestNewVariable(t *testing.T) {
	expectedName := "x1"
	expectedValue := 3.0

	value := 3.0
	variable := NewVariable("x1", &value)

	resultName := variable.Name
	resultValue := *variable.Value
	if resultName != expectedName || resultValue != expectedValue {
		t.Errorf("Expected %s %f but got %s %f", expectedName, expectedValue, resultName, resultValue)
	}
}

/***********************************************************************************************************************
***********************************************************************************************************************/

func TestGulpRun(t *testing.T) {
	result := Gulp()
	fmt.Println(result)
}
