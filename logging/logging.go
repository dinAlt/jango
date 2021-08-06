package logging

import (
	"fmt"
	"reflect"
	"strings"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type NoopLogger struct{}

func (l NoopLogger) Printf(format string, v ...interface{}) {}

type StructLogger struct {
	path string
	l    Logger
}

func (l *StructLogger) Printf(method string, format string, v ...interface{}) {
	if _, ok := l.l.(NoopLogger); ok {
		return
	}

	l.l.Printf(l.path+"."+method+": "+format, v...)
}

func NewStructLogger(strct interface{}, logger Logger, trimPathPref string) *StructLogger {
	if logger == nil {
		panic("logger is nil")
	}

	t := reflect.TypeOf(strct)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("StructLogger can be created for struct only")
	}

	pkgPath := strings.TrimPrefix(t.PkgPath(), trimPathPref)
	path := fmt.Sprintf("%s.%s", pkgPath, t.Name())

	return &StructLogger{path, logger}
}
