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
	"syscall"
	"time"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "healthcheck" {
		os.Exit(probe())
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("module", "runtime", "request_id", "system")
	if e := run(logger); e != nil {
		logger.Error("startup_or_shutdown_failed", "event", "failure", "reason", e.Error())
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
	listener, e := net.Listen("tcp", c.HTTPAddr)
	if e != nil {
		return fmt.Errorf("HTTP listener unavailable")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	s := &http.Server{Handler: httpapi.Handler(service.Health{DB: db, Version: c.Version, Timeout: c.DBTimeout, Now: time.Now}, logger), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 16 << 10}
	logger.Info("started", "event", "startup", "version", c.Version)
	e = lifecycle.Serve(ctx, s, listener, c.ShutdownTimeout)
	logger.Info("stopped", "event", "shutdown")
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
