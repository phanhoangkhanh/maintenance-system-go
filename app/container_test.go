package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestServeHTTPCanceledContext(t *testing.T) {
	listener := newTestListener(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)
	go func() { done <- serveHTTP(ctx, &http.Server{}, listener) }()
	if err := waitForServer(t, done); err != nil {
		t.Fatalf("canceled context should shut down normally: %v", err)
	}
}

func TestServeHTTPWaitsForActiveRequest(t *testing.T) {
	listener := newTestListener(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	started := make(chan struct{})
	release := make(chan struct{})
	var releaseOnce sync.Once
	finishRequest := func() { releaseOnce.Do(func() { close(release) }) }
	defer finishRequest()
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusNoContent)
	})}
	defer server.Close()
	done := make(chan error, 1)
	go func() { done <- serveHTTP(ctx, server, listener) }()

	transport := &http.Transport{}
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: transport, Timeout: 3 * time.Second}
	requestDone := make(chan error, 1)
	go func() {
		resp, err := client.Get("http://" + listener.Addr().String())
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode != http.StatusNoContent {
				err = errors.New("unexpected response status")
			}
		}
		requestDone <- err
	}()

	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("request did not reach the handler")
	}
	cancel()
	select {
	case err := <-done:
		t.Fatalf("server returned before the active request finished: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	finishRequest()
	if err := waitForServer(t, requestDone); err != nil {
		t.Fatalf("active request was interrupted: %v", err)
	}
	if err := waitForServer(t, done); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}
}

func TestServeHTTPReturnsListenerError(t *testing.T) {
	listener := newTestListener(t)
	listener.Close()
	done := make(chan error, 1)
	go func() { done <- serveHTTP(context.Background(), &http.Server{}, listener) }()
	if err := waitForServer(t, done); !errors.Is(err, net.ErrClosed) {
		t.Fatalf("expected the underlying listener error, got %v", err)
	}
}

func newTestListener(t *testing.T) net.Listener {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { listener.Close() })
	return listener
}

func waitForServer(t *testing.T, done <-chan error) error {
	t.Helper()
	select {
	case err := <-done:
		return err
	case <-time.After(3 * time.Second):
		t.Fatal("server operation did not finish")
		return nil
	}
}
