package settings

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	env string
	s   *Settings
)

// S returns global application settings
func S() *Settings {
	if s == nil {
		log.Fatalln("application settings not loaded")
	}
	return s
}

// Env returns the current application environment
func Env() string {
	return env
}

// Load global application settings, specify altPaths to add alternative config search paths
func Load() (*Settings, error) {
	env = getEnv("APP_ENV", EnvUnitTest)

	if env == EnvUnitTest {
		return &Settings{}, nil
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	// Parse environment variables into Settings struct
	settings := &Settings{
		App: AppSettings{
			Name:        getEnv("APP_NAME", "MyApp"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Description: getEnv("APP_DESCRIPTION", ""),
			Debug:       getEnvAsBool("APP_DEBUG", false),
			UploadDir:   getEnv("UPLOAD_DIR", "uploads"),
			Server: ServerOptions{
				APIBase:       getEnv("SERVER_API_BASE", "/api"),
				DomainName:    getEnv("SERVER_DOMAIN", "localhost"),
				HTTPAddress:   getEnv("SERVER_HTTP_ADDRESS", ":8080"),
				CorsWhitelist: getEnvAsSlice("SERVER_CORS_WHITELIST", ","),
			},
			Auth: AuthOptions{
				EncryptKeys: getEnv("AUTH_ENCRYPT_KEYS", ""),
			},
		},
		Conn: ConnectionSettings{
			Postgres: PostgresOption{
				Driver:  getEnv("POSTGRES_DRIVER", "postgres"),
				URI:     getEnv("POSTGRES_URI", ""),
				Enabled: getEnvAsBool("POSTGRES_ENABLED", true),
			},
		},
	}

	return settings, nil
}

// Helper functions to get env variables with default values
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		durationValue, err := time.ParseDuration(value)
		if err == nil {
			return durationValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key, sep string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, sep)
	}
	return []string{}
}
