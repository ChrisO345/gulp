package gulp

import (
	"fmt"
	"math"
)

type Tableau struct {
	// Rows
	NamesRow       []string
	ObjectiveRow   Row
	ConstraintRows []Row

	// Basis
	BasisNames  []string
	BasisColumn Column
	BColumn     Column

	// Solutions
	ZRow  Row
	CZRow Row

	// Optimality
	TableauValue float64

	Sense LpSense

	Variables []LpVariable
}

type Row struct {
	Values []float64
}

type Column struct {
	Values []float64
}

func NewTableau(lp *LinearProgram) *Tableau {
	tableau := &Tableau{}

	// Create the names row and objective row
	tableau.NamesRow = make([]string, len(lp.ObjectiveFunction.Terms))
	tableau.ObjectiveRow = Row{Values: make([]float64, len(lp.ObjectiveFunction.Terms))}
	for i, v := range lp.ObjectiveFunction.Terms {
		tableau.NamesRow[i] = v.Variable.Name
		tableau.ObjectiveRow.Values[i] = v.Coefficient
	}

	tableau.ConstraintRows = make([]Row, len(lp.Constraints))
	tableau.BasisNames = make([]string, len(lp.Constraints))
	tableau.BasisColumn = Column{Values: make([]float64, len(lp.Constraints))}
	tableau.BColumn = Column{Values: make([]float64, len(lp.Constraints))}

	// Create the constraint rows
	for i, v := range lp.Constraints {
		tableau.ConstraintRows[i] = Row{Values: make([]float64, len(lp.ObjectiveFunction.Terms))}
		tableau.BColumn.Values[i] = lp.Constraints[i].RightHandSide
		for _, p := range v.Terms {
			for k, o := range lp.ObjectiveFunction.Terms {
				if o.Variable.Name == p.Variable.Name {
					tableau.ConstraintRows[i].Values[k] = p.Coefficient
				}
			}
		}
	}

	// Create the basis column names and values
	for i, r := range tableau.ConstraintRows {
		for j, v := range r.Values {
			if v == 1 {
				boolean := true
				for k := 0; k < len(tableau.BasisColumn.Values); k++ {
					if k == i {
						continue
					}
					if tableau.ConstraintRows[k].Values[j] != 0 {
						boolean = false
						break
					}
				}
				if boolean {
					tableau.BasisNames[i] = tableau.NamesRow[j]
					tableau.BasisColumn.Values[i] = tableau.ObjectiveRow.Values[j]
				}
			}
		}
	}

	// Create the Z row
	tableau.ZRow = Row{Values: make([]float64, len(lp.ObjectiveFunction.Terms))}
	tableau.CZRow = Row{Values: make([]float64, len(lp.ObjectiveFunction.Terms))}
	for i, r := range tableau.ConstraintRows {
		for j, v := range r.Values {
			tableau.ZRow.Values[j] += v * tableau.BasisColumn.Values[i]
		}
	}

	// Create the CZ row
	for i, v := range tableau.ZRow.Values {
		tableau.CZRow.Values[i] = tableau.ObjectiveRow.Values[i] - v
	}

	// Calculate the tableau value
	for i, v := range tableau.BColumn.Values {
		tableau.TableauValue += v * tableau.BasisColumn.Values[i]
	}

	// Add all the variables to the tableau
	for _, v := range lp.ObjectiveFunction.Terms {
		tableau.Variables = append(tableau.Variables, v.Variable)
	}

	tableau.Sense = lp.hiddenSense
	return tableau
}

// TODO: Implement the String() method for the Tableau struct
func (t *Tableau) String() string {
	fmt.Println("Not implemented")
	return "Not implemented"
}

// Pivot TODO: Take the largest negative value in the CZRow as the pivot column
func (t *Tableau) Pivot() {
	// Find the pivot column
	pivotColumnIndex := 0
	for i, v := range t.CZRow.Values {
		boolean := v > t.CZRow.Values[pivotColumnIndex]

		if boolean {
			pivotColumnIndex = i
		}
	}

	// Find the pivot row
	pivotRowIndex := 0
	// TODO: this might need to change depending on the sense of the problem
	optimumColumnRatio := math.Inf(1)
	for i, v := range t.BColumn.Values {
		if t.ConstraintRows[i].Values[pivotColumnIndex] > 0 {
			ratio := v / t.ConstraintRows[i].Values[pivotColumnIndex]
			boolean := ratio < optimumColumnRatio && ratio > 0
			if boolean {
				optimumColumnRatio = ratio
				pivotRowIndex = i
			}
		}
	}

	oldBasisName := t.BasisNames[pivotRowIndex]
	// Update the basis
	t.BasisNames[pivotRowIndex] = t.NamesRow[pivotColumnIndex]
	t.BasisColumn.Values[pivotRowIndex] = t.ObjectiveRow.Values[pivotColumnIndex]

	// If variable being replaced is artificial, remove it from the basis
	// TODO: Have a better way to check if a variable is artificial, as this is a hack
	boolean := false
	for _, v := range t.Variables {
		if v.Name == oldBasisName {
			boolean = v.IsArtificial
			break
		}
	}

	if boolean {
		for i, v := range t.NamesRow {
			if v == oldBasisName {
				t.ObjectiveRow.Values[i] = 0
				for j := range t.ConstraintRows {
					t.ConstraintRows[j].Values[i] = 0
				}
				t.ZRow.Values[i] = 0
				t.CZRow.Values[i] = 0
			}
		}
	}

	// Update the pivot row
	pivotRow := t.ConstraintRows[pivotRowIndex]
	pivotRowValue := pivotRow.Values[pivotColumnIndex]
	for i := range pivotRow.Values {
		pivotRow.Values[i] /= pivotRowValue
	}
	t.BColumn.Values[pivotRowIndex] /= pivotRowValue

	// Update the pivot column
	for i := range t.ConstraintRows {
		if i != pivotRowIndex {
			multiplier := t.ConstraintRows[i].Values[pivotColumnIndex]
			for j := range t.ConstraintRows[i].Values {
				t.ConstraintRows[i].Values[j] -= multiplier * pivotRow.Values[j]
			}
			t.BColumn.Values[i] -= multiplier * t.BColumn.Values[pivotRowIndex]
		}
	}

	// Update the Z row
	for i := range t.ZRow.Values {
		val := 0.0
		for j := range t.ConstraintRows {
			val += t.ConstraintRows[j].Values[i] * t.BasisColumn.Values[j]
		}
		t.ZRow.Values[i] = val
	}

	// Update the objective row
	multiplier := t.CZRow.Values[pivotColumnIndex]
	for i := range t.CZRow.Values {
		t.CZRow.Values[i] -= multiplier * pivotRow.Values[i]
	}

	// Update the tableau value
	t.TableauValue = 0
	for i, v := range t.BColumn.Values {
		t.TableauValue += v * t.BasisColumn.Values[i]
	}
}

func (t *Tableau) IsOptimal() bool {
	for _, v := range t.CZRow.Values {
		if v > 0 {
			return false
		}
	}
	return true
}

func (t *Tableau) GetSolution() map[string]float64 {
	solution := make(map[string]float64)
	for i, v := range t.BasisNames {
		solution[v] = t.BColumn.Values[i]
	}
	return solution
}

func (t *Tableau) Log() {
	// Print Debug
	fmt.Println("Objective Names:", t.NamesRow)
	fmt.Println("Objective Row:", t.ObjectiveRow)

	for _, v := range t.ConstraintRows {
		fmt.Println("Constraint Row:", v)
	}

	fmt.Println("Basis Names:", t.BasisNames)
	fmt.Println("Basis Column:", t.BasisColumn)
	fmt.Println("B Column:", t.BColumn)

	fmt.Println("Z Row:", t.ZRow)
	fmt.Println("CZ Row:", t.CZRow)

	fmt.Println("Tableau Value:", t.TableauValue)
	fmt.Println()
}
