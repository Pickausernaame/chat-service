package logger

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

func NewSentryClient(dsn, env, version string) (*sentry.Client, error) {
	return sentry.NewClient(sentry.ClientOptions{
		Dsn:         dsn,
		Release:     version,
		Environment: env,
		HTTPTransport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
			},
		},
	})
}

func NewSentryCore(opts Options) (zapcore.Core, error) {
	sentryFactory := func() (*sentry.Client, error) { return NewSentryClient(opts.sentryDSN, opts.env, opts.version) }
	sentryCfg := zapsentry.Configuration{
		Level: zapcore.WarnLevel,
	}
	core, err := zapsentry.NewCore(sentryCfg, sentryFactory)
	if err != nil {
		return nil, fmt.Errorf("sentry init error: %v", err)
	}
	return core, nil
}
