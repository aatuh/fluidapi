package clause

// Projection represents a projected column in a query.
type Projection struct {
	Table  string
	Column string
	Alias  string
}
