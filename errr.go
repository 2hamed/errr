package errr

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var (
	MessageDelimiter = " -> "
)

type ErrorValue interface {
	error
	Log(lf logFunc, mf msgFunc)
	WithValue(key string, value interface{}) ErrorValue
	WithValues(map[string]interface{}) ErrorValue
	Msgf(format string, args ...interface{}) ErrorValue
	Wrap(error) ErrorValue
}

type logFunc func(string, interface{})
type msgFunc func(string)

func WithValue(key string, value interface{}) ErrorValue {
	e := errorValue{
		fields: make(map[string]interface{}),
	}
	_, file, lineno, ok := runtime.Caller(1)
	if ok {
		e.file = file
		e.lineno = lineno
	}
	return e.WithValue(key, value)
}

func WithValues(values map[string]interface{}) ErrorValue {
	e := errorValue{
		fields: values,
	}
	_, file, lineno, ok := runtime.Caller(1)
	if ok {
		e.file = file
		e.lineno = lineno
	}
	return e
}

type location struct {
	file   string
	lineno int
}
type errorValue struct {
	err        error
	fields     map[string]interface{}
	msg        string
	stacktrace []runtime.StackRecord
	location
}

func (e errorValue) Error() string {
	return fmt.Sprintf("%s %v %v", e.msg, e.err, e.fields)
}

func (e errorValue) Unwrap() error {
	return e.err
}

func (e errorValue) Log(lf logFunc, mf msgFunc) {
	for k, v := range e.fields {
		lf(k, v)
	}
	lf("file", e.file)
	lf("lineno", e.lineno)
	var builder strings.Builder
	msgf := func(msg string) {
		builder.WriteString(msg)
		builder.WriteString(MessageDelimiter)
	}
	var ev ErrorValue
	if errors.As(e.err, &ev) {
		ev.Log(lf, msgf)
	}
	builder.WriteString(e.msg)
	mf(builder.String())
}

// WithValue wraps an error inside an `errorValue` to provide a way to hold on to
// arbitrary key-value pairs to provide more details for upper levels of the call chain
// we use the same concept that is used in `context.WithValue`
func (e errorValue) WithValue(key string, value interface{}) ErrorValue {
	e.fields[key] = value
	return e
}

func (e errorValue) WithValues(fields map[string]interface{}) ErrorValue {
	var ev ErrorValue = e
	for k, v := range fields {
		ev = e.WithValue(k, v)
	}
	return ev
}

func (e errorValue) Msgf(format string, args ...interface{}) ErrorValue {
	e.msg = fmt.Sprintf(format, args...)
	return e
}
func (e errorValue) Wrap(err error) ErrorValue {
	e.err = err
	return e
}