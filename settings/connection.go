package settings

type PostgresOption struct {
	Driver  string
	URI     string
	Enabled bool
}

type ConnectionSettings struct {
	Postgres PostgresOption
}
