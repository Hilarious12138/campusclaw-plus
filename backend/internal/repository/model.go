package repository

import (
	"bytes"
	"campusclaw-plus/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type ChatModel struct{ BaseURL, APIKey, Model string }

func (m ChatModel) Stream(ctx context.Context, messages []domain.ChatMessage, emit func(string) error) error {
	if !validBase(m.BaseURL) || m.APIKey == "" || m.Model == "" {
		return errors.New("model configuration unavailable")
	}
	raw, _ := json.Marshal(map[string]any{"model": m.Model, "messages": messages, "stream": true, "temperature": 0, "max_tokens": 700})
	r, e := http.NewRequestWithContext(ctx, "POST", strings.TrimRight(m.BaseURL, "/")+"/chat/completions", bytes.NewReader(raw))
	if e != nil {
		return e
	}
	r.Header.Set("Authorization", "Bearer "+m.APIKey)
	r.Header.Set("Content-Type", "application/json")
	client := InternalClient(60 * time.Second)
	res, e := client.Do(r)
	if e != nil {
		return errors.New("model unavailable")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 || !strings.HasPrefix(res.Header.Get("Content-Type"), "text/event-stream") {
		return errors.New("model stream unavailable")
	}
	return parseModelStream(io.LimitReader(res.Body, 1<<20), emit)
}
