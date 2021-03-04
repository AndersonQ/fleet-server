// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package logger

import (
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "ts",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

type zapStub struct {
}

func (z zapStub) Enabled(zapLevel zapcore.Level) bool {

	zeroLevel := log.Logger.GetLevel()

	switch zapLevel {
	case zapcore.DebugLevel:
		return zeroLevel == zerolog.DebugLevel
	case zapcore.InfoLevel:
		return zeroLevel <= zerolog.InfoLevel
	case zapcore.WarnLevel:
		return zeroLevel <= zerolog.WarnLevel
	case zapcore.ErrorLevel:
		return zeroLevel <= zerolog.ErrorLevel
	case zapcore.FatalLevel:
		return zeroLevel <= zerolog.FatalLevel
	case zapcore.DPanicLevel, zapcore.PanicLevel:
		return zeroLevel <= zerolog.PanicLevel
	}

	return true
}

func (z zapStub) Sync() error {
	return nil
}

func (z zapStub) Write(p []byte) (n int, err error) {
	log.Log().RawJSON("zap", p).Msg("")
	return 0, nil
}

func NewZapStub(selector string) *logp.Logger {

	wrapF := func(zapcore.Core) zapcore.Core {
		enc := zapcore.NewJSONEncoder(encoderConfig())
		stub := zapStub{}
		return zapcore.NewCore(enc, stub, stub)
	}

	return logp.NewLogger(selector, zap.WrapCore(wrapF))
}
