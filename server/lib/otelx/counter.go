package otelx

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/metric"
)

type (
	counter struct {
		value map[string]metric.Int64Counter
		mu    sync.Mutex
	}
)

var counterValue = newCounter()

func Count(ctx context.Context, name string) error {
	c := counterValue
	counter, err := c.GetOrAdd(name)
	if err != nil {
		return err
	}
	counter.Add(ctx, 1)
	return nil
}

func newCounter() *counter {
	return &counter{
		value: make(map[string]metric.Int64Counter),
	}
}

func (c *counter) GetOrAdd(name string) (metric.Int64Counter, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.value[name]; !ok {
		m, err := Meter.Int64Counter(name)
		if err != nil {
			return nil, err
		}
		c.value[name] = m
	}
	return c.value[name], nil
}
