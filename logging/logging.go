package logging

import (
	"context"
	"github.com/google/uuid"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"os"
)

type Logger struct {
	*zap.Logger
	LevelId string
}

const (
	loggerKey        string = "logger"
	sugaredLoggerKey string = "sugaredLogger"
	// TODO: other consts?
	StatusCode      string = "statusCode"
	TaskTime        string = "taskTime"
	TraceId         string = "traceId"
	RequestPath     string = "requestPath"
	ClientKey       string = "headers.X-Client-Key"
	ApiKey          string = "apiKey"
	RequestId       string = "requestId"
	ParentRequestId string = "parentRequestId"
)

// TODO: USE SetLogger(ctx, log) AFTER!
func NewLogger(ctx context.Context, serviceName string) *Logger {
	return LoggerFactoryFor(serviceName)
}

func (logger *Logger) Fork(name string) (*Logger, string) {
	newId := uuid.New().String()
	return &Logger{
		logger.Logger.With(zap.String(name, newId)),
		newId,
	}, newId
}

func (logger *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		logger.Logger.With(fields...),
		logger.LevelId,
	}
}

func (logger *Logger) WithTraceId(ctx context.Context, traceId string) *Logger {
	return &Logger{
		logger.Logger.With(
			zap.String(TraceId, traceId), // TODO: string ok?
		),
		logger.LevelId,
	}
}
func (logger *Logger) WithRequestId(ctx context.Context, requestId string) *Logger {
	return &Logger{
		logger.Logger.With(
			zap.String(RequestId, requestId), // TODO: string ok?
		),
		logger.LevelId,
	}
}

func (logger *Logger) WithStringKVP(key, value string) *Logger {
	return &Logger{
		logger.Logger.With(zap.String(key, value)),
		logger.LevelId,
	}
}

func (logger *Logger) WithInt(key string, value int) *Logger {
	return &Logger{
		logger.Logger.With(zap.Int(key, value)),
		logger.LevelId,
	}
}

func (logger *Logger) WithField(field zap.Field) *Logger {
	return &Logger{
		logger.Logger.With(field),
		logger.LevelId,
	}
}

func LoggerFactoryFor(component string) *Logger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	newId := uuid.New().String()
	logger = logger.With(
		zap.String("component", component),
		zap.String("id", newId), // TODO: ????
	)
	return &Logger{logger, newId}
}

func SetLogger(parent context.Context, logger *Logger) context.Context {
	return context.WithValue(parent, loggerKey, logger)
}

func GetLogger(ctx context.Context) *Logger {
	v := ctx.Value(loggerKey)
	if logger, ok := v.(*Logger); ok {
		return logger
	}
	//Default is a Noop logger
	return &Logger{Logger: zap.NewNop()}
}

func SetSugaredLogger(parent context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(parent, sugaredLoggerKey, logger)
}

func GetSugaredLogger(ctx context.Context) *zap.SugaredLogger {
	v := ctx.Value(sugaredLoggerKey)
	if logger, ok := v.(*zap.SugaredLogger); ok {
		return logger
	}
	return GetLogger(ctx).Sugar()
}
