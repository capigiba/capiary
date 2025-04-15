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
	Storage  StorageConfig
	CORS     CORSConfig
}

type StorageConfig struct {
	AwsRegion      string `mapstructure:"aws_region"`
	AwsBucket      string `mapstructure:"aws_bucket"`
	AwsAccessKeyID string `mapstructure:"aws_access_key_id"`
	AwsSecretKey   string `mapstructure:"aws_secret_key"`
}

// ServerConfig holds server-related configurations.
type ServerConfig struct {
	Port      string
	JWTSecret string `mapstructure:"jwt_secret"`
	// TODO: implement later
	// Language string
	// I18NPath string `mapstructure:"i18n_path"`
}

// DatabaseConfig holds database-related configurations.
type DatabaseConfig struct {
	PostgresURL    string `mapstructure:"postgres_url"`
	RdsPostgresURL string `mapstructure:"rds_postgres_url"`
	MongodbURI     string `mapstructure:"mongodb_uri"`
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

	// Bind environment variables to specific config keys
	v.BindEnv("database.postgres_url", "POSTGRES_URL")
	v.BindEnv("database.rds_postgres_url", "RDS_POSTGRES_ENDPOINT")
	v.BindEnv("database.mongodb_uri", "MONGODB_ENDPOINT")
	v.BindEnv("server.jwt_secret", "JWT_SECRET")
	v.BindEnv("storage.aws_region", "AWS_REGION")
	v.BindEnv("storage.aws_bucket", "AWS_BUCKET")
	v.BindEnv("storage.aws_access_key_id", "AWS_ACCESS_KEY_ID")
	v.BindEnv("storage.aws_secret_key", "AWS_SECRET_KEY")

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
