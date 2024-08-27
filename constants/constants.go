package constants

func Runner() bool {
	return true
}

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
	LpStatusNotSolved  = LpStatus(0)
	LpStatusOptimal    = LpStatus(1)
	LpStatusInfeasible = LpStatus(2)
	LpStatusUnbounded  = LpStatus(3)
	LpStatusUndefined  = LpStatus(4)
)

// TODO: Status Mapping ?? dictionary??
// TODO: SolutionStatus, and Status to SolutionStatus Mapping

type LpConstraint int

const (
	LpConstraintLE = LpConstraint(-1)
	LpConstraintEQ = LpConstraint(0)
	LpConstraintGE = LpConstraint(1)
)

// TODO: Constraint Mapping ?? dictionary??
