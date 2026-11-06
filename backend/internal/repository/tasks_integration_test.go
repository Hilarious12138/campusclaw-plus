package repository

import (
	"campusclaw-plus/internal/domain"
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"
)

func TestTaskConcurrencyAndFencing(t *testing.T) {
	if os.Getenv("PLUS_INTEGRATION") != "1" {
		t.Skip("requires isolated MySQL, executed by Docker grader B04")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db, e := Open("db:3306", os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), 5*time.Second)
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	u, e := db.LoginUser(ctx, "teacher_a")
	if e != nil {
		t.Fatal(e)
	}
	key := NewID()
	name := "fence-" + key + ".md"
	event := domain.Audit{ActorID: u.ID, ClassID: u.ClassID, Operation: "integration.task", HTTPStatus: 202, Outcome: "ACCEPTED", RequestID: key}
	var wg sync.WaitGroup
	ids := make(chan string, 20)
	errs := make(chan error, 20)
	for n := 0; n < 20; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, task, e := db.CreateMaterial(ctx, u.Identity, name, "concurrency fixture", key, event)
			if e != nil {
				errs <- e
			} else {
				ids <- task.ID
			}
		}()
	}
	wg.Wait()
	close(ids)
	close(errs)
	for e := range errs {
		t.Fatal(e)
	}
	one := ""
	for id := range ids {
		if one == "" {
			one = id
		}
		if one != id {
			t.Fatal("duplicate business task")
		}
	}
	old, _, e := db.ClaimTask(ctx, one, NewID())
	if e != nil {
		t.Fatal(e)
	}
	if _, e = db.db.ExecContext(ctx, "UPDATE tasks SET lease_until=TIMESTAMPADD(SECOND,-1,UTC_TIMESTAMP(6)) WHERE id=?", one); e != nil {
		t.Fatal(e)
	}
	if e = db.RecoverTasks(ctx); e != nil {
		t.Fatal(e)
	}
	next, _, e := db.ClaimTask(ctx, one, NewID())
	if e != nil {
		t.Fatal(e)
	}
	chunks := []domain.Chunk{{ID: NewID(), Index: 1, Start: 0, End: 7, Text: "fixture"}}
	if e = db.CompleteTask(ctx, old, chunks); !errors.Is(e, domain.ErrConflict) {
		t.Fatal("stale worker committed", e)
	}
	if e = db.CompleteTask(ctx, next, chunks); e != nil {
		t.Fatal(e)
	}
	var n int
	if e = db.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM task_results WHERE task_id=?", one).Scan(&n); e != nil || n != 1 {
		t.Fatal("result uniqueness", e, n)
	}
	if _, e = db.db.ExecContext(ctx, "DELETE FROM task_outbox WHERE task_id=?", one); e != nil {
		t.Fatal(e)
	}
}
