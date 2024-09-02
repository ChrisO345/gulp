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
}

type Row struct {
	Values []float64
}

type Column struct {
	Values []float64
}

func NewTableau(lp *LinearProgram) *Tableau {
	tableau := &Tableau{}

	tableau.NamesRow = make([]string, len(lp.ObjectiveFunction.Pairs))
	tableau.ObjectiveRow = Row{Values: make([]float64, len(lp.ObjectiveFunction.Pairs))}
	for i, v := range lp.ObjectiveFunction.Pairs {
		tableau.NamesRow[i] = v.Variable.Name
		tableau.ObjectiveRow.Values[i] = v.Coefficient
	}

	tableau.ConstraintRows = make([]Row, len(lp.Constraints))
	tableau.BasisNames = make([]string, len(lp.Constraints))
	tableau.BasisColumn = Column{Values: make([]float64, len(lp.Constraints))}
	tableau.BColumn = Column{Values: make([]float64, len(lp.Constraints))}

	for i, v := range lp.Constraints {
		tableau.ConstraintRows[i] = Row{Values: make([]float64, len(lp.ObjectiveFunction.Pairs))}
		tableau.BasisNames[i] = v.Pairs[len(v.Pairs)-1].Variable.Name
		tableau.BasisColumn.Values[i] = 0
		tableau.BColumn.Values[i] = lp.Constraints[i].RightHandSide
		for _, p := range v.Pairs {
			for k, o := range lp.ObjectiveFunction.Pairs {
				if o.Variable.Name == p.Variable.Name {
					tableau.ConstraintRows[i].Values[k] = p.Coefficient
				}
			}
		}
	}

	tableau.ZRow = Row{Values: make([]float64, len(lp.ObjectiveFunction.Pairs))}
	tableau.CZRow = Row{Values: make([]float64, len(lp.ObjectiveFunction.Pairs))}
	for i := range lp.ObjectiveFunction.Pairs {
		tableau.ZRow.Values[i] = 0
		tableau.CZRow.Values[i] = lp.ObjectiveFunction.Pairs[i].Coefficient
	}

	tableau.Sense = lp.ObjectiveFunction.ObjectiveType
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
		if t.Sense == LpMinimise {
			boolean = v < t.CZRow.Values[pivotColumnIndex]
		}

		if boolean {
			pivotColumnIndex = i
		}
	}

	// Find the pivot row
	pivotRowIndex := 0
	// TODO: this might need to change depending on the sense of the problem
	optimumColumnRatio := 0.0
	if t.Sense == LpMinimise {
		optimumColumnRatio = math.Inf(1)
	}
	for i, v := range t.BColumn.Values {
		if t.ConstraintRows[i].Values[pivotColumnIndex] > 0 {
			ratio := v / t.ConstraintRows[i].Values[pivotColumnIndex]
			boolean := ratio < optimumColumnRatio
			if t.Sense == LpMinimise {
				boolean = ratio > optimumColumnRatio
			}
			if boolean || i == 0 {
				optimumColumnRatio = ratio
				pivotRowIndex = i
			}
		}
	}

	// Update the basis
	t.BasisNames[pivotRowIndex] = t.NamesRow[pivotColumnIndex]
	t.BasisColumn.Values[pivotRowIndex] = t.ObjectiveRow.Values[pivotColumnIndex]

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
