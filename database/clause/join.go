package clause

// JoinType represents the type of join
type JoinType string

const (
	JoinTypeInner JoinType = "INNER"
	JoinTypeLeft  JoinType = "LEFT"
	JoinTypeRight JoinType = "RIGHT"
	JoinTypeFull  JoinType = "FULL"
)

// Join represents a database join clause
type Join struct {
	Type    JoinType
	Table   string
	OnLeft  ColumnSelector
	OnRight ColumnSelector
}

// ColumnSelector represents a columnn selector
type ColumnSelector struct {
	Table   string
	Columnn string
}
