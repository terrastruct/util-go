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

	serverClosed := make(chan struct{})
	var serverError error
	go func() {
		serverError = s.Serve(l)
		close(serverClosed)
	}()

	select {
	case <-serverClosed:
		return serverError
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(xcontext.WithoutCancel(ctx), shutdownTimeout)
		defer cancel()

		err := s.Shutdown(shutdownCtx)
		<-serverClosed // Wait for server to exit
		if err != nil {
			return err
		}
		return serverError
	}
}
