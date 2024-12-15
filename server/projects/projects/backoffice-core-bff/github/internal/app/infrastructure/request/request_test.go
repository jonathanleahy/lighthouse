package request

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestContext_GetValueInCustomHeaders(t *testing.T) {
	rctx := RequestContext{
		CustomHeaders: map[string]string{
			HeaderXEmail: "pismo@pismo.io",
		},
	}

	result := rctx.GetValueInCustomHeaders(HeaderXEmail)

	assert.Equal(t, "pismo@pismo.io", result)
}

func TestRequestContext_HasAnyRole(t *testing.T) {
	type args struct {
		roles []string
	}
	roles := []string{"owner", "admin"}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should validate at least one role",
			args: args{[]string{"owner", "user"}},
			want: true,
		},
		{
			name: "should invalidate all roles",
			args: args{[]string{"user", "operator"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rctx := &RequestContext{
				Roles: roles,
			}
			if got := rctx.HasAnyRole(tt.args.roles...); got != tt.want {
				t.Errorf("HasAnyRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestContext_HasRole(t *testing.T) {
	type args struct {
		role string
	}
	roles := []string{"owner", "admin"}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should validate the role",
			args: args{"owner"},
			want: true,
		},
		{
			name: "should invalidate the role",
			args: args{"user"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rctx := &RequestContext{
				Roles: roles,
			}
			if got := rctx.HasRole(tt.args.role); got != tt.want {
				t.Errorf("HasRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

