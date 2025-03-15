package settings

type PostgresOption struct {
	Driver  string
	URI     string
	Enabled bool
}

type SMTPOption struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type ConnectionSettings struct {
	Postgres PostgresOption
	SMTP     SMTPOption
}
