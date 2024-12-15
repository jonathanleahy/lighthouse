package logger

import (
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

type Fields map[string]interface{}

var (
	singleton *zap.Logger
	App       string
	Env       string
	Version   string
	once      sync.Once
)

const (
	FieldCid     = "Cid"
	FieldOrgId   = "OrgId"
	FieldApp     = "service.name"
	FieldVersion = "service.version"
	FieldEnv     = "service.env"
)

func Init() {
	once.Do(func() {
		singleton, _ = NewPismo()
		zap.ReplaceGlobals(singleton)
	})
}

func Defer() {
	if err := zap.L().Sync(); err != nil {
		return
	}
}

func init() {
	App = env.GetEnvWithDefaultAsString(env.AppName, env.DefaultAppName)
	Version = env.GetEnvWithDefaultAsString(env.Version, env.DefaultVersion)
	Env = env.GetEnvWithDefaultAsString(env.Env, env.DefaultEnv)
}

func makeFields(cid string, orgId string, fields Fields) []zapcore.Field {
	f := []zapcore.Field{
		zap.Any("Resource", map[string]interface{}{
			FieldApp:     App,
			FieldVersion: Version,
			FieldEnv:     Env,
		}),
	}

	if fields != nil {
		fields[FieldCid] = cid
		fields[FieldOrgId] = orgId
		f = append(f, zap.Any("Attributes", fields))
	}

	return f
}

func Debug(message string, cid string, orgId string, fields Fields) {
	zap.L().Debug(message, makeFields(cid, orgId, fields)...)
}

func Info(message string, cid string, orgId string, fields Fields) {
	zap.L().Info(message, makeFields(cid, orgId, fields)...)
}

func Warn(message string, cid string, orgId string, fields Fields) {
	zap.L().Warn(message, makeFields(cid, orgId, fields)...)
}

func Error(message string, cid string, orgId string, fields Fields) {
	if env.GetEnvWithDefaultAsBoolean(env.StackTraceEnabled, env.DefaultStackTraceEnabled) {
		zap.L().Error(message, makeFields(cid, orgId, fields)...)
	} else {
		zap.L().WithOptions(zap.AddStacktrace(zap.DPanicLevel)).Error(message, makeFields(cid, orgId, fields)...)
	}
}

func Panic(message string, cid string, orgId string, fields Fields) {
	zap.L().Panic(message, makeFields(cid, orgId, fields)...)
}

