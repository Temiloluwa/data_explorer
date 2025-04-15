package domain

// Question represents a user's query processed internally.
type Question struct {
	Query   string
	// Sources defines specific sources requested by the user, if any.
	// If empty, the engine might query all configured/applicable sources.
	Sources []DataSourceSpec
	// Options provides additional parameters, potentially source-specific.
	Options map[string]string
}

// DataSourceSpec identifies a requested data source type and potentially a specific instance.
type DataSourceSpec struct {
	Type       string // e.g., "web", "local", "db_postgres"
	Identifier string // Optional: e.g., specific DB connection name, file path pattern
}
