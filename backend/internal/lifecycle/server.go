package lifecycle

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

func Serve(ctx context.Context, s *http.Server, listener net.Listener, budget time.Duration) error {
	done := make(chan error, 1)
	go func() { done <- s.Serve(listener) }()
	select {
	case e := <-done:
		if errors.Is(e, http.ErrServerClosed) {
			return nil
		}
		return e
	case <-ctx.Done():
		stop, cancel := context.WithTimeout(context.Background(), budget)
		defer cancel()
		if e := s.Shutdown(stop); e != nil {
			_ = s.Close()
			return e
		}
		e := <-done
		if errors.Is(e, http.ErrServerClosed) {
			return nil
		}
		return e
	}
}
