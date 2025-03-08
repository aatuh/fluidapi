package database

// GetOptions is used for get queries.
type GetOptions struct {
	Selectors   Selectors
	Orders      Orders
	Page        *Page
	Joins       Joins
	Projections Projections
	Lock        bool
}

// CountOptions is used for count queries.
type CountOptions struct {
	Selectors Selectors
	Page      *Page
	Joins     Joins
}

// DeleteOptions is used for delete queries.
type DeleteOptions struct {
	Limit  int
	Orders Orders
}

// ColumnDefinition defines the properties for a table column in a table.
// creation query.
type ColumnDefinition struct {
	Name          string  // Column name
	Type          string  // Data type (with length/precision, e.g. "CHAR(36)")
	NotNull       bool    // Whether to add NOT NULL (if false, NULL is allowed)
	Default       *string // Optional default value (pass nil if not needed)
	AutoIncrement bool    // Whether to add AUTO_INCREMENT
	Extra         string  // Extra column options (e.g. "CHARACTER SET utf8mb4 COLLATE utf8mb4_bin")
	PrimaryKey    bool    // Marks this column as primary key (inline)
	Unique        bool    // Marks this column as UNIQUE (unless already primary key)
}

// TableOptions holds additional options for a table creation query.
type TableOptions struct {
	Engine  string // e.g. "InnoDB"
	Charset string // e.g. "utf8mb4"
	Collate string // e.g. "utf8mb4_bin"
}

// InsertedValuesFn defines a function that returns column names and values for
// an insert. This allows deferred evaluation of values and consistent ordering
// of parameters.
type InsertedValuesFn func() ([]string, []any)

// QueryBuilder defines an interface for building SQL queries dynamically.
// Implementations handle specifics for different SQL dialects.
type QueryBuilder interface {
	// Insert builds an INSERT statement for a single row.
	// The insertedValuesFunc should produce the column names and values for the row.
	Insert(table string, insertedValuesFunc InsertedValuesFn) (query string, params []any)
	// InsertMany builds a batch INSERT for multiple rows.
	InsertMany(table string, valuesFuncs []InsertedValuesFn) (query string, params []any)
	// UpsertMany builds an UPSERT (insert or update) statement for multiple rows.
	UpsertMany(table string, valuesFuncs []InsertedValuesFn, updateProjections []Projection) (query string, params []any)
	// Get builds a SELECT statement with optional filtering, ordering, and limits.
	Get(table string, options *GetOptions) (query string, params []any)
	// Count builds a SELECT COUNT(*) statement with optional filters.
	Count(table string, options *CountOptions) (query string, params []any)
	// UpdateQuery builds an UPDATE statement for given selectors and update fields.
	UpdateQuery(table string, updateFields []UpdateField, selectors []Selector) (query string, params []any)
	// Delete builds a DELETE statement for given selectors.
	Delete(table string, selectors []Selector, opts *DeleteOptions) (query string, params []any)
	// CreateDatabaseQuery builds a CREATE DATABASE statement.
	CreateDatabaseQuery(dbName string, ifNotExists bool, charset string, collate string) (string, []any, error)
	// CreateTableQuery builds a CREATE TABLE statement.
	CreateTableQuery(tableName string, ifNotExists bool, columns []ColumnDefinition, constraints []string, options TableOptions) (string, []any, error)
	// UseDatabaseQuery builds a USE DATABASE statement.
	UseDatabaseQuery(dbName string) (string, []any, error)
	// SetVariableQuery builds a SET statement for a variable.
	SetVariableQuery(variable string, value string) (string, []any, error)
	// AdvisoryLock builds an advisory lock statement.
	AdvisoryLock(lockName string, timeout int) (string, []any, error)
	// AdvisoryUnlock builds an advisory unlock statement.
	AdvisoryUnlock(lockName string) (string, []any, error)
}
