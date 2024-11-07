package config

import (
    "time"
)

type Config struct {
    Server    ServerConfig    `mapstructure:"server"`
    Database  DatabaseConfig  `mapstructure:"database"`
    JWT       JWTConfig      `mapstructure:"jwt"`
    CORS      CORSConfig     `mapstructure:"cors"`
    RateLimit RateLimitConfig `mapstructure:"rate_limit"`
    Security  SecurityConfig  `mapstructure:"security"`
    App       AppConfig      `mapstructure:"app"`
    Logging   LoggingConfig  `mapstructure:"logging"`
    Cache     CacheConfig    `mapstructure:"cache"`
    Monitoring MonitoringConfig `mapstructure:"monitoring"`
    Email     EmailConfig    `mapstructure:"email"`
    Storage   StorageConfig  `mapstructure:"storage"`
    Features  FeaturesConfig `mapstructure:"features"`
}

type ServerConfig struct {
    Port            string          `mapstructure:"port"`
    Host            string          `mapstructure:"host"`
    Timeout         TimeoutConfig   `mapstructure:"timeout"`
    ShutdownTimeout time.Duration   `mapstructure:"shutdown_timeout"`
}

type TimeoutConfig struct {
    Read  time.Duration `mapstructure:"read"`
    Write time.Duration `mapstructure:"write"`
    Idle  time.Duration `mapstructure:"idle"`
}

type DatabaseConfig struct {
    Host            string        `mapstructure:"host"`
    Port            string        `mapstructure:"port"`
    User            string        `mapstructure:"user"`
    Password        string        `mapstructure:"password"`
    Name            string        `mapstructure:"name"`
    MaxOpenConns    int          `mapstructure:"max_open_conns"`
    MaxIdleConns    int          `mapstructure:"max_idle_conns"`
    ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
    SSLMode         string        `mapstructure:"ssl_mode"`
    MaxRetries      int          `mapstructure:"max_retries"`
}

type JWTConfig struct {
    SecretKey             string        `mapstructure:"secret_key"`
    RefreshKey             string        `mapstructure:"refresh_key"`
    AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry"`
    RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
    Issuer            string        `mapstructure:"issuer"`
    Audience          string        `mapstructure:"audience"`
}

type CORSConfig struct {
    AllowedOrigins    []string      `mapstructure:"allowed_origins"`
    AllowedMethods    []string      `mapstructure:"allowed_methods"`
    AllowedHeaders    []string      `mapstructure:"allowed_headers"`
    ExposedHeaders    []string      `mapstructure:"exposed_headers"`
    AllowCredentials  bool          `mapstructure:"allow_credentials"`
    MaxAge           int           `mapstructure:"max_age"`
}

type RateLimitConfig struct {
    Enabled   bool          `mapstructure:"enabled"`
    Requests  int           `mapstructure:"requests"`
    Duration  time.Duration `mapstructure:"duration"`
    Type      string        `mapstructure:"type"`
}

type SecurityConfig struct {
    BCryptCost           int                    `mapstructure:"bcrypt_cost"`
    MinPasswordLength    int                    `mapstructure:"min_password_length"`
    MaxPasswordLength    int                    `mapstructure:"max_password_length"`
    AllowedSpecialChars  string                 `mapstructure:"allowed_special_chars"`
    PasswordRequirements PasswordRequirements   `mapstructure:"password_requirements"`
}

type PasswordRequirements struct {
    MinUppercase int `mapstructure:"min_uppercase"`
    MinLowercase int `mapstructure:"min_lowercase"`
    MinNumbers   int `mapstructure:"min_numbers"`
    MinSpecial   int `mapstructure:"min_special"`
}

type AppConfig struct {
    Name        string    `mapstructure:"name"`
    Environment string    `mapstructure:"environment"`
    Debug       bool      `mapstructure:"debug"`
    API         APIConfig `mapstructure:"api"`
    Timezone    string    `mapstructure:"timezone"`
}

type APIConfig struct {
    Version string `mapstructure:"version"`
    Prefix  string `mapstructure:"prefix"`
}

type LoggingConfig struct {
    Level    string         `mapstructure:"level"`
    Format   string         `mapstructure:"format"`
    Output   string         `mapstructure:"output"`
    FilePath string         `mapstructure:"file_path"`
    Include  LoggingInclude `mapstructure:"include"`
}

type LoggingInclude struct {
    Timestamp bool `mapstructure:"timestamp"`
    Level     bool `mapstructure:"level"`
    Caller    bool `mapstructure:"caller"`
    RequestID bool `mapstructure:"request_id"`
}

type CacheConfig struct {
    Enabled          bool          `mapstructure:"enabled"`
    Type            string        `mapstructure:"type"`
    TTL             time.Duration `mapstructure:"ttl"`
    CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

type MonitoringConfig struct {
    Enabled         bool   `mapstructure:"enabled"`
    MetricsPath    string `mapstructure:"metrics_path"`
    HealthCheckPath string `mapstructure:"health_check_path"`
}

type EmailConfig struct {
    Enabled      bool         `mapstructure:"enabled"`
    SMTP         SMTPConfig   `mapstructure:"smtp"`
    From         FromConfig   `mapstructure:"from"`
    TemplatesDir string       `mapstructure:"templates_dir"`
}

type SMTPConfig struct {
    Host       string `mapstructure:"host"`
    Port       int    `mapstructure:"port"`
    Username   string `mapstructure:"username"`
    Password   string `mapstructure:"password"`
    Encryption string `mapstructure:"encryption"`
}

type FromConfig struct {
    Name  string `mapstructure:"name"`
    Email string `mapstructure:"email"`
}

type StorageConfig struct {
    Type   string           `mapstructure:"type"`
    Local  LocalStorage     `mapstructure:"local"`
    S3     S3Storage       `mapstructure:"s3"`
}

type LocalStorage struct {
    Path    string `mapstructure:"path"`
    MaxSize string `mapstructure:"max_size"`
}

type S3Storage struct {
    Region    string `mapstructure:"region"`
    Bucket    string `mapstructure:"bucket"`
    AccessKey string `mapstructure:"access_key"`
    SecretKey string `mapstructure:"secret_key"`
}

type FeaturesConfig struct {
    EnableRefreshTokens     bool `mapstructure:"enable_refresh_tokens"`
    EnableSocialLogin      bool `mapstructure:"enable_social_login"`
    EnableMFA             bool `mapstructure:"enable_mfa"`
    EnablePasswordReset    bool `mapstructure:"enable_password_reset"`
    EnableEmailVerification bool `mapstructure:"enable_email_verification"`
    EnableUserDeletion     bool `mapstructure:"enable_user_deletion"`
}
