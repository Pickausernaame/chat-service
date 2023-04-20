package logger

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Pickausernaame/chat-service/internal/validator"
)

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options
type Options struct {
	level     LogLevelOption
	env       string `option:"mandatory" validate:"required,oneof=dev stage prod"`
	sentryDSN string `validate:"http_url,omitempty"`
	version   string `validate:"semver,omitempty"`
}

//go:generate options-gen -out-filename=logger_lvl_option.gen.go -from-struct=LogLevelOption
type LogLevelOption struct {
	value string `option:"mandatory" validate:"required,oneof=debug info warn error"`
}

var globalLogLevel = zap.NewAtomicLevel()

func MustInit(opts Options) {
	if err := Init(opts); err != nil {
		panic(err)
	}
}

func Init(opts Options) error {
	if err := setLogLevel(opts.level); err != nil {
		return fmt.Errorf("set log level error: %v", err)
	}

	encoderConfig := zapcore.EncoderConfig{
		LevelKey:    "level",
		MessageKey:  "msg",
		NameKey:     "component",
		TimeKey:     "T",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	if opts.env == "stage" || opts.env == "prod" {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, globalLogLevel),
	}
	if opts.sentryDSN != "" {
		core, err := NewSentryCore(opts)
		if err != nil {
			return fmt.Errorf("sentry core creating error: %v", err)
		}
		cores = append(cores, core)
	}

	l := zap.New(zapcore.NewTee(cores...))

	if opts.version != "" {
		l = l.With(zap.String("version", opts.version))
	}

	if opts.env != "" {
		l = l.With(zap.String("env", opts.env))
	}

	zap.ReplaceGlobals(l)
	return nil
}

func setLogLevel(opt LogLevelOption) error {
	if err := validator.Validator.Struct(opt); err != nil {
		return fmt.Errorf("validation logger level error: %v", err)
	}

	lvl, err := zapcore.ParseLevel(opt.value)
	if err != nil {
		return fmt.Errorf("parse logger level error: %v", err)
	}

	globalLogLevel.SetLevel(lvl)
	return nil
}

func SetLogLevel(opt LogLevelOption) error {
	return setLogLevel(opt)
}

func LogLevel() string {
	return strings.ToUpper(globalLogLevel.String())
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
