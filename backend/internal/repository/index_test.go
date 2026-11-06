package repository

import (
	"campusclaw-plus/internal/domain"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEmbeddingShapeAndBoundary(t *testing.T) {
	for _, response := range []string{`{"data":[{"index":0,"embedding":[1,0]}]}`, `{"data":[{"index":0,"embedding":[0,0]}]}`, `{"data":[]}`, `{"data":[{"index":9,"embedding":[1,0]}]}`} {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer artificial-only" {
				t.Error("missing auth")
			}
			w.Write([]byte(response))
		}))
		i := Index{EmbeddingURL: s.URL, EmbeddingKey: "artificial-only", EmbeddingModel: "test", Dimension: 2, Client: InternalClient()}
		v, e := i.Embed(context.Background(), []string{"文本"})
		if strings.Contains(response, `[1,0]`) && !strings.Contains(response, `index":9`) {
			if e != nil || len(v) != 1 {
				t.Fatal(e)
			}
		} else if e == nil {
			t.Fatal("invalid embedding accepted")
		}
		s.Close()
	}
}
func TestQdrantDoesNotContainTextAndFiltersBothScopes(t *testing.T) {
	var bodies []map[string]any
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Write([]byte(`{}`))
			return
		}
		var b map[string]any
		json.NewDecoder(r.Body).Decode(&b)
		bodies = append(bodies, b)
		if strings.Contains(r.URL.Path, "search") {
			w.Write([]byte(`{"result":[]}`))
		} else {
			w.Write([]byte(`{}`))
		}
	}))
	defer s.Close()
	i := Index{QdrantURL: s.URL, Dimension: 2, Client: InternalClient()}
	e := i.Upsert(context.Background(), "class", "material", []domain.Chunk{{ID: "id", Text: "never-copy-body", Version: "v"}}, [][]float64{{1, 0}})
	if e != nil {
		t.Fatal(e)
	}
	_, e = i.Search(context.Background(), "class", []float64{1, 0}, 5, []string{"material"})
	if e != nil {
		t.Fatal(e)
	}
	raw, _ := json.Marshal(bodies)
	if strings.Contains(string(raw), "never-copy-body") {
		t.Fatal("body in vector payload")
	}
	if !strings.Contains(string(raw), "class_id") || !strings.Contains(string(raw), "material_id") {
		t.Fatal("missing scope filters")
	}
}
func TestConfiguredTimeoutCoversResponseHeaders(t *testing.T) {
	c := InternalClient(35 * time.Second)
	transport := c.Transport.(*http.Transport)
	if c.Timeout != 35*time.Second || transport.ResponseHeaderTimeout != c.Timeout {
		t.Fatal("header timeout silently kept shorter than configured embedding deadline")
	}
}
