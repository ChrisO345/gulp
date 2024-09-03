package gulp

import (
	"fmt"
	"math"
	"testing"
)

/* *********************************************************************************************************************
Variables
********************************************************************************************************************* */

func TestNewVariable(t *testing.T) {
	expected := LpVariable{
		Name:         "Apples",
		Value:        0,
		IsSlack:      false,
		IsArtificial: false,
	}

	result := NewVariable("Apples")
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

/* *********************************************************************************************************************
Terms
********************************************************************************************************************* */

func TestNewTerm(t *testing.T) {
	expected := LpTerm{
		Coefficient: 7,
		Variable:    LpVariable{Name: "Apples", Value: 0, IsSlack: false, IsArtificial: false},
	}

	result := NewTerm(7, NewVariable("Apples"))
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

/* *********************************************************************************************************************
Expressions
********************************************************************************************************************* */

func TestNewExpression(t *testing.T) {
	expected := LpExpression{
		Terms: []LpTerm{
			NewTerm(7, NewVariable("Apples")),
			NewTerm(6, NewVariable("Bananas")),
		},
	}

	result := NewExpression([]LpTerm{
		NewTerm(7, NewVariable("Apples")),
		NewTerm(6, NewVariable("Bananas")),
	})

	for i, v := range result.Terms {
		if v != expected.Terms[i] {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	}
}

/* *********************************************************************************************************************
LinearProgram
********************************************************************************************************************* */

func TestNewLinearProgram(t *testing.T) {
	expected := LinearProgram{
		Status: LpStatusNotSolved,
	}

	lp := NewLinearProgram()
	if lp.Status != expected.Status {
		t.Errorf("Expected %v, got %v", expected, lp)
	}
}

func TestAddObjective(t *testing.T) {
	lp := NewLinearProgram()
	objective := NewExpression([]LpTerm{
		NewTerm(7, NewVariable("Apples")),
		NewTerm(6, NewVariable("Bananas")),
	})

	expectedLP := LinearProgram{
		ObjectiveFunction: objective,
	}
	expectedSense := LpMaximise

	lp.AddObjective(LpMaximise, objective)
	for i, v := range lp.ObjectiveFunction.Terms {
		if v != expectedLP.ObjectiveFunction.Terms[i] {
			t.Errorf("Expected %v, got %v", expectedLP, lp)
		}
	}
	if lp.Sense != expectedSense {
		t.Errorf("Expected %v, got %v", expectedSense, lp.Sense)
	}
}

func TestAddConstraint(t *testing.T) {
	fmt.Println("TestAddConstraint")
}

func TestAddConstraintNoObjective(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, got nil")
		}
	}()

	lp := NewLinearProgram()
	constraint := NewExpression([]LpTerm{
		NewTerm(2, NewVariable("Apples")),
		NewTerm(4, NewVariable("Bananas")),
	})

	lp.AddConstraint(constraint, LpConstraintLE, 16)
}

func TestSolve(t *testing.T) {
	expectedSense := LpMaximise
	expectedStatus := LpStatusOptimal
	expectedOptimalValue := 32.0
	expectedSolution := map[string]float64{
		"Apples":  2,
		"Bananas": 3,
	}
	// Objective: Maximise 7 * x1 + 6 * x2
	// Constraints: 2 * x1 + 4 * x2 <= 16
	// Constraints: 3 * x1 + 2 * x2 <= 12
	variables := []LpVariable{
		NewVariable("Apples"),
		NewVariable("Bananas"),
	}

	terms := []LpTerm{
		NewTerm(7, variables[0]),
		NewTerm(6, variables[1]),
	}
	objective := NewExpression(terms)

	terms2 := []LpTerm{
		NewTerm(2, variables[0]),
		NewTerm(4, variables[1]),
	}

	terms3 := []LpTerm{
		NewTerm(3, variables[0]),
		NewTerm(2, variables[1]),
	}

	lp := NewLinearProgram()
	lp.AddObjective(LpMaximise, objective).
		AddConstraint(NewExpression(terms2), LpConstraintLE, 16).
		AddConstraint(NewExpression(terms3), LpConstraintLE, 12).
		Solve()

	if lp.hiddenSense != expectedSense {
		t.Errorf("Expected %v, got %v", expectedSense, lp.hiddenSense)
	}
	if lp.Status != expectedStatus {
		t.Errorf("Expected %v, got %v", expectedStatus, lp.Status)
	}
	if lp.OptimalValue != expectedOptimalValue {
		t.Errorf("Expected %v, got %v", expectedOptimalValue, lp.OptimalValue)
	}
	for k, v := range lp.Solution {
		if math.Abs(v-expectedSolution[k]) > 0.0001 {
			t.Errorf("Expected %v, got %v", expectedSolution, lp.Solution)
		}
	}
}

func TestSolveMinimum(t *testing.T) {
	expectedSense := LpMinimise
	expectedStatus := LpStatusOptimal
	expectedOptimalValue := 16.0
	expectedSolution := map[string]float64{
		"x1": 3,
		"x2": 0,
		"x3": 8.5,
	}
	// Objective: Minimise - 6 * x1 + 7 * x2 + 4 * x3
	// Constraints: 2 * x1 + 5 * x2 - 1 * x3 <= 18
	// Constraints: 1 * x1 - 1 * x2 - 2 * x3 <= -14
	// Constraints: 3 * x1 + 2 * x2 + 2 * x3 = 26

	variables := []LpVariable{
		NewVariable("x1"),
		NewVariable("x2"),
		NewVariable("x3"),
	}

	terms := []LpTerm{
		NewTerm(-6, variables[0]),
		NewTerm(7, variables[1]),
		NewTerm(4, variables[2]),
	}
	objective := NewExpression(terms)

	terms2 := []LpTerm{
		NewTerm(2, variables[0]),
		NewTerm(5, variables[1]),
		NewTerm(-1, variables[2]),
	}

	terms3 := []LpTerm{
		NewTerm(1, variables[0]),
		NewTerm(-1, variables[1]),
		NewTerm(-2, variables[2]),
	}

	terms4 := []LpTerm{
		NewTerm(3, variables[0]),
		NewTerm(2, variables[1]),
		NewTerm(2, variables[2]),
	}

	lp := NewLinearProgram()
	lp.AddObjective(LpMinimise, objective).
		AddConstraint(NewExpression(terms2), LpConstraintLE, 18).
		AddConstraint(NewExpression(terms3), LpConstraintLE, -14).
		AddConstraint(NewExpression(terms4), LpConstraintEQ, 26).
		Solve()

	if lp.hiddenSense != expectedSense {
		t.Errorf("Expected %v, got %v", expectedSense, lp.hiddenSense)
	}

	if lp.Status != expectedStatus {
		t.Errorf("Expected %v, got %v", expectedStatus, lp.Status)
	}

	if lp.OptimalValue != expectedOptimalValue {
		t.Errorf("Expected %v, got %v", expectedOptimalValue, lp.OptimalValue)
	}

	for k, v := range lp.Solution {
		if math.Abs(v-expectedSolution[k]) > 0.0001 {
			t.Errorf("Expected %v, got %v", expectedSolution, lp.Solution)
		}
	}
}

func TestGulpRun(t *testing.T) {
	Gulp()
}
