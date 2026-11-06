package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

func parseModelStream(input io.Reader, emit func(string) error) error {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 4096), 1<<20)
	data := []string{}
	finished := false
	done := false
	dispatch := func() error {
		if len(data) == 0 {
			return nil
		}
		raw := strings.Join(data, "\n")
		data = nil
		if strings.TrimSpace(raw) == "[DONE]" {
			done = true
			return nil
		}
		var frame struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
				Finish *string `json:"finish_reason"`
			} `json:"choices"`
			Error json.RawMessage `json:"error"`
		}
		if json.Unmarshal([]byte(raw), &frame) != nil || len(frame.Error) > 0 {
			return errors.New("invalid stream")
		}
		for _, choice := range frame.Choices {
			if choice.Delta.Content != "" {
				if e := emit(choice.Delta.Content); e != nil {
					return e
				}
			}
			if choice.Finish != nil {
				if *choice.Finish != "stop" {
					return errors.New("incomplete generation")
				}
				finished = true
			}
		}
		return nil
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if e := dispatch(); e != nil {
				return e
			}
			if done {
				return nil
			}
			continue
		}
		if strings.HasPrefix(line, "data:") {
			data = append(data, strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
		}
	}
	if e := scanner.Err(); e != nil {
		return e
	}
	if e := dispatch(); e != nil {
		return e
	}
	if done || finished {
		return nil
	}
	return errors.New("stream ended without completion")
}
