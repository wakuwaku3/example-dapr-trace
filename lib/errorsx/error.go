package errorsx

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
)

type Body struct {
	Err   error  `json:"error,omitempty"`
	Stack string `json:"stack,omitempty"`
	Args  []any  `json:"args,omitempty"`
}

var (
	ErrInRecover = errors.New("error in recover")
)

func new(err error, args ...any) error {
	return &Body{
		Err:   err,
		Stack: stackTrace(),
		Args:  args,
	}
}

func (b *Body) Error() string {
	return b.Err.Error()
}

func Wrap(err error, args ...any) error {
	if err == nil {
		return nil
	}

	converted := &Body{}
	if ok := errors.As(err, &converted); ok {
		converted.Args = append(converted.Args, args...)
		return converted
	}
	return new(err, args...)
}

func Recover() error {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			return Wrap(err)
		} else {
			return Wrap(ErrInRecover, r)
		}
	}
	return nil
}

func stackTrace() string {
	const skipCount = 4
	const depth = 16

	pcs := make([]uintptr, depth)

	frameCount := runtime.Callers(skipCount, pcs)
	frames := runtime.CallersFrames(pcs[:frameCount])

	var sb strings.Builder
	sb.WriteString("goroutine [running]:\n")

	for {
		frame, more := frames.Next()

		sb.WriteString(frame.Function)
		sb.WriteString("\n\t")
		sb.WriteString(frame.File)
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(frame.Line))
		sb.WriteString("\n")

		if !more {
			break
		}
	}

	return sb.String()
}

func IsError(err error, target error) bool {
	converted := &Body{}
	if ok := errors.As(err, &converted); ok {
		return errors.Is(converted.Err, target)
	}
	return errors.Is(err, target)
}

func StackTrace(err error) string {
	converted := &Body{}
	if ok := errors.As(err, &converted); ok {
		return converted.Stack
	}
	return stackTrace()
}
