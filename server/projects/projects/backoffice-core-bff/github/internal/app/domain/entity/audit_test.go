package entity

import (
	"reflect"
	"testing"
)

func TestNewAudit(t *testing.T) {
	tests := []struct {
		name string
		want Audit
	}{
		{
			name: "should create entity audit",
			want: Audit{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAudit(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAudit() = %v, want %v", got, tt.want)
			}
		})
	}
}

