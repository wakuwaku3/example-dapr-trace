package order

import (
	"errors"
	"io"
	"net/http"

	"github.com/wakuwaku3/example-dapr-trace/lib/errorsx"
	"github.com/wakuwaku3/example-dapr-trace/lib/logx"
)

func Get(w http.ResponseWriter, r *http.Request) error {
	logger := logx.Provider.Get()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(r.Context(), errorsx.Wrap(err))
	}

	logger.Info(r.Context(), "get order", logx.KeyValue{
		Key:   "data",
		Value: string(data),
	})

	logger.Error(r.Context(), errorsx.Wrap(errors.New("example error")))

	_, err = w.Write(data)
	if err != nil {
		logger.Error(r.Context(), errorsx.Wrap(err))
	}

	return nil
}
