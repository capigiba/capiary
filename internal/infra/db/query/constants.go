package query

type OperationType string

const (
	OpEqual       OperationType = "=="
	OpNotEqual    OperationType = "!="
	OpGreaterThan OperationType = ">"
	OpLessThan    OperationType = "<"
	OpGTE         OperationType = ">="
	OpLTE         OperationType = "<="
)
