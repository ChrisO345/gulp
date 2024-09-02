package gulp

import (
	"fmt"
	"testing"
)

/***********************************************************************************************************************
CONSTRAINTS
***********************************************************************************************************************/

func TestConstraint(t *testing.T) {
	expected := "1 * x1 + 2 * x2 >= 3"

	pairs := []Pair{{1, NewVariable("x1", nil)}, {2, NewVariable("x2", nil)}}
	constraint := NewConstraint(LpConstraintGE, pairs, 3)

	result := constraint.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestConstraint_AddSlackVariable_EQ(t *testing.T) {
	expected := "1 * x1 + 2 * x2 = 3"

	pairs := []Pair{{1, NewVariable("x1", nil)}, {2, NewVariable("x2", nil)}}
	constraint := NewConstraint(LpConstraintEQ, pairs, 3)
	constraint.AddSlackVariable()

	result := constraint.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestConstraint_AddSlackVariable_GE(t *testing.T) {
	expected := "1 * x1 + 2 * x2 - 1 * s1 = 3"

	pairs := []Pair{{1, NewVariable("x1", nil)}, {2, NewVariable("x2", nil)}}
	constraint := NewConstraint(LpConstraintGE, pairs, 3)
	constraint.AddSlackVariable()
	constraint.ShowSlack()

	result := constraint.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestConstraint_AddSlackVariable_LE(t *testing.T) {
	expected := "- 1 * x1 + 2 * x2 + 1 * s1 = 3"

	pairs := []Pair{{-1, NewVariable("x1", nil)}, {2, NewVariable("x2", nil)}}
	constraint := NewConstraint(LpConstraintLE, pairs, 3)
	constraint.AddSlackVariable()
	constraint.ShowSlack()

	result := constraint.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

/***********************************************************************************************************************
OBJECTIVES
***********************************************************************************************************************/

func TestNewObjective(t *testing.T) {
	expected := "Min: 3 * x1 + 4 * x2"

	pairs := []Pair{{3, NewVariable("x1", nil)}, {4, NewVariable("x2", nil)}}
	objective := NewObjective(LpMinimise, pairs)

	result := objective.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

/***********************************************************************************************************************
LINEAR PROGRAM
***********************************************************************************************************************/

func TestNewLinearProgram(t *testing.T) {
	expected := "Min: - 3 * x1 + 4 * x2\n\t1 * x1 - 2 * x2 >= 3\n\t1 * x1 - 2 * x2 = 3\n\t1 * x1 - 2 * x2 <= 3\n"

	objective := Objective{LpMinimise, []Pair{{-3, NewVariable("x1", nil)}, {4, NewVariable("x2", nil)}}}
	pairs := []Pair{{1, NewVariable("x1", nil)}, {-2, NewVariable("x2", nil)}}
	constraint1 := NewConstraint(LpConstraintGE, pairs, 3)
	constraint2 := NewConstraint(LpConstraintEQ, pairs, 3)
	constraint3 := NewConstraint(LpConstraintLE, pairs, 3)
	constraints := []Constraint{constraint1, constraint2, constraint3}
	lp := NewLinearProgram(objective, constraints)

	result := lp.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

/***********************************************************************************************************************
PAIR
***********************************************************************************************************************/

func TestNewPair(t *testing.T) {
	expected := "3 * x1"

	pair := NewPair(3, NewVariable("x1", nil))

	result := pair.String()
	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}

func TestPair_Result(t *testing.T) {
	expected := 9.0

	variableValue := 3.0
	pair := NewPair(3, NewVariable("x1", &variableValue))

	result := pair.Result()
	if result != expected {
		t.Errorf("Expected %f but got %f", expected, result)
	}
}

/***********************************************************************************************************************
VARIABLE
***********************************************************************************************************************/

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

/***********************************************************************************************************************
***********************************************************************************************************************/

func TestGulpRun(t *testing.T) {
	result := Gulp()
	fmt.Println(result)
}
