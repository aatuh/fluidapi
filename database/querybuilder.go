package database

// InsertedValues returns the columns and values to insert.
type InsertedValues func() (columns []string, values []any)

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

// ColumnDefinition defines the properties for a table column in a
// table creation query.
type ColumnDefinition struct {
	Name          string  // Column name
	Type          string  // Data type (with length/precision, e.g. "CHAR(36)")
	NotNull       bool    // Whether to add NOT NULL (if false, NULL is allowed)
	Default       *string // Optional default value (pass nil if not needed)
	AutoIncrement bool    // Whether to add AUTO_INCREMENT
	Extra         string  // Extra column options (e.g. "CHARACTER SET utf8mb4 COLLATE utf8mb4_bin")
	PrimaryKey    bool    // If true, marks this column as primary key (inline)
	Unique        bool    // If true, marks this column as UNIQUE (unless already primary key)
}

// TableOptions holds additional options for a  table creation query.
type TableOptions struct {
	Engine  string // e.g. "InnoDB"
	Charset string // e.g. "utf8mb4"
	Collate string // e.g. "utf8mb4_bin"
}

// QueryBuilder defines an interface for generating SQL queries.
// This interface abstracts all query-generation logic.
// In all methods an error is returned if the query generation failed but in
// case the query is not supported, an empty string is returned without error.
type QueryBuilder interface {
	Insert(tableName string, insertedValues InsertedValues) (string, []any)
	InsertMany(
		tableName string,
		insertedValues []InsertedValues,
	) (string, []any)
	UpsertMany(
		tableName string,
		insertedValues []InsertedValues,
		updateProjections []Projection,
	) (string, []any)
	Get(tableName string, options *GetOptions) (string, []any)
	Count(tableName string, options *CountOptions) (string, []any)
	UpdateQuery(
		tableName string,
		updateFields []UpdateField,
		selectors []Selector,
	) (string, []any)
	Delete(
		tableName string,
		selectors []Selector,
		opts *DeleteOptions,
	) (string, []any)
	CreateDatabaseQuery(
		dbName string,
		ifNotExists bool,
		charset string,
		collate string,
	) (string, []any, error)
	CreateTableQuery(
		tableName string,
		ifNotExists bool,
		columns []ColumnDefinition,
		constraints []string,
		options TableOptions,
	) (string, []any, error)
	UseDatabaseQuery(dbName string) (string, []any, error)
	SetVariableQuery(variable string, value string) (string, []any, error)
	AdvisoryLock(lockName string, timeout int) (string, []any, error)
	AdvisoryUnlock(lockName string) (string, []any, error)
}
