package gulp

func Gulp() {
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
		AddConstraint(NewExpression(terms4), LpConstraintEQ, 26)

	//fmt.Println(lp.String())

	lp.Solve().PrintSolution()

	// lp.Shadows()
}
