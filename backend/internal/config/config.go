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
	AIEnabled                                                                                     bool
	EmbeddingTimeoutSeconds                                                                       int
	Environment, Version, HTTPAddr, DBAddr, DBName, DBUser, DBPassword                            string
	BrokerAddr, BrokerUser, BrokerPassword, EmbeddingURL, EmbeddingKey, EmbeddingModel, QdrantURL string
	ChatURL, ChatKey, ChatModel, MCPToken, MCPAllowed                                             string
	TaskTimeoutSeconds, TaskLeaseSeconds, TaskMaxAttempts, TaskRetryBaseMS, TaskPollMS            int
	EmbeddingDim                                                                                  int
	SeedPasswords                                                                                 map[string]string
	JWTSecret, Origin                                                                             string
	DBTimeout, ShutdownTimeout                                                                    time.Duration
}

func Load(get func(string) string) (Config, error) {
	var c Config
	for _, v := range []struct {
		k string
		p *string
	}{{"APP_ENV", &c.Environment}, {"APP_VERSION", &c.Version}, {"HTTP_ADDR", &c.HTTPAddr}, {"DB_ADDR", &c.DBAddr}, {"DB_NAME", &c.DBName}, {"DB_USER", &c.DBUser}, {"DB_PASSWORD", &c.DBPassword}, {"JWT_SECRET", &c.JWTSecret}, {"APP_ORIGIN", &c.Origin}, {"BROKER_ADDR", &c.BrokerAddr}, {"BROKER_USER", &c.BrokerUser}, {"BROKER_PASSWORD", &c.BrokerPassword}, {"QDRANT_URL", &c.QdrantURL}, {"MCP_AUTH_TOKEN", &c.MCPToken}, {"MCP_ALLOWED_BASE_URLS", &c.MCPAllowed}} {
		*v.p = get(v.k)
		if strings.TrimSpace(*v.p) == "" {
			return c, fmt.Errorf("required config: %s", v.k)
		}
	}
	if c.Environment != "local" && c.Environment != "test" && c.Environment != "production" {
		return c, fmt.Errorf("invalid config: APP_ENV")
	}
	for k, v := range map[string]string{"HTTP_ADDR": c.HTTPAddr, "DB_ADDR": c.DBAddr, "BROKER_ADDR": c.BrokerAddr} {
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
	c.SeedPasswords = map[string]string{}
	seenPasswords := map[string]bool{}
	for _, name := range []string{"teacher_a", "student_a", "teacher_b", "student_b"} {
		key := "SEED_" + strings.ToUpper(name) + "_PASSWORD"
		value := get(key)
		if len(value) < 16 || seenPasswords[value] {
			return c, fmt.Errorf("invalid independent seed password: %s", key)
		}
		seenPasswords[value] = true
		c.SeedPasswords[name] = value
	}
	if len(c.JWTSecret) < 32 {
		return c, fmt.Errorf("invalid credential configuration")
	}
	if len(c.DBPassword) < 24 {
		return c, fmt.Errorf("invalid config: DB_PASSWORD")
	}
	switch get("AI_ENABLED") {
	case "true":
		c.AIEnabled = true
	case "false":
	default:
		return c, fmt.Errorf("required boolean config: AI_ENABLED")
	}
	for key, target := range map[string]*string{"EMBEDDING_BASE_URL": &c.EmbeddingURL, "EMBEDDING_API_KEY": &c.EmbeddingKey, "EMBEDDING_MODEL": &c.EmbeddingModel, "CHAT_BASE_URL": &c.ChatURL, "CHAT_API_KEY": &c.ChatKey, "CHAT_MODEL": &c.ChatModel} {
		*target = get(key)
		if c.AIEnabled && strings.TrimSpace(*target) == "" {
			return c, fmt.Errorf("required config: %s", key)
		}
	}
	if c.AIEnabled {
		dim, de := strconv.Atoi(get("EMBEDDING_DIM"))
		if de != nil || dim < 2 || dim > 8192 {
			return c, fmt.Errorf("invalid config: EMBEDDING_DIM")
		}
		c.EmbeddingDim = dim
		n, e := strconv.Atoi(get("EMBEDDING_TIMEOUT_SECONDS"))
		if e != nil || n < 1 || n > 60 {
			return c, fmt.Errorf("invalid config: EMBEDDING_TIMEOUT_SECONDS")
		}
		c.EmbeddingTimeoutSeconds = n
	}
	bases := map[string]string{"QDRANT_URL": c.QdrantURL}
	if c.AIEnabled {
		bases["EMBEDDING_BASE_URL"] = c.EmbeddingURL
		bases["CHAT_BASE_URL"] = c.ChatURL
	}
	for key, raw := range bases {
		u, e := url.Parse(raw)
		if e != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Scheme != "http" && u.Scheme != "https") {
			return c, fmt.Errorf("invalid config: %s", key)
		}
	}
	if len(c.MCPToken) < 32 {
		return c, fmt.Errorf("invalid config: MCP_AUTH_TOKEN")
	}
	if len(c.BrokerPassword) < 24 {
		return c, fmt.Errorf("invalid config: BROKER_PASSWORD")
	}
	for key, target := range map[string]*int{"TASK_TIMEOUT_SECONDS": &c.TaskTimeoutSeconds, "TASK_LEASE_SECONDS": &c.TaskLeaseSeconds, "TASK_MAX_ATTEMPTS": &c.TaskMaxAttempts, "TASK_RETRY_BASE_MS": &c.TaskRetryBaseMS, "TASK_POLL_MS": &c.TaskPollMS} {
		n, e := strconv.Atoi(get(key))
		if e != nil {
			return c, fmt.Errorf("invalid config: %s", key)
		}
		*target = n
	}
	if c.TaskTimeoutSeconds < 10 || c.TaskTimeoutSeconds > 600 || c.TaskLeaseSeconds < c.TaskTimeoutSeconds+5 || c.TaskLeaseSeconds > 3600 || c.TaskMaxAttempts < 1 || c.TaskMaxAttempts > 3 || c.TaskRetryBaseMS < 100 || c.TaskRetryBaseMS > 60000 || c.TaskPollMS < 50 || c.TaskPollMS > 10000 {
		return c, fmt.Errorf("invalid task policy configuration")
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
