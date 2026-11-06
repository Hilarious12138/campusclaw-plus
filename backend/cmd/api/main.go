package main

import (
	"campusclaw-plus/internal/config"
	"campusclaw-plus/internal/httpapi"
	"campusclaw-plus/internal/lifecycle"
	"campusclaw-plus/internal/repository"
	"campusclaw-plus/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "healthcheck" {
		os.Exit(probe())
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if e := run(logger); e != nil {
		logger.Error("startup_or_shutdown_failed", "module", "runtime", "request_id", "system", "event", "failure", "reason", e.Error())
		os.Exit(1)
	}
}
func run(logger *slog.Logger) error {
	c, e := config.Load(os.Getenv)
	if e != nil {
		return e
	}
	db, e := repository.Open(c.DBAddr, c.DBName, c.DBUser, c.DBPassword, c.DBTimeout)
	if e != nil {
		return fmt.Errorf("database configuration invalid")
	}
	defer db.Close()
	db.TaskLeaseSeconds = c.TaskLeaseSeconds
	db.TaskMaxAttempts = c.TaskMaxAttempts
	db.TaskRetryBaseMS = c.TaskRetryBaseMS
	credentials := repository.Credentials{Secret: []byte(c.JWTSecret), Now: time.Now}
	migrationCtx, migrationCancel := context.WithTimeout(context.Background(), 30*time.Second)
	e = db.Migrate(migrationCtx, credentials, c.SeedPassword)
	migrationCancel()
	if e != nil {
		return fmt.Errorf("platform migration failed")
	}
	dummy, e := credentials.Hash(repository.NewID())
	if e != nil {
		return fmt.Errorf("credential initialization failed")
	}
	auth := service.Auth{Store: db, Credentials: credentials, Now: time.Now, FingerprintKey: []byte(c.JWTSecret), DummyHash: dummy}
	platform := &httpapi.Platform{Log: logger, Auth: auth, Origin: c.Origin, Resources: service.Resources{Store: db}, Materials: service.Materials{Store: db}, Retrieval: service.Retrieval{Store: db, Index: &repository.Index{EmbeddingURL: c.EmbeddingURL, EmbeddingKey: c.EmbeddingKey, EmbeddingModel: c.EmbeddingModel, Dimension: c.EmbeddingDim, QdrantURL: c.QdrantURL, Client: repository.InternalClient()}}}
	platform.MCP = service.MCP{Store: db, Remote: repository.MCPRemote{AllowedURLs: strings.Split(c.MCPAllowed, ","), Token: c.MCPToken}}
	platform.Conversations = service.Conversations{Store: db, Resources: db, MCP: db, Remote: platform.MCP.Remote, Audit: db, Retrieval: platform.Retrieval, Model: repository.ChatModel{BaseURL: c.ChatURL, APIKey: c.ChatKey, Model: c.ChatModel}}
	seedCtx, seedCancel := context.WithTimeout(context.Background(), 10*time.Second)
	e = db.SeedApplications(seedCtx)
	seedCancel()
	if e != nil {
		return fmt.Errorf("application seed failed")
	}
	listener, e := net.Listen("tcp", c.HTTPAddr)
	if e != nil {
		return fmt.Errorf("HTTP listener unavailable")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	s := &http.Server{Handler: httpapi.Handler(service.Health{DB: db, Version: c.Version, Timeout: c.DBTimeout, Now: time.Now}, logger, platform), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 16 << 10}
	logger.Info("started", "module", "runtime", "request_id", "system", "event", "startup", "version", c.Version)
	e = lifecycle.Serve(ctx, s, listener, c.ShutdownTimeout)
	logger.Info("stopped", "module", "runtime", "request_id", "system", "event", "shutdown")
	return e
}
func probe() int {
	_, port, e := net.SplitHostPort(os.Getenv("HTTP_ADDR"))
	if e != nil {
		return 1
	}
	c := http.Client{Timeout: 3 * time.Second}
	r, e := c.Get("http://127.0.0.1:" + port + "/health")
	if e != nil {
		return 1
	}
	defer r.Body.Close()
	var v struct {
		Code string `json:"code"`
	}
	if r.StatusCode != 200 || json.NewDecoder(r.Body).Decode(&v) != nil || v.Code != "OK" {
		return 1
	}
	return 0
}
