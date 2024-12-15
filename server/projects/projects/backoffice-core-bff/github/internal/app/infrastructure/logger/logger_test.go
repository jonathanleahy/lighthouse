package logger

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDebug(t *testing.T) {
	type args struct {
		message string
		cid     string
		orgId   string
		fields  Fields
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "should log debug",
			args: args{
				message: "message", cid: "cid", orgId: "orgId",
			},
		},
		{
			name: "should log debug with fields",
			args: args{
				message: "message", cid: "cid", orgId: "orgId", fields: Fields{
					FieldApp: App,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debug(tt.args.message, tt.args.cid, tt.args.orgId, tt.args.fields)
		})
	}
}

func TestInfo(t *testing.T) {
	assert.NotPanics(t, func() {
		Info("message", "cid", "orgId", nil)
	})
}

func TestWarn(t *testing.T) {
	assert.NotPanics(t, func() {
		Warn("message", "cid", "orgId", nil)
	})
}

func TestError(t *testing.T) {
	assert.NotPanics(t, func() {
		Error("message", "cid", "orgId", nil)
	})
}

func TestErrorStackTrace(t *testing.T) {
	var err = os.Setenv("STACK_TRACE_ENABLED", "true")
	if err != nil {
		return
	}
	assert.NotPanics(t, func() {
		Error("message", "cid", "orgId", nil)
	})
}

func TestPanic(t *testing.T) {
	assert.Panics(t, func() {
		Panic("message", "cid", "orgId", nil)
	})
}

func TestInit(t *testing.T) {
	assert.NotPanics(t, func() {
		Init()
	})
}

func TestDefer(t *testing.T) {
	assert.NotPanics(t, func() {
		Defer()
	})
}

