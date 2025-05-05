package query

// Filter represents a filtering condition like "field operator value".
type Filter struct {
	Field    string
	Operator OperationType
	Value    interface{}
}

// Sort represents sorting instructions like "field ASC/DESC".
type Sort struct {
	Field string
	Desc  bool
}

// QueryOptions is a container for all query customizations.
type QueryOptions struct {
	Filters []Filter
	Sorts   []Sort
	Fields  []string
	Skip    int64
	Limit   int64
}
