package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app"`
	Auth     Auth     `mapstructure:"auth"`
	Postgres Postgres `mapstructure:"postgres"`
}

type App struct {
	Environment  string        `envconfig:"APP_ENVIRONMENT" default:"dev" required:"true" mapstructure:"environment"`
	ServiceName  string        `envconfig:"APP_NAME" default:"Movie" required:"true" mapstructure:"app_name"`
	Host         string        `envconfig:"APP_HOST" default:"localhost" required:"true" mapstructure:"host"`
	Port         int           `envconfig:"APP_PORT" default:"8000" required:"true" mapstructure:"port"`
	LogLevel     string        `envconfig:"APP_LOG_LEVEL" default:"info" required:"true" mapstructure:"log_level"`
	GracePeriod  time.Duration `envconfig:"APP_GRACE_PERIOD" default:"4s" required:"true" mapstructure:"grace_period"`
	Debug        bool          `envconfig:"APP_DEBUG" default:"false" mapstructure:"debug"`
	ReadTimeout  time.Duration `envconfig:"APP_READ_TIMEOUT" default:"10s" mapstructure:"read_timeout"`
	WriteTimeout time.Duration `envconfig:"APP_WRITE_TIMEOUT" default:"10s" mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `envconfig:"APP_IDLE_TIMEOUT" default:"120s" mapstructure:"idle_timeout"`
}

type Auth struct {
	AccessTTL        time.Duration `envconfig:"AUTH_ACCESS_TTL" default:"30m" required:"true" mapstructure:"access_ttl"`
	RefreshTTL       time.Duration `envconfig:"AUTH_REFRESH_TTL" default:"2h" required:"true" mapstructure:"refresh_ttl"`
	OtpTTL           time.Duration `envconfig:"AUTH_OTP_TTL" default:"1m" required:"true" mapstructure:"otp_ttl"`
	CodeLength       int           `envconfig:"AUTH_CODE_LENGTH" default:"6" required:"true" mapstructure:"code_length"`
	SecretKey        string        `envconfig:"AUTH_KEY" default:"auth_secret_key" required:"true" mapstructure:"secret_key"`
	MaxLoginAttempts int           `envconfig:"AUTH_MAX_LOGIN_ATTEMPTS" default:"5" mapstructure:"max_login_attempts"`
	LockoutDuration  time.Duration `envconfig:"AUTH_LOCKOUT_DURATION" default:"15m" mapstructure:"lockout_duration"`
}

type Postgres struct {
	Host            string        `envconfig:"POSTGRES_HOST" default:"localhost" required:"true" mapstructure:"host"`
	Port            int           `envconfig:"POSTGRES_PORT" default:"5432" required:"true" mapstructure:"port"`
	User            string        `envconfig:"POSTGRES_USER" default:"postgres" required:"true" mapstructure:"user"`
	Password        string        `envconfig:"POSTGRES_PASSWORD" default:"password" required:"true" mapstructure:"password"`
	Database        string        `envconfig:"POSTGRES_DATABASE" default:"i_tv_task" required:"true" mapstructure:"database"`
	SslMode         string        `envconfig:"POSTGRES_SSLMODE" default:"disable" required:"true" mapstructure:"sslmode"`
	MaxOpenConns    int           `envconfig:"POSTGRES_MAX_OPEN_CONNS" default:"25" mapstructure:"max_open_conns"`
	MaxIdleConns    int           `envconfig:"POSTGRES_MAX_IDLE_CONNS" default:"5" mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `envconfig:"POSTGRES_CONN_MAX_LIFETIME" default:"5m" mapstructure:"conn_max_lifetime"`
}

const configDir = "config"

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	configPaths := []string{configDir, "."}
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	env := os.Getenv("APP_ENVIRONMENT")
	if env == "" {
		env = "dev"
	}
	v.SetConfigName(env)

	envFile := fmt.Sprintf("%s/%s.env", configDir, env)
	if err := godotenv.Load(envFile); err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.Wrapf(err, "failed to load %s.env file", env)
		}
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.Wrap(err, "failed to read config")
		}
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.Wrap(err, "failed to process env config")
	}

	if v.ConfigFileUsed() != "" {
		if err := v.Unmarshal(&cfg); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal config")
		}
	}

	return &cfg, nil
}

func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.Database,
		c.Postgres.SslMode,
	)
}

func (c *Config) GetAppAddr() string {
	return fmt.Sprintf("%s:%d", c.App.Host, c.App.Port)
}