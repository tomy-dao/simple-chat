package logger

import (
	"fmt"
	"local/model"
	"os"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// loggerImpl is a structured logger with OpenTelemetry integration
type loggerImpl struct {
	zerolog zerolog.Logger
}

var (
	// loggerInstance is the global logger instance
	loggerInstance *loggerImpl
)

func init() {
	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	
	// Use JSON format for Loki, or console for local development
	var zlog zerolog.Logger
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "json" || os.Getenv("ENV") == "production" || os.Getenv("ENV") == "docker" {
		// JSON format for structured logging (Loki)
		zlog = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		// Console format for local development
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
		zlog = zerolog.New(output).With().Timestamp().Logger()
	}
	
	loggerInstance = &loggerImpl{zerolog: zlog}
}

// fromContext creates a zerolog logger from RequestContext with trace information
func (l *loggerImpl) fromContext(reqCtx *model.RequestContext) zerolog.Logger {
	if reqCtx == nil {
		return l.zerolog
	}

	logger := l.zerolog

	// Add trace ID if available
	if traceID := reqCtx.TraceID(); traceID != "" {
		logger = logger.With().
			Str("trace_id", traceID).
			Logger()
	}

	// Add span ID if available
	if spanID := reqCtx.SpanID(); spanID != "" {
		logger = logger.With().
			Str("span_id", spanID).
			Logger()
	}

	// Add user ID if available
	if reqCtx.UserID > 0 {
		logger = logger.With().
			Uint("user_id", reqCtx.UserID).
			Logger()
	}

	// Add session ID if available
	if reqCtx.SessionID != "" {
		logger = logger.With().
			Str("session_id", reqCtx.SessionID).
			Logger()
	}

	return logger
}

// logEvent logs a message and adds it to OpenTelemetry span
func (l *loggerImpl) logEvent(reqCtx *model.RequestContext, level zerolog.Level, message string, fields ...map[string]interface{}) {
	logger := l.fromContext(reqCtx)
	
	// Create event
	event := logger.WithLevel(level)
	
	// Add fields if provided
	var fieldsMap map[string]interface{}
	if len(fields) > 0 && fields[0] != nil {
		fieldsMap = fields[0]
		for k, v := range fieldsMap {
			event = event.Interface(k, v)
		}
	}
	
	// Log message to console (zerolog)
	event.Msg(message)

	// Add to OpenTelemetry span root tracer - this will be exported to stdout
	// The event will appear in the trace output when span ends
	if reqCtx != nil && reqCtx.Span() != nil {
		attrs := []attribute.KeyValue{
			attribute.String("log.level", level.String()),
			attribute.String("log.message", message),
		}
		
		// Add fields as attributes
		for k, v := range fieldsMap {
			attrs = append(attrs, attribute.String(fmt.Sprintf("log.field.%s", k), fmt.Sprintf("%v", v)))
		}

		// Add event to span - will be exported by OpenTelemetry stdout exporter
		reqCtx.Span().AddEvent(message, trace.WithAttributes(attrs...))
	}
}

// Debug logs a debug message
func Debug(reqCtx *model.RequestContext, message string, fields ...map[string]interface{}) {
	loggerInstance.logEvent(reqCtx, zerolog.DebugLevel, message, fields...)
}

// Info logs an info message
func Info(reqCtx *model.RequestContext, message string, fields ...map[string]interface{}) {
	loggerInstance.logEvent(reqCtx, zerolog.InfoLevel, message, fields...)
}

// Warn logs a warning message
func Warn(reqCtx *model.RequestContext, message string, fields ...map[string]interface{}) {
	loggerInstance.logEvent(reqCtx, zerolog.WarnLevel, message, fields...)
}

// Error logs an error message
func Error(reqCtx *model.RequestContext, message string, err error, fields ...map[string]interface{}) {
	var fieldsMap map[string]interface{}
	if len(fields) > 0 && fields[0] != nil {
		fieldsMap = fields[0]
	} else {
		fieldsMap = make(map[string]interface{})
	}
	if err != nil {
		fieldsMap["error"] = err.Error()
	}
	loggerInstance.logEvent(reqCtx, zerolog.ErrorLevel, message, fieldsMap)
}

// Fatal logs a fatal message and exits
func Fatal(reqCtx *model.RequestContext, message string, fields ...map[string]interface{}) {
	logger := loggerInstance.fromContext(reqCtx)
	event := logger.Fatal()
	
	if len(fields) > 0 && fields[0] != nil {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	
	event.Msg(message)
}

// SetLevel sets the global log level
func SetLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}

