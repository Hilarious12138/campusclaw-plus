package service

import (
	"campusclaw-plus/internal/domain"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode/utf8"
)

type Materials struct{ Store domain.TaskStore }

func (s Materials) Create(ctx context.Context, u domain.Identity, title, body, key string, event domain.Audit) (domain.Material, domain.Task, error) {
	if u.Role != "teacher" {
		return domain.Material{}, domain.Task{}, domain.ErrForbidden
	}
	if len(body) == 0 || len(body) > 2<<20 || !utf8.ValidString(body) || strings.TrimSpace(body) == "" || len([]rune(title)) > 255 || len(title) == 0 || strings.ContainsAny(title, "/\\\r\n\x00") || (!strings.HasSuffix(strings.ToLower(title), ".txt") && !strings.HasSuffix(strings.ToLower(title), ".md")) || len(key) < 1 || len(key) > 128 {
		return domain.Material{}, domain.Task{}, domain.ErrInvalid
	}
	for _, c := range key {
		if c < 33 || c > 126 {
			return domain.Material{}, domain.Task{}, domain.ErrInvalid
		}
	}
	return s.Store.CreateMaterial(ctx, u, title, body, key, event)
}
func (s Materials) List(ctx context.Context, u domain.Identity, limit, offset int) ([]domain.Material, error) {
	return s.Store.Materials(ctx, u, limit, offset)
}
func (s Materials) Read(ctx context.Context, u domain.Identity, id string) (domain.Material, error) {
	return s.Store.Material(ctx, u, id)
}
func (s Materials) Task(ctx context.Context, u domain.Identity, id string) (domain.Task, error) {
	return s.Store.Task(ctx, u, id)
}

type Runner struct {
	VectorEnabled        bool
	Timeout, Lease, Poll time.Duration
	Store                domain.TaskStore
	Index                domain.Index
	Owner                string
	Log                  *slog.Logger
}

func (r Runner) Dispatch(ctx context.Context, q domain.Queue) error {
	ticker := time.NewTicker(r.Poll)
	defer ticker.Stop()
	ticks := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			ticks++
			if ticks%60 == 0 {
				if e := r.Store.CleanupExpired(ctx); e != nil {
					return e
				}
				materials, e := r.Store.Orphans(ctx)
				if e != nil {
					return e
				}
				for _, m := range materials {
					_, _, parts := splitText(m.Content, ChunkOptions{Strategy: "auto", MaxLength: 800, OverlapRatio: 0.1})
					ids := []string{}
					for i := range parts {
						ids = append(ids, chunkID(m.ID, m.InputHash, i+1))
					}
					if e = r.Index.Delete(ctx, ids); e != nil {
						return e
					}
					if e = r.Store.MarkCleaned(ctx, m.TaskID); e != nil {
						return e
					}
				}
			}
			if ticks%5 == 0 {
				if e := r.Store.RecoverTasks(ctx); e != nil {
					return e
				}
			}
			for n := 0; n < 100; n++ {
				v, e := r.Store.ClaimOutbox(ctx, r.Owner)
				if errors.Is(e, domain.ErrNotFound) {
					break
				}
				if e != nil {
					return e
				}
				pub, cancel := context.WithTimeout(ctx, 5*time.Second)
				e = q.Publish(pub, v)
				cancel()
				if e != nil {
					return e
				}
				if e = r.Store.Published(ctx, v); e != nil {
					return e
				}
			}
		}
	}
}
func chunkID(mid, version string, index int) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d", mid, version, index)))
	x := hex.EncodeToString(h[:16])
	return x[:8] + "-" + x[8:12] + "-" + x[12:16] + "-" + x[16:20] + "-" + x[20:]
}
func (r Runner) Execute(ctx context.Context, id string) error {
	t, m, e := r.Store.ClaimTask(ctx, id, r.Owner)
	if errors.Is(e, domain.ErrConflict) || errors.Is(e, domain.ErrNotFound) {
		return nil
	}
	if e != nil {
		if errors.Is(e, domain.ErrForbidden) && t.Status == "running" {
			return r.Store.FailTask(ctx, t, "PERMISSION_REVOKED", false)
		}
		return e
	}
	job, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()
	done := make(chan struct{})
	defer close(done)
	go func() {
		ticker := time.NewTicker(min(10*time.Second, r.Lease/3))
		defer ticker.Stop()
		for {
			select {
			case <-job.Done():
				return
			case <-done:
				return
			case <-ticker.C:
				if e := r.Store.RenewTask(job, t); e != nil {
					cancel()
					return
				}
			}
		}
	}()
	start := time.Now()
	chunks := []domain.Chunk{}
	text := []string{}
	_, _, parts := splitText(m.Content, ChunkOptions{Strategy: "auto", MaxLength: 800, OverlapRatio: 0.1})
	for i, p := range parts {
		chunks = append(chunks, domain.Chunk{ID: chunkID(m.ID, t.Hash, i+1), MaterialID: m.ID, Index: i + 1, Start: p.Start, End: p.End, Text: p.Text, Version: t.Hash, Title: m.Title})
		text = append(text, p.Text)
	}
	var vectors [][]float64
	if len(text) == 0 {
		e = domain.ErrInvalid
	} else if r.VectorEnabled {
		vectors, e = r.Index.Embed(job, text)
	}
	if e == nil && r.VectorEnabled {
		e = r.Index.Upsert(job, t.ClassID, t.MaterialID, chunks, vectors)
	}
	if e == nil {
		e = r.Store.CompleteTask(job, t, chunks)
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if e != nil {
		if errors.Is(e, domain.ErrConflict) {
			return e
		}
		retry := !errors.Is(e, domain.ErrInvalid) && !errors.Is(e, domain.ErrForbidden) && !errors.Is(e, domain.ErrUnauthenticated)
		failure, c := context.WithTimeout(ctx, 5*time.Second)
		defer c()
		code := "INDEX_DEPENDENCY_FAILED"
		if !retry {
			code = "INDEX_INPUT_REJECTED"
		}
		if fe := r.Store.FailTask(failure, t, code, retry); fe != nil {
			return fe
		}
	}
	r.Log.Info("task_finished", "event", "task_attempt", "module", "worker", "request_id", id, "task_id", id, "attempts", t.Attempts, "duration_ms", time.Since(start).Milliseconds(), "failed", e != nil)
	return nil
}
