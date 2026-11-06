package repository

import (
	"bytes"
	"campusclaw-plus/internal/domain"
	"context"
	"encoding/json"
	"errors"

	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Index struct {
	EmbeddingURL, EmbeddingKey, EmbeddingModel, QdrantURL string
	Dimension                                             int
	Client                                                *http.Client
}

func InternalClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second, Transport: &http.Transport{Proxy: nil, ResponseHeaderTimeout: 30 * time.Second}, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
}
func validBase(s string) bool {
	u, e := url.Parse(s)
	return e == nil && u.Host != "" && (u.Scheme == "http" || u.Scheme == "https") && u.User == nil && u.RawQuery == "" && u.Fragment == ""
}
func (i *Index) call(ctx context.Context, method, target, key string, body, out any) error {
	if !strings.HasPrefix(target, strings.TrimRight(i.EmbeddingURL, "/")+"/") && !strings.HasPrefix(target, strings.TrimRight(i.QdrantURL, "/")+"/") {
		return domain.ErrInvalid
	}
	raw, e := json.Marshal(body)
	if e != nil {
		return e
	}
	r, e := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(raw))
	if e != nil {
		return e
	}
	r.Header.Set("Content-Type", "application/json")
	if key != "" {
		r.Header.Set("Authorization", "Bearer "+key)
	}
	res, e := i.Client.Do(r)
	if e != nil {
		return errors.New("index dependency unavailable")
	}
	defer res.Body.Close()
	data, e := io.ReadAll(io.LimitReader(res.Body, 16<<20))
	if e != nil {
		return e
	}
	if res.StatusCode == 401 || res.StatusCode == 403 {
		return domain.ErrUnauthenticated
	}
	if res.StatusCode == 400 {
		return domain.ErrInvalid
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.New("index dependency unavailable")
	}
	if out != nil && json.Unmarshal(data, out) != nil {
		return errors.New("invalid index response")
	}
	return nil
}
func (i *Index) Embed(ctx context.Context, texts []string) ([][]float64, error) {
	if i.EmbeddingKey == "" || i.Dimension < 1 || !validBase(i.EmbeddingURL) {
		return nil, errors.New("embedding not configured")
	}
	all := [][]float64{}
	for start := 0; start < len(texts); start += 16 {
		end := start + 16
		if end > len(texts) {
			end = len(texts)
		}
		var out struct {
			Data []struct {
				Index     int       `json:"index"`
				Embedding []float64 `json:"embedding"`
			} `json:"data"`
		}
		if e := i.call(ctx, "POST", strings.TrimRight(i.EmbeddingURL, "/")+"/embeddings", i.EmbeddingKey, map[string]any{"model": i.EmbeddingModel, "input": texts[start:end]}, &out); e != nil {
			return nil, e
		}
		if len(out.Data) != end-start {
			return nil, errors.New("invalid embedding count")
		}
		batch := make([][]float64, end-start)
		for _, v := range out.Data {
			if v.Index < 0 || v.Index >= len(batch) || batch[v.Index] != nil || len(v.Embedding) != i.Dimension {
				return nil, errors.New("invalid embedding shape")
			}
			norm := 0.0
			for _, x := range v.Embedding {
				if math.IsNaN(x) || math.IsInf(x, 0) {
					return nil, errors.New("invalid embedding value")
				}
				norm += x * x
			}
			if norm == 0 {
				return nil, errors.New("invalid embedding magnitude")
			}
			batch[v.Index] = v.Embedding
		}
		all = append(all, batch...)
	}
	return all, nil
}
func (i *Index) ensure(ctx context.Context) error {
	target := strings.TrimRight(i.QdrantURL, "/") + "/collections/campusclaw"
	r, e := http.NewRequestWithContext(ctx, "GET", target, nil)
	if e != nil {
		return e
	}
	res, e := i.Client.Do(r)
	if e != nil {
		return errors.New("vector unavailable")
	}
	io.Copy(io.Discard, res.Body)
	res.Body.Close()
	if res.StatusCode == 200 {
		return nil
	}
	if res.StatusCode != 404 {
		return errors.New("vector unavailable")
	}
	e = i.call(ctx, "PUT", target, "", map[string]any{"vectors": map[string]any{"size": i.Dimension, "distance": "Cosine"}}, nil)
	if errors.Is(e, domain.ErrInvalid) {
		r, _ = http.NewRequestWithContext(ctx, "GET", target, nil)
		res, e = i.Client.Do(r)
		if e != nil {
			return e
		}
		defer res.Body.Close()
		if res.StatusCode == 200 {
			return nil
		}
	}
	return e
}
func (i *Index) Upsert(ctx context.Context, cid, mid string, chunks []domain.Chunk, vectors [][]float64) error {
	if len(chunks) != len(vectors) {
		return domain.ErrInvalid
	}
	if e := i.ensure(ctx); e != nil {
		return e
	}
	for start := 0; start < len(chunks); start += 32 {
		end := start + 32
		if end > len(chunks) {
			end = len(chunks)
		}
		points := []any{}
		for n := start; n < end; n++ {
			c := chunks[n]
			points = append(points, map[string]any{"id": c.ID, "vector": vectors[n], "payload": map[string]any{"class_id": cid, "material_id": mid, "version": c.Version, "chunk_id": c.ID}})
		}
		if e := i.call(ctx, "PUT", strings.TrimRight(i.QdrantURL, "/")+"/collections/campusclaw/points?wait=true", "", map[string]any{"points": points}, nil); e != nil {
			return e
		}
	}
	return nil
}
func (i *Index) Search(ctx context.Context, cid string, vector []float64, limit int, materials []string) ([]domain.ScoredID, error) {
	if len(vector) != i.Dimension {
		return nil, domain.ErrInvalid
	}
	// An empty collection is a valid empty knowledge base, not a swallowed outage.
	if e := i.ensure(ctx); e != nil {
		return nil, e
	}
	var out struct {
		Result []struct {
			ID    string  `json:"id"`
			Score float64 `json:"score"`
		} `json:"result"`
	}
	must := []any{map[string]any{"key": "class_id", "match": map[string]any{"value": cid}}}
	if materials != nil {
		if len(materials) == 0 {
			return []domain.ScoredID{}, nil
		}
		must = append(must, map[string]any{"key": "material_id", "match": map[string]any{"any": materials}})
	}
	e := i.call(ctx, "POST", strings.TrimRight(i.QdrantURL, "/")+"/collections/campusclaw/points/search", "", map[string]any{"vector": vector, "limit": limit, "score_threshold": 0.35, "filter": map[string]any{"must": must}}, &out)
	if e != nil {
		return nil, e
	}
	ids := []domain.ScoredID{}
	for _, v := range out.Result {
		if v.Score >= 0.35 {
			ids = append(ids, domain.ScoredID{ID: v.ID, Score: v.Score})
		}
	}
	return ids, nil
}
