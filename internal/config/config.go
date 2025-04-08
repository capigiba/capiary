package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds the entire application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
}

// ServerConfig holds server-related configurations.
type ServerConfig struct {
	Port string

	// TODO: implement later
	// Language string
	// I18NPath string `mapstructure:"i18n_path"`
}

// DatabaseConfig holds database-related configurations.
type DatabaseConfig struct {
	PostgresURL string `mapstructure:"postgres_url"`
}

// CORSConfig holds CORS-related configurations.
type CORSConfig struct {
	AllowedOrigins   []string      `mapstructure:"allowed_origins"`
	AllowedMethods   []string      `mapstructure:"allowed_methods"`
	AllowedHeaders   []string      `mapstructure:"allowed_headers"`
	ExposeHeaders    []string      `mapstructure:"expose_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           time.Duration `mapstructure:"max_age"`
}

var AppConfig *Config

// LoadConfig initializes the configuration by reading from environment variables and config files.
func LoadConfig() (*Config, error) {
	env := "development"
	envFile := fmt.Sprintf(".env.%s", env)
	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Printf("No %s file found, relying on existing environment variables\n", envFile)
	}

	v := viper.New()

	// Set the file name of the configurations file
	v.SetConfigName("config")

	// Set the path to look for the configurations file
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read in the config file if it exists
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("No config file found, relying on environment variables")
	}

	// Define default values for non-sensitive fields
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.language", "en")
	v.SetDefault("server.i18n_path", "../../i18n")
	v.SetDefault("cors.allowed_origins", []string{"http://localhost:3000"})
	v.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allowed_headers", []string{"Origin", "Content-Type", "Authorization"})
	v.SetDefault("cors.expose_headers", []string{"Content-Length"})
	v.SetDefault("cors.allow_credentials", true)
	v.SetDefault("cors.max_age", 43200) // 12 hours in seconds

	// Bind environment variables to specific config keys
	v.BindEnv("database.postgres_url", "POSTGRES_URL")

	// Unmarshal the config into the Config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %v", err)
	}

	// Convert CORS.MaxAge from seconds to time.Duration
	config.CORS.MaxAge = time.Duration(config.CORS.MaxAge) * time.Second

	// Validate required configurations
	missing := []string{}
	if config.Database.PostgresURL == "" {
		missing = append(missing, "POSTGRES_URL")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	AppConfig = &config
	return AppConfig, nil
}
