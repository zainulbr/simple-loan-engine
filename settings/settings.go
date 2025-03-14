package settings

const (
	EnvDev      = "dev"
	EnvTest     = "test"
	EnvUnitTest = "unit-test"
	EnvProd     = "prd"
)

type Settings struct {
	App  AppSettings
	Conn ConnectionSettings
}
