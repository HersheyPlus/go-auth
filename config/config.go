package config

import (
    "fmt"
    "github.com/spf13/viper"
)

func LoadConfig() (*Config, error) {
    v := viper.New()
    setDefaults(v)
    v.SetConfigName("config")
    v.SetConfigType("yaml")

    v.AddConfigPath(".")
    v.AddConfigPath("./config")

    v.AutomaticEnv()
    v.SetEnvPrefix("APP")

    if err := v.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, fmt.Errorf("error reading config file: %s", err)
        }
    }

    var config Config
    if err := v.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("unable to decode config into struct: %s", err)
    }

    if err := validateConfig(&config); err != nil {
        return nil, err
    }
    return &config, nil
}


func setDefaults(v *viper.Viper) {
    // Server defaults
    v.SetDefault("server.port", "8080")
    v.SetDefault("server.host", "0.0.0.0")
    v.SetDefault("server.timeout.read", "15s")
    v.SetDefault("server.timeout.write", "15s")
    v.SetDefault("server.timeout.idle", "60s")
    v.SetDefault("server.shutdown_timeout", "30s")

    // Database defaults
    v.SetDefault("database.host", "localhost")
    v.SetDefault("database.port", "5432")
    v.SetDefault("database.max_open_conns", 100)
    v.SetDefault("database.max_idle_conns", 10)
    v.SetDefault("database.conn_max_lifetime", "1h")
    v.SetDefault("database.ssl_mode", "disable")

    // Rate limit defaults
    v.SetDefault("rate_limit.enabled", true)
    v.SetDefault("rate_limit.requests", 100)
    v.SetDefault("rate_limit.duration", "1m")
    v.SetDefault("rate_limit.type", "ip")

    // Security defaults
    v.SetDefault("security.bcrypt_cost", 12)
    v.SetDefault("security.min_password_length", 8)
    v.SetDefault("security.max_password_length", 72)

    // App defaults
    v.SetDefault("app.environment", "development")
    v.SetDefault("app.debug", true)
    v.SetDefault("app.timezone", "UTC")
    v.SetDefault("app.api.version", "v1")
    v.SetDefault("app.api.prefix", "/api")

    // Logging defaults
    v.SetDefault("logging.level", "info")
    v.SetDefault("logging.format", "json")
    v.SetDefault("logging.output", "stdout")
}

func validateConfig(config *Config) error {
    if config.JWT.Secret == "" {
        return fmt.Errorf("jwt secret is required")
    }

    if config.Database.User == "" || config.Database.Password == "" || config.Database.Name == "" {
        return fmt.Errorf("database configuration is incomplete")
    }

    // Validate rate limit configuration
    if config.RateLimit.Enabled {
        if config.RateLimit.Requests <= 0 {
            return fmt.Errorf("rate limit requests must be greater than 0")
        }
        if config.RateLimit.Duration <= 0 {
            return fmt.Errorf("rate limit duration must be greater than 0")
        }
    }

    return nil
}

func (c *Config) GetDBConnString() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Database.Host,
        c.Database.Port,
        c.Database.User,
        c.Database.Password,
        c.Database.Name,
        c.Database.SSLMode,
    )
}
