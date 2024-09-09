package logging

import (
	"log"
	"os"
	"sync"
	"time"

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
	mu            sync.RWMutex
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

	SamplingTick       time.Duration
	SamplingFirst      int
	SamplingThereafter int

	// Core is the logger core. If not set, the default core will be used.
	// This option is useful for testing purposes.
	Core zapcore.Core

	BufferingEnabled       bool
	BufferingSize          int
	BufferingFlushInterval time.Duration
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

	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}

	if cfg.BufferingEnabled {
		size := cfg.BufferingSize
		if size < 1 {
			size = 256 * 1024
		}

		flushInterval := cfg.BufferingFlushInterval
		if flushInterval < 1 {
			flushInterval = 30 * time.Second
		}

		ws, _, e := zap.Open(zapConfig.OutputPaths...)
		bufferedWriteSyncer := &zapcore.BufferedWriteSyncer{
			WS:            ws,
			Size:          size,
			FlushInterval: flushInterval,
		}
		if e != nil {
			return e
		}
		logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
				bufferedWriteSyncer,
				zapConfig.Level,
			)
		}))
	}
	if cfg.Core != nil {
		logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return cfg.Core
		}))
	}

	samplingTick := cfg.SamplingTick
	if samplingTick < 1 {
		samplingTick = time.Second
	}
	samplingFirst := cfg.SamplingFirst
	if samplingFirst < 1 {
		samplingFirst = 100
	}
	samplingThereafter := cfg.SamplingThereafter
	if samplingThereafter < 1 {
		samplingThereafter = 100
	}
	logger = zap.New(
		zapcore.NewSamplerWithOptions(
			logger.Core(),
			samplingTick,
			samplingFirst,
			samplingThereafter,
		),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	SetLogger(logger)
	return nil
}

// SetLogger concurrently safe sets the default application logger.
// Avoid using this function directly, and prefer setting the logger using the common context instead,
// i.e: someContext.WithLogger(ctx, logger)
func SetLogger(logger *zap.Logger) {
	mu.Lock()
	defaultLogger = logger
	mu.Unlock()
}

// GetLogger concurrently safe gets the default application logger.
// Avoid using this function directly, and prefer getting the logger from the context instead,
// i.e: someContext.Logger(ctx)
func GetLogger() *zap.Logger {
	mu.RLock()
	l := defaultLogger
	mu.RUnlock()
	return l
}

func Sync() {
	mu.Lock()
	defer mu.Unlock()
	_ = defaultLogger.Sync()
}

func HTTPErrorLogger() *log.Logger {
	return log.New(&httpErrorLogger{logger: GetLogger()}, "", 0)
}

type httpErrorLogger struct {
	logger *zap.Logger
}

func (l *httpErrorLogger) Write(p []byte) (n int, err error) {
	l.logger.Error(string(p))
	return len(p), nil
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
	zapConfig.Sampling = nil

	return zapConfig, nil
}

func newKey(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
