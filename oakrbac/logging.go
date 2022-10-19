package oakrbac

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
)

type (
	// Logger captures an [Event] stream. It should serve as a primary observability mechanism for your domain logic application layer.
	Logger func(context.Context, *Event) error

	// An Event is an authorization record consumed by a [Logger].
	Event struct {
		IsAllowed bool
		Intent    *Intent
		Role      string
		Policy    Policy
		Error     error
	}
)

// PolicyName returns the name of the [Policy] function by using reflection. Use only for debugging.
func (e *Event) PolicyName() string {
	if e.Policy == nil {
		return "<nil>"
	}
	return runtime.FuncForPC(reflect.ValueOf(e.Policy).Pointer()).Name()
}

// PolicyNameFileLine returns the name of the [Policy] function, the name of the its file, and the line number by using reflection. Use only for debugging.
func (e *Event) PolicyNameFileLine() (name string, file string, line int) {
	if e.Policy == nil {
		return "<nil>", "<nil>", 0
	}
	definition := runtime.FuncForPC(reflect.ValueOf(e.Policy).Pointer())
	name = definition.Name()
	file, line = definition.FileLine(definition.Entry())
	return
}

func (e *Event) String() string {
	output := &strings.Builder{}
	if e.IsAllowed {
		output.WriteString("[ALLOW] to ")
	} else if e.Error == nil {
		output.WriteString("[DENY] to ")
	} else {
		output.WriteString("[FAILED] to ")
	}
	output.WriteString(e.Intent.String())

	output.WriteString(" by role ")
	output.WriteString(e.Role)

	if e.Policy != nil {
		output.WriteString(" (policy ")
		name, file, line := e.PolicyNameFileLine()
		output.WriteString(name)
		output.WriteString(" defined in ")
		output.WriteString(file)
		output.WriteString(", line ")
		output.WriteString(fmt.Sprintf("%d", line))
		output.WriteString(")")
	}

	if e.Error != nil {
		output.WriteString(" Reason for failure: ")
		output.WriteString(fmt.Sprintf("%s", e.Error))
	}

	return output.String()
}

// StandardLogger consumes [Event] records using the standard library logger. If uninitialized logger is provided, `log.Default()` is used instead. Do not use in production, because this logger relies on reflection.
func StandardLogger(l *log.Logger) Logger {
	if l == nil {
		l = log.Default()
	}
	return func(_ context.Context, e *Event) error {
		l.Println(e.String())
		return nil
	}
}

// MultiLogger sequentially combines a [Logger] list.
func MultiLogger(ls ...Logger) Logger {
	if len(ls) == 0 {
		panic("cannot combine an empty list of loggers")
	}
	for _, l := range ls {
		if l == nil {
			panic("logger list includes an uninitialized logger")
		}
	}

	return func(ctx context.Context, e *Event) (err error) {
		for _, l := range ls {
			if err = l(ctx, e); err != nil {
				return
			}
		}
		return nil
	}
}
