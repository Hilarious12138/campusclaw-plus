package main

import (
	"campusclaw-plus/internal/config"
	"campusclaw-plus/internal/repository"
	"campusclaw-plus/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		client := http.Client{Timeout: 3 * time.Second}
		res, e := client.Get("http://127.0.0.1" + os.Getenv("JOB_HEALTH_ADDR") + "/health")
		if e != nil {
			os.Exit(1)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			os.Exit(1)
		}
		return
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if run(logger) != nil {
		logger.Error("job_failed", "event", "failure", "request_id", "system", "module", "job")
		os.Exit(1)
	}
}
func run(log *slog.Logger) error {
	c, e := config.Load(os.Getenv)
	if e != nil {
		return e
	}
	db, e := repository.Open(c.DBAddr, c.DBName, c.DBUser, c.DBPassword, c.DBTimeout)
	if e != nil {
		return e
	}
	defer db.Close()
	db.TaskLeaseSeconds = c.TaskLeaseSeconds
	db.TaskMaxAttempts = c.TaskMaxAttempts
	db.TaskRetryBaseMS = c.TaskRetryBaseMS
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	indexClient := repository.InternalClient()
	if c.AIEnabled {
		indexClient.Timeout = time.Duration(c.EmbeddingTimeoutSeconds) * time.Second
	}
	index := &repository.Index{EmbeddingURL: c.EmbeddingURL, EmbeddingKey: c.EmbeddingKey, EmbeddingModel: c.EmbeddingModel, QdrantURL: c.QdrantURL, Dimension: c.EmbeddingDim, Client: indexClient}
	addr := os.Getenv("JOB_HEALTH_ADDR")
	if addr == "" {
		return fmt.Errorf("required config: JOB_HEALTH_ADDR")
	}
	var ready atomic.Bool
	health := &http.Server{Addr: addr, ReadHeaderTimeout: 3 * time.Second, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/health" {
			w.WriteHeader(404)
			return
		}
		probe, done := context.WithTimeout(r.Context(), 2*time.Second)
		defer done()
		ok := ready.Load() && db.Ping(probe) == nil
		if !ok {
			w.WriteHeader(503)
		}
		json.NewEncoder(w).Encode(map[string]any{"ready": ok})
	})}
	go func() {
		if e := health.ListenAndServe(); e != nil && e != http.ErrServerClosed {
			log.Error("job_health_failed", "event", "failure", "module", "job", "request_id", "system")
			cancel()
		}
	}()
	defer health.Close()
	runner := service.Runner{VectorEnabled: c.AIEnabled, Timeout: time.Duration(c.TaskTimeoutSeconds) * time.Second, Lease: time.Duration(c.TaskLeaseSeconds) * time.Second, Poll: time.Duration(c.TaskPollMS) * time.Millisecond, Store: db, Index: index, Owner: repository.NewID(), Log: log}
	address := url.URL{Scheme: "amqp", Host: c.BrokerAddr, User: url.UserPassword(c.BrokerUser, c.BrokerPassword), Path: "/"}
	for ctx.Err() == nil {
		q, e := repository.OpenRabbit(address.String())
		if e != nil {
			log.Warn("broker_retry", "module", "job", "event", "broker_retry", "request_id", "system")
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(5 * time.Second):
				continue
			}
		}
		if raw := os.Getenv("TEST_ACK_DELAY_MS"); raw != "" {
			n, e := strconv.Atoi(raw)
			if c.Environment != "test" || e != nil || n < 0 || n > 30000 {
				q.Close()
				return fmt.Errorf("invalid test-only configuration")
			}
			q.AckDelay = time.Duration(n) * time.Millisecond
		}
		ready.Store(true)
		if len(os.Args) > 1 && os.Args[1] == "dispatcher" {
			e = runner.Dispatch(ctx, q)
		} else {
			e = q.Consume(ctx, runner.Execute)
		}
		ready.Store(false)
		q.Close()
		if ctx.Err() != nil {
			return nil
		}
		if e != nil {
			log.Warn("job_retry", "module", "job", "event", "retry", "request_id", "system")
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(5 * time.Second):
			}
		}
	}
	return nil
}
