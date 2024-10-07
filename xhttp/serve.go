// Package xhttp implements http helpers.
//
package xhttp

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"oss.terrastruct.com/util-go/xcontext"
)

func NewServer(log *log.Logger, h http.Handler) *http.Server {
	return &http.Server{
		MaxHeaderBytes: 1 << 18, // 262,144B
		ReadTimeout:    time.Minute,
		WriteTimeout:   time.Minute,
		IdleTimeout:    time.Hour,
		ErrorLog:       log,
		Handler:        http.MaxBytesHandler(h, 1<<20), // 1,048,576B
	}
}

type safeServer struct {
	*http.Server
	running int32
	mu      sync.Mutex
}

func newSafeServer(s *http.Server) *safeServer {
	return &safeServer{
		Server: s,
	}
}

func (s *safeServer) ListenAndServe(l net.Listener) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
		return errors.New("server is already running")
	}
	defer atomic.StoreInt32(&s.running, 0)

	return s.Serve(l)
}

func (s *safeServer) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if atomic.LoadInt32(&s.running) == 0 {
		return nil
	}

	return s.Server.Shutdown(ctx)
}

func Serve(ctx context.Context, shutdownTimeout time.Duration, s *http.Server, l net.Listener) error {
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}

	ss := newSafeServer(s)

	serverClosed := make(chan struct{})
	var serverError error
	go func() {
		serverError = ss.ListenAndServe(l)
		close(serverClosed)
	}()

	select {
	case <-serverClosed:
		return serverError
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(xcontext.WithoutCancel(ctx), shutdownTimeout)
		defer cancel()
		err := ss.Shutdown(shutdownCtx)
		<-serverClosed // Wait for server to exit
		if err != nil {
			return err
		}
		return serverError
	}
}
