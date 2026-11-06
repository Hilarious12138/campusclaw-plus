package config

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment, Version, HTTPAddr, DBAddr, DBName, DBUser, DBPassword string
	JWTSecret, Origin, SeedPassword                                    string
	DBTimeout, ShutdownTimeout                                         time.Duration
}

func Load(get func(string) string) (Config, error) {
	var c Config
	for _, v := range []struct {
		k string
		p *string
	}{{"APP_ENV", &c.Environment}, {"APP_VERSION", &c.Version}, {"HTTP_ADDR", &c.HTTPAddr}, {"DB_ADDR", &c.DBAddr}, {"DB_NAME", &c.DBName}, {"DB_USER", &c.DBUser}, {"DB_PASSWORD", &c.DBPassword}, {"JWT_SECRET", &c.JWTSecret}, {"APP_ORIGIN", &c.Origin}, {"SEED_PASSWORD", &c.SeedPassword}} {
		*v.p = get(v.k)
		if strings.TrimSpace(*v.p) == "" {
			return c, fmt.Errorf("required config: %s", v.k)
		}
	}
	if c.Environment != "local" && c.Environment != "test" && c.Environment != "production" {
		return c, fmt.Errorf("invalid config: APP_ENV")
	}
	for k, v := range map[string]string{"HTTP_ADDR": c.HTTPAddr, "DB_ADDR": c.DBAddr} {
		_, port, e := net.SplitHostPort(v)
		n, ne := strconv.Atoi(port)
		if e != nil || ne != nil || n < 1 || n > 65535 {
			return c, fmt.Errorf("invalid config: %s", k)
		}
	}
	u, ue := url.Parse(c.Origin)
	if ue != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return c, fmt.Errorf("invalid config: APP_ORIGIN")
	}
	if len(c.JWTSecret) < 32 || len(c.SeedPassword) < 16 {
		return c, fmt.Errorf("invalid credential configuration")
	}
	if len(c.DBPassword) < 24 {
		return c, fmt.Errorf("invalid config: DB_PASSWORD")
	}
	n, e := strconv.Atoi(get("DB_TIMEOUT_MS"))
	if e != nil || n < 100 || n > 10000 {
		return c, fmt.Errorf("invalid config: DB_TIMEOUT_MS")
	}
	c.DBTimeout = time.Duration(n) * time.Millisecond
	n, e = strconv.Atoi(get("SHUTDOWN_TIMEOUT_SECONDS"))
	if e != nil || n < 1 || n > 10 {
		return c, fmt.Errorf("invalid config: SHUTDOWN_TIMEOUT_SECONDS")
	}
	c.ShutdownTimeout = time.Duration(n) * time.Second
	return c, nil
}
