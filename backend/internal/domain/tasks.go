package domain

import (
	"context"
	"time"
)

type Material struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	ClassID       string    `json:"class_id"`
	CreatedBy     string    `json:"created_by"`
	TaskID        string    `json:"task_id"`
	Status        string    `json:"status"`
	ActiveVersion string    `json:"active_version,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	Content       string    `json:"-"`
	InputHash     string    `json:"-"`
}
type Task struct {
	ID         string     `json:"id"`
	MaterialID string     `json:"material_id"`
	ClassID    string     `json:"-"`
	CreatedBy  string     `json:"-"`
	Status     string     `json:"status"`
	Stage      string     `json:"stage"`
	Attempts   int        `json:"attempts"`
	Deliveries int        `json:"delivery_count"`
	Fence      int        `json:"-"`
	Owner      string     `json:"-"`
	ErrorCode  string     `json:"error_code,omitempty"`
	NextRunAt  time.Time  `json:"next_run_at"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	Hash       string     `json:"-"`
}
type Outbox struct{ ID, TaskID, Kind, Owner string }
type Chunk struct {
	ID         string  `json:"chunk_id"`
	MaterialID string  `json:"material_id"`
	Title      string  `json:"material_title"`
	Index      int     `json:"chunk_index"`
	Start      int     `json:"start_char"`
	End        int     `json:"end_char"`
	Text       string  `json:"excerpt"`
	Version    string  `json:"-"`
	Score      float64 `json:"score"`
}
type TaskStore interface {
	CreateMaterial(context.Context, Identity, string, string, string, Audit) (Material, Task, error)
	Materials(context.Context, Identity, int, int) ([]Material, error)
	Material(context.Context, Identity, string) (Material, error)
	Task(context.Context, Identity, string) (Task, error)
	ClaimOutbox(context.Context, string) (Outbox, error)
	Published(context.Context, Outbox) error
	ClaimTask(context.Context, string, string) (Task, Material, error)
	RenewTask(context.Context, Task) error
	CompleteTask(context.Context, Task, []Chunk) error
	FailTask(context.Context, Task, string, bool) error
	RecoverTasks(context.Context) error
	Orphans(context.Context) ([]Material, error)
	MarkCleaned(context.Context, string) error
	CleanupExpired(context.Context) error
}
type Queue interface {
	Publish(context.Context, Outbox) error
	Consume(context.Context, func(context.Context, string) error) error
	Close() error
}
type Index interface {
	Delete(context.Context, []string) error
	Embed(context.Context, []string) ([][]float64, error)
	Upsert(context.Context, string, string, []Chunk, [][]float64) error
	Search(context.Context, string, []float64, int, []string) ([]ScoredID, error)
}
type ScoredID struct {
	ID    string
	Score float64
}
type RetrievalStore interface {
	Keyword(context.Context, Identity, string, []string, int) ([]Chunk, error)
	ChunksByID(context.Context, Identity, []string, []string) ([]Chunk, error)
	Chunk(context.Context, Identity, string) (Chunk, error)
}
