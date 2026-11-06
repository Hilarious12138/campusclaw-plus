package service

import (
	"campusclaw-plus/internal/domain"
	"context"
	"sort"
	"strings"
)

type Retrieval struct {
	Store domain.RetrievalStore
	Index domain.Index
}
type SearchResult struct {
	Hits    []domain.Chunk `json:"hits"`
	Message string         `json:"message"`
}

func (s Retrieval) Search(ctx context.Context, u domain.Identity, query, mode string, materials []string) (SearchResult, error) {
	out := SearchResult{Hits: []domain.Chunk{}}
	query = strings.TrimSpace(query)
	if query == "" || len([]rune(query)) > 1000 {
		return out, domain.ErrInvalid
	}
	if mode == "" {
		mode = "hybrid"
	}
	if mode != "keyword" && mode != "vector" && mode != "hybrid" {
		return out, domain.ErrInvalid
	}
	if materials != nil && len(materials) == 0 {
		out.Message = "资料中未找到相关内容"
		return out, nil
	}
	keyword := []domain.Chunk{}
	var e error
	if mode != "vector" {
		keyword, e = s.Store.Keyword(ctx, u, query, materials, 20)
		if e != nil {
			return out, e
		}
	}
	vector := []domain.Chunk{}
	if mode != "keyword" {
		vectors, e := s.Index.Embed(ctx, []string{query})
		if e != nil {
			return out, e
		}
		ids, e := s.Index.Search(ctx, u.ClassID, vectors[0], 100, materials)
		if e != nil {
			return out, e
		}
		keys := []string{}
		for _, v := range ids {
			keys = append(keys, v.ID)
		}
		valid, e := s.Store.ChunksByID(ctx, u, keys, materials)
		if e != nil {
			return out, e
		}
		byID := map[string]domain.Chunk{}
		for _, v := range valid {
			byID[v.ID] = v
		}
		for _, v := range ids {
			if c, ok := byID[v.ID]; ok {
				c.Score = v.Score
				vector = append(vector, c)
			}
		}
	}
	if mode == "keyword" {
		out.Hits = keyword
	} else if mode == "vector" {
		out.Hits = vector
	} else {
		scores := map[string]float64{}
		chunks := map[string]domain.Chunk{}
		for _, list := range [][]domain.Chunk{keyword, vector} {
			for rank, c := range list {
				scores[c.ID] += 1 / float64(60+rank+1)
				chunks[c.ID] = c
			}
		}
		for id, c := range chunks {
			c.Score = scores[id]
			out.Hits = append(out.Hits, c)
		}
		sort.Slice(out.Hits, func(i, j int) bool {
			if out.Hits[i].Score == out.Hits[j].Score {
				return out.Hits[i].ID < out.Hits[j].ID
			}
			return out.Hits[i].Score > out.Hits[j].Score
		})
	}
	if len(out.Hits) > 5 {
		out.Hits = out.Hits[:5]
	}
	used := 0
	for n := range out.Hits {
		rs := []rune(out.Hits[n].Text)
		remain := 6000 - used
		if len(rs) > remain {
			out.Hits[n].Text = string(rs[:remain])
			rs = rs[:remain]
		}
		used += len(rs)
	}
	if len(out.Hits) == 0 {
		out.Message = "资料中未找到相关内容"
	}
	return out, nil
}
