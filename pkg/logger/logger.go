package logger

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"io"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// stackTracer the stackTracer interface
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Data type, used to pass to `ErrorWithData`.
type Data map[string]interface{}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func getCallerInfo() (pc uintptr, file string, line int, ok bool) {
	pc, file, line, ok = runtime.Caller(3)
	return
}

// SetLevel sets the standard entry level
func SetLevel(levelStr string) error {
	level, err := log.ParseLevel(levelStr)
	if err != nil {
		return err
	}

	log.SetLevel(level)
	return nil
}

// SetOutput sets the standard logger output.
func SetOutput(output io.Writer) {
	log.SetOutput(output)
}

// newEntry create new entry
func newEntry(file, line, function string) *log.Entry {
	return log.WithFields(log.Fields{
		"on":       fmt.Sprintf("%s:%s", file, line),
		"function": function,
	})
}

// An entry is the final or intermediate logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Debug, Info,
// Warn, Error, Fatal or Panic is called on it. These objects can be reused and
// passed around as much as you wish to avoid field duplication.
func entry() *log.Entry {
	pc, file, line, _ := getCallerInfo()
	fpc := uintptr(pc) - 1
	name := runtime.FuncForPC(fpc).Name()
	return newEntry(file, fmt.Sprintf("%d", line), funcName(name))
}

// entryFromError construct an entry from stackTracer error.
func entryFromError(err error) *log.Entry {
	if err, ok := err.(stackTracer); ok {
		f := err.StackTrace()[0]
		file := strings.Split(fmt.Sprintf("%+s", f), "\n\t")[1]
		function := fmt.Sprintf("%n", f)
		line := fmt.Sprintf("%d", f)
		return newEntry(file, line, function)
	}

	pc, file, line, _ := getCallerInfo()
	fpc := uintptr(pc) - 1
	name := runtime.FuncForPC(fpc).Name()
	return newEntry(file, fmt.Sprintf("%d", line), funcName(name))
}

//Debug the debug log
func Debug(ctx context.Context, msg ...interface{}) {
	fields := initFieldWithTraceID(ctx)
	entry().WithFields(fields).Debug(msg...)
}

//Info the info log
func Info(ctx context.Context, msg ...interface{}) {
	fields := initFieldWithTraceID(ctx)
	entry().WithFields(fields).Info(msg...)
}

//Warn the warn log
func Warn(ctx context.Context, msg ...interface{}) {
	fields := initFieldWithTraceID(ctx)
	entry().WithFields(fields).Warn(msg...)
}

//Error the error log
func Error(ctx context.Context, msg ...interface{}) {
	fields := initFieldWithTraceID(ctx)
	entry().WithFields(fields).Error(msg...)
}

//Fatal the fatal log
func Fatal(ctx context.Context, msg ...interface{}) {
	fields := initFieldWithTraceID(ctx)
	entry().WithFields(fields).Fatal(msg...)
}

// ErrorWithData the error log with data
func ErrorWithData(ctx context.Context, data Data, msg ...interface{}) {
	fields := initFieldWithTraceID(ctx)
	fields["data"] = data

	for _, m := range msg {
		if err, ok := m.(error); ok {
			entryFromError(err).WithFields(fields).Error(msg...)
			return
		}
	}

	log.WithFields(fields).Error(msg...)
}

func initFieldWithTraceID(ctx context.Context) log.Fields {
	fields := make(log.Fields)
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		dd := make(map[string]string)
		fields["dd"] = dd
	}

	return fields
}

// funcName removes the path prefix component of a function's name reported by func.Name().
func funcName(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
