package database

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// TODO: Add SSL configuration.

type Config struct {
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USER"`
	Host     string `env:"DB_HOST, default=localhost"`
	Port     string `env:"DB_PORT, default=5432"`
	Password string `env:"DB_PASSWORD" json:"-"`

	PoolMaxConns    int           `env:"DB_POOL_MAX_CONNS, default=10"`
	PoolMinConns    int           `env:"DB_POOL_MIN_CONNS, default=3"`
	PoolMaxConnLife time.Duration `env:"DB_POOL_MAX_CONN_LIFETIME, default=1h"`
	PoolMaxConnIdle time.Duration `env:"DB_POOL_MAX_CONN_IDLE_TIME, default=30m"`
	PoolHealthCheck time.Duration `env:"DB_POOL_HEALTH_CHECK_PERIOD, default=1m"`
}

func (cfg *Config) connectionURL() string {
	if cfg == nil {
		return ""
	}

	host := cfg.Host
	if port := cfg.Port; port != "" {
		host = host + ":" + port
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   cfg.Name,
	}

	if cfg.User != "" || cfg.Password != "" {
		u.User = url.UserPassword(cfg.User, cfg.Password)
	}

	return u.String()
}

func (cfg *Config) dsn() string {
	values := cfg.mapValues()

	pairs := make([]string, 0, len(values))

	for key, value := range values {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(pairs, " ")
}

func (cfg *Config) mapValues() map[string]string {
	values := map[string]string{}

	setIfNotEmpty := func(key, val string) {
		if val != "" {
			values[key] = val
		}
	}

	setIfPositive := func(key string, val int) {
		if val > 0 {
			values[key] = fmt.Sprintf("%d", val)
		}
	}

	setIfPositiveDuration := func(key string, d time.Duration) {
		if d > 0 {
			values[key] = d.String()
		}
	}

	// Standard connection parameters
	setIfNotEmpty("dbname", cfg.Name)
	setIfNotEmpty("user", cfg.User)
	setIfNotEmpty("host", cfg.Host)
	setIfNotEmpty("password", cfg.Password)
	setIfNotEmpty("port", cfg.Port)

	// Pool parameters
	setIfPositive("pool_max_conns", cfg.PoolMaxConns)
	setIfPositive("pool_min_conns", cfg.PoolMinConns)
	setIfPositiveDuration("pool_max_conn_lifetime", cfg.PoolMaxConnLife)
	setIfPositiveDuration("pool_max_conn_idle_time", cfg.PoolMaxConnIdle)
	setIfPositiveDuration("pool_health_check_period", cfg.PoolHealthCheck)

	return values
}
