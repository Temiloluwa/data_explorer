package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/yourusername/data-explorer/internal/platform/logger" // Adjust path if needed
)

// Config holds the application configuration.
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Log         logger.Config     `mapstructure:"log"`
	DataSources DataSourcesConfig `mapstructure:"datasources"`
	// Add other top-level config sections if needed
}

// ServerConfig holds server-related settings.
type ServerConfig struct {
	GRPCPort string `mapstructure:"grpc_port"`
	RESTPort string `mapstructure:"rest_port"`
	// Add TLS config etc. if needed (e.g., CertFile, KeyFile)
}

// DataSourcesConfig holds configurations for various data sources.
// Use map[string]interface{} for flexibility or define specific structs.
type DataSourcesConfig struct {
	WebSearch  WebSearchConfig  `mapstructure:"web_search"`
	LocalFiles LocalFilesConfig `mapstructure:"local_files"`
	// Add structs for DB, external APIs etc. following the pattern:
	// DBPostgres PostgresConfig `mapstructure:"db_postgres"`
	// APISomeAPI SomeAPIConfig  `mapstructure:"api_someapi"`
}

// WebSearchConfig example
type WebSearchConfig struct {
	Enable bool   `mapstructure:"enable"`
	APIKey string `mapstructure:"api_key"` // Consider loading secrets securely (env var, vault)
	// Other options like search engine choice, results limit
}

// LocalFilesConfig example
type LocalFilesConfig struct {
	Enable     bool     `mapstructure:"enable"`
	WatchPaths []string `mapstructure:"watch_paths"` // Directories to scan
	// Other options like allowed extensions, indexing strategy
}

// PostgresConfig example
// type PostgresConfig struct {
//  Enable           bool   `mapstructure:"enable"`
// 	ConnectionString string `mapstructure:"connection_string"` // Load securely!
// }

// SomeAPIConfig example
// type SomeAPIConfig struct {
//  Enable   bool   `mapstructure:"enable"`
// 	BaseURL  string `mapstructure:"base_url"`
// 	APIKey   string `mapstructure:"api_key"` // Load securely!
// }


// LoadConfig loads configuration from file, env vars etc. using Viper.
func LoadConfig(log logger.Logger, configPath string) (*Config, error) {
	v := viper.New() // Use a local viper instance

	// 1. Set Defaults
	setDefaults(v)

	// 2. Read from config file
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.AddConfigPath("./configs") // Look in ./configs first
		v.AddConfigPath(".")         // Then look in current directory
		v.SetConfigName("config")    // Name of config file (without extension)
		v.SetConfigType("yaml")      // Or json, toml etc.
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired and rely on defaults/env vars
			log.Warnf("Config file not found (checked paths './configs', '.'); using defaults/env vars.")
		} else {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("failed to read config file '%s': %w", v.ConfigFileUsed(), err)
		}
	} else {
		log.Infof("Using config file: %s", v.ConfigFileUsed())
	}

	// 3. Read from Environment Variables
	v.SetEnvPrefix("DATAEXPLORER") // e.g., DATAEXPLORER_SERVER_GRPC_PORT=...
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // Replace . and - with _ for env vars
	v.AutomaticEnv()

	// 4. Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	// 5. Post-load validation (optional but recommended)
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}


	// Log loaded sensitive info carefully (e.g., only if API key is set, not the value)
	logIfSet := func(key, value string) {
		if value != "" {
			log.Debugf("Config: %s is set (from file or env)", key)
		} else {
			log.Debugf("Config: %s is NOT set", key)
		}
	}
	logIfSet("datasources.web_search.api_key", cfg.DataSources.WebSearch.APIKey)
    // Add similar logging for other sensitive fields if needed


	return &cfg, nil
}

// setDefaults defines the default configuration values.
func setDefaults(v *viper.Viper) {
	v.SetDefault("log.level", "info")
	v.SetDefault("server.grpc_port", ":50051")
	v.SetDefault("server.rest_port", ":8080") // Default REST port
	v.SetDefault("datasources.web_search.enable", false)
	v.SetDefault("datasources.local_files.enable", false)
	v.SetDefault("datasources.local_files.watch_paths", []string{"./documents"}) // Default local path
    // Set defaults for other data sources if added
}

// validateConfig performs basic validation checks on the loaded config.
func validateConfig(cfg *Config) error {
	if cfg.Server.GRPCPort == "" {
		return fmt.Errorf("server.grpc_port must be set")
	}
	if cfg.Server.RESTPort == "" {
		return fmt.Errorf("server.rest_port must be set")
	}
	if cfg.DataSources.WebSearch.Enable && cfg.DataSources.WebSearch.APIKey == "" {
		// Warn instead of error? Depends on requirements.
		// return fmt.Errorf("datasources.web_search is enabled but api_key is missing")
        fmt.Println("WARN: datasources.web_search is enabled but api_key is missing (set DATAEXPLORER_DATASOURCES_WEB_SEARCH_API_KEY)")
	}
    if cfg.DataSources.LocalFiles.Enable && len(cfg.DataSources.LocalFiles.WatchPaths) == 0 {
        fmt.Println("WARN: datasources.local_files is enabled but watch_paths is empty")
    }
	// Add more validation rules as needed
	return nil
}
