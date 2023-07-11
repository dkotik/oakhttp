package oakpolicy

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"golang.org/x/exp/slog"
)

type (
	// Policy returns [Allow] sentinel error if the session is permitted to interact with the context.
	// Policy returns [Deny] sentinel error to interrupt the matching loop.
	// Policy returns `nil` if it did not match, but another policy might match.
	//
	// In order to check a predicate assertion inside the policy, run a an anonymous interface type check on the [Resource].
	//
	//    predicated, ok := r.(interface{
	//      IsOwnedBySession(ctx) (bool, error)
	//    })
	Policy func(context.Context, Action, Resource) error
)

// Name returns the name of the [Policy] function by using reflection.
func (p Policy) Name() string {
	if p == nil {
		return "<no policy matched>"
	}
	return runtime.FuncForPC(reflect.ValueOf(p).Pointer()).Name()
}

// File returns the path to the file containing [Policy] function definition by using reflection.
func (p Policy) File() string {
	if p == nil {
		return "<nil>"
	}
	definition := runtime.FuncForPC(reflect.ValueOf(p).Pointer())
	file, _ := definition.FileLine(definition.Entry())
	return file
}

// Line returns the line number of the [Policy] function in its file by using reflection.
func (p Policy) Line() int {
	if p == nil {
		return 0
	}
	definition := runtime.FuncForPC(reflect.ValueOf(p).Pointer())
	_, line := definition.FileLine(definition.Entry())
	return line
}

func (p Policy) String() string {
	if p == nil {
		return "<no policy matched>"
	}
	definition := runtime.FuncForPC(reflect.ValueOf(p).Pointer())
	name := definition.Name()
	file, line := definition.FileLine(definition.Entry())

	output := strings.Builder{}
	output.WriteString("policy `")
	output.WriteString(name)
	output.WriteString("` defined in ")
	output.WriteString(file)
	output.WriteString(", line ")
	output.WriteString(fmt.Sprintf("%d", line))
	return output.String()
}

func (p Policy) Logger(l *slog.Logger) *slog.Logger {
	return l.With(slog.Any("policy", p))
}

func (p Policy) LogValue() slog.Value {
	if p == nil {
		return slog.StringValue("<nil policy>")
	}
	definition := runtime.FuncForPC(reflect.ValueOf(p).Pointer())
	name := definition.Name()
	file, line := definition.FileLine(definition.Entry())

	return slog.GroupValue(
		slog.String("name", name),
		slog.String("file", file),
		slog.Int("line", line),
	)
}

// AllowAll authorizes any action on any resource. Use cautiously.
func AllowAll(_ context.Context, _ Action, _ Resource) error {
	return Allow
}

// DenyAll denies authorization for any action on any resource.
func DenyAll(_ context.Context, _ Action, _ Resource) error {
	return Deny
}
