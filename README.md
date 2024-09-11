# gulp

*A Simplex Method Linear Programming Solver in Go.*

**gulp** is a minimal, educational linear programming solver built using the simplex method. It’s designed to be easy to understand and straightforward to use.

___

## Features

- **Simplex Method**: Uses the simplex method to solve linear programming problems.
- **Minimization and Maximization**: Can solve both minimization and maximization problems.
- **Simple Interface**: Designed to be easy to use and understand.

___

## Installation

To install **gulp**, you need to have Go installed on your machine. Then, you can run the following command:

```bash
go get github.com/chriso345/gulp
```

___

## Usage

Lets set up a basic linear programming problem using **gulp**:

```go
package main

import "github.com/chriso345/gulp"

func main() {
    // Create decision variables
    x1 := gulp.NewVariable("x1")
    x2 := gulp.NewVariable("x2")
    x3 := gulp.NewVariable("x3")

    // Objective function: Minimize -6 * x1 + 7 * x2 + 4 * x3
    objective := gulp.NewExpression([]gulp.LpTerm{
        gulp.NewTerm(-6, x1),
        gulp.NewTerm(7, x2),
        gulp.NewTerm(4, x3),
    })

    // Set up the LP problem
    lp := gulp.NewLinearProgram()
    lp.AddObjective(gulp.LpMinimise, objective)

    // Add constraints
    lp.AddConstraint(gulp.NewExpression([]gulp.LpTerm{
        gulp.NewTerm(2, x1),
        gulp.NewTerm(5, x2),
        gulp.NewTerm(-1, x3),
    }), gulp.LpConstraintLE, 18)

    lp.AddConstraint(gulp.NewExpression([]gulp.LpTerm{
        gulp.NewTerm(1, x1),
        gulp.NewTerm(-1, x2),
        gulp.NewTerm(-2, x3),
    }), gulp.LpConstraintLE, -14)

    lp.AddConstraint(gulp.NewExpression([]gulp.LpTerm{
        gulp.NewTerm(3, x1),
        gulp.NewTerm(2, x2),
        gulp.NewTerm(2, x3),
    }), gulp.LpConstraintEQ, 26)

    // Solve it
    lp.Solve().PrintSolution()
}
```

### What's Happening Here?

We're defining a linear programming problem where:

1. **Objective Function**: We want to minimize $-6x_1 + 7x_2 + 4x_3$
2. **Constraints**:
    - $2x_1 + 5x_2 - x_3 \leq 18$
    - $x_1 - x_2 - 2x_3 \leq -14$
    - $3x_1 + 2x_2 + 2x_3 = 26$

When we run this program, **gulp** will solve the linear programming problem and print the solution, and optimal values for $x_1$, $x_2$, and $x_3$.

___

## The API

### Decision Variables

These are the variables whose values we wish to determine.

```go
x1 := gulp.NewVariable("x1")
x2 := gulp.NewVariable("x2")    
```

Each variable is created using `gulp.NewVariable()` with a name that uniquely identifies it. In this case `x1` and `x2` represent the decision variables.

### Objective Function

The objective function is the mathematical expression that we want to minimize or maximize.

```go
objective := gulp.NewExpression([]gulp.LpTerm{
    gulp.NewTerm(-6, x1),
    gulp.NewTerm(7, x2),
    gulp.NewTerm(4, x3),
})
```

- We use `gulp.NewTerm()` to specify each term in the objective function. Each term consists of:
  - A coefficient (the weight or importance of the variable in the objective function)
  - A variable (the decision variable it corresponds to)
- We then use `gulp.NewExpression()` to create the objective function using the terms we defined.

```go
lp := gulp.NewLinearProgram()
lp.AddObjective(gulp.LpMinimise, objective)
```

- After creating the objective function, we add it to the linear program using `lp.AddObjective()`.
- The first argument specifies whether we want to minimize or maximize the objective function. We use `gulp.LpMinimise` or `gulp.LpMaximise` for this.

### Constraints

Constraints are the conditions that the decision variables must satisfy. They can be equality or inequality constraints.

```go
lp.AddConstraint(gulp.NewExpression([]gulp.LpTerm{
    gulp.NewTerm(2, x1),
    gulp.NewTerm(5, x2),
    gulp.NewTerm(-1, x3),
}), gulp.LpConstraintLE, 18)
```

Each constraint is composed of:
- **Expression**: Defined similarly to the objective function using `gulp.NewTerm()` and `gulp.NewExpression()`.
- **Type**: The type of constraint. We use `gulp.LpConstraintLE` for less than or equal to ($\leq$), `gulp.LpConstraintGE` for greater than or equal to ($\geq$), and `gulp.LpConstraintEQ` for equality ($=$).
- **Right-hand Side**: The value on the right-hand side of the constraint.

### Solving the Problem

```go
solution := lp.Solve()
solution.PrintSolution()
```

- `lp.Solve()` will run the simplex algorithm to find the optimal solution based on the objective function and constraints.
- `solution.PrintSolution()` will print the optimal values of the decision variables and the optimal value of the objective function.

___ 

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
