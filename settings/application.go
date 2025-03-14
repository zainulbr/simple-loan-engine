package settings

type AuthOptions struct {
	EncryptKeys string
}

type ServerOptions struct {
	APIBase       string
	DomainName    string
	HTTPAddress   string
	CorsWhitelist []string
}

type AppSettings struct {
	Name        string
	Version     string
	Description string
	Debug       bool
	UploadDir   string
	Server      ServerOptions
	Auth        AuthOptions
}
