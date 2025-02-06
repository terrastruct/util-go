// Package xhttp implements http helpers.
package xhttp

import (
	"context"
	"log"
	"net"
	"net/http"
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

func Serve(ctx context.Context, shutdownTimeout time.Duration, s *http.Server, l net.Listener) error {
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}

	done := make(chan error, 1)
	go func() {
		done <- s.Serve(l)
	}()

	select {
	case err := <-done:
		return err

	case <-ctx.Done():

		shutdownCtx := xcontext.WithoutCancel(ctx)
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, shutdownTimeout)
		defer cancel()
		shutdownErr := s.Shutdown(shutdownCtx)
		serveErr := <-done
		if serveErr != nil && serveErr != http.ErrServerClosed {
			return serveErr
		}
		
		return shutdownErr
	}
}
