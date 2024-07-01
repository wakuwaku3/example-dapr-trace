package httpx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wakuwaku3/example-dapr-trace/server/lib/errorsx"
	"github.com/wakuwaku3/example-dapr-trace/server/lib/logx"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type (
	server struct {
		mux         *http.ServeMux
		middlewares []MiddlewareFunc
		option      *ServerOption
	}
	ServerOption struct {
		Port          string
		ReadTimeout   time.Duration
		WriteTimeout  time.Duration
		IdleTimeout   time.Duration
		CancelTimeout time.Duration
	}
	HandleFunc     func(w http.ResponseWriter, r *http.Request) error
	MiddlewareFunc func(w http.ResponseWriter, r *http.Request, next HandleFunc) error
)

var (
	ErrTimeout = errors.New("timeout")
)

func NewServer(option *ServerOption, middlewares ...MiddlewareFunc) *server {
	mux := http.NewServeMux()

	return &server{mux, middlewares, option}
}

func (s *server) HandleFunc(pattern string, handler HandleFunc, middlewares ...MiddlewareFunc) {
	s.mux.HandleFunc(pattern, otelhttp.WithRouteTag(pattern, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx, span, err := startSpan(r)
			if err != nil {
				s.handleError(w, r, err)
				return
			}
			defer span.End()

			r = r.WithContext(ctx)

			done := make(chan struct{})
			errorChan := make(chan error, 1)
			if err := func() error {
				go func(w http.ResponseWriter, r *http.Request) {
					defer func() {
						if err := errorsx.Recover(); err != nil {
							errorChan <- err
						}
					}()

					if err := s.handleWithMiddleware(append(s.middlewares, middlewares...), handler, w, r, 0); err != nil {
						errorChan <- err
						return
					}
					close(done)
				}(w, r)

				select {
				case err := <-errorChan:
					return err
				case <-done:
					break
				case <-r.Context().Done():
					return errorsx.Wrap(ErrTimeout)
				}

				return nil
			}(); err != nil {
				s.handleError(w, r, err)
			}
		})).ServeHTTP)
}

func (s *server) handleError(w http.ResponseWriter, r *http.Request, err error) {
	logger := logx.Provider.Get()
	logger.Error(r.Context(), err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *server) handleWithMiddleware(middlewares []MiddlewareFunc, handler HandleFunc, w http.ResponseWriter, r *http.Request, index int) error {
	middlewareCount := len(middlewares)
	if index+1 < middlewareCount {
		middleware := middlewares[index]
		return middleware(w, r, func(w http.ResponseWriter, r *http.Request) error {
			return s.handleWithMiddleware(middlewares, handler, w, r, index+1)
		})
	} else if index+1 == middlewareCount {
		middleware := middlewares[index]
		return middleware(w, r, handler)
	}

	return handler(w, r)
}

func (s *server) Serve(ctx context.Context) error {
	logger := logx.Provider.Get()
	srv := &http.Server{
		ReadTimeout:  s.option.ReadTimeout,
		WriteTimeout: s.option.WriteTimeout,
		IdleTimeout:  s.option.IdleTimeout,
		Handler:      s.mux,
		Addr:         fmt.Sprintf(":%s", s.option.Port),
	}

	logger.Info(ctx, fmt.Sprintf("server started on %s", s.option.Port))
	go srv.ListenAndServe()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), s.option.CancelTimeout)
	defer cancel()

	return srv.Shutdown(ctx)
}
