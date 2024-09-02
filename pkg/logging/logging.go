package logging

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger encoder configuration keys that can be passed via environment variables.
const (
	NameKeyEnvVar    = "LOGGING_NAME_KEY"
	MessageKeyEnvVar = "LOGGING_MESSAGE_KEY"
	LevelKeyEnvVar   = "LOGGING_LEVEL_KEY"
	TimeKeyEnvVar    = "LOGGING_TIME_KEY"
	CallerKeyEnvVar  = "LOGGING_CALLER_KEY"
)

// Default *zap.Logger configuration values.
const (
	DefaultLoggingOutput = "stdout"
	DefaultNameKey       = "logger"
	DefaultMessageKey    = "message"
	DefaultLevelKey      = "level"
	DefaultTimeKey       = "time"
	DefaultCallerKey     = "caller"
	DefaultTraceIDKey    = "trace_id"
	DefaultSpanIDKey     = "span_id"

	DefaultEncoding = JSONEncoding
	JSONEncoding    = "json"
	ConsoleEncoding = "console"
)

const (
	goVersionLogKey = "go_version"
	revisionLogKey  = "revision"
)

var (
	defaultLogger = zap.NewExample()
	mu            sync.Mutex
)

// Config represents logging configuration
type Config struct {
	// LoggingLevel is the logger logging level.
	// Can be one of: "debug", "info", "warn", "error", "dpanic", "panic", or "fatal". (default "info")
	LoggingLevel string

	// LoggingOutput is the logger output. Can be "stdout", "stderr" or a list of other files. (default "stdout")
	LoggingOutput []string

	// Encoding is the logger encoding format used when logging.
	// Can be one of: "json" or "console". (default "json")
	Encoding string

	// Core is the logger core. If not set, the default core will be used.
	// This option is useful for testing purposes.
	Core zapcore.Core
}

// Init initializes application logger
func Init(cfg Config) error {
	encoding := cfg.Encoding
	if encoding == "" {
		encoding = DefaultEncoding
	}

	zapConfig, err := newConfig(cfg.LoggingLevel, encoding, cfg.LoggingOutput)
	if err != nil {
		return err
	}

	logger, err := zapConfig.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		if cfg.Core != nil {
			return cfg.Core
		}
		return core
	}))
	if err != nil {
		return err
	}
	defer func() { _ = logger.Sync() }()

	SetLogger(logger)
	return nil
}

// SetLogger concurrently safe sets the default application logger.
// Avoid using this function directly, and prefer setting the logger using the common context instead,
// i.e: someContext.WithLogger(ctx, logger)
func SetLogger(logger *zap.Logger) {
	mu.Lock()
	defer mu.Unlock()
	defaultLogger = logger
}

// GetLogger concurrently safe gets the default application logger.
// Avoid using this function directly, and prefer getting the logger from the context instead,
// i.e: someContext.Logger(ctx)
func GetLogger() *zap.Logger {
	mu.Lock()
	defer mu.Unlock()
	return defaultLogger
}

func newConfig(level, encoding string, output []string) (zap.Config, error) {
	var logLevel zapcore.Level
	err := logLevel.Set(level)
	if err != nil {
		return zap.Config{}, err
	}
	if len(output) == 0 {
		output = []string{DefaultLoggingOutput}
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(logLevel)
	zapConfig.OutputPaths = output
	zapConfig.Encoding = encoding
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig.NameKey = newKey(NameKeyEnvVar, DefaultNameKey)
	zapConfig.EncoderConfig.MessageKey = newKey(MessageKeyEnvVar, DefaultMessageKey)
	zapConfig.EncoderConfig.LevelKey = newKey(LevelKeyEnvVar, DefaultLevelKey)
	zapConfig.EncoderConfig.TimeKey = newKey(TimeKeyEnvVar, DefaultTimeKey)
	zapConfig.EncoderConfig.CallerKey = newKey(CallerKeyEnvVar, DefaultCallerKey)

	return zapConfig, nil
}

func newKey(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
