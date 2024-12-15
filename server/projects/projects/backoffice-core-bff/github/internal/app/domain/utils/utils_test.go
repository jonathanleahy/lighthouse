package utils

import (
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/network"
	"reflect"
	"testing"
)

func TestGetQueryParameter(t *testing.T) {
	type args struct {
		page    *int
		perPage *int
	}
	tests := []struct {
		name string
		args args
		want []*network.QueryParameter
	}{
		{
			name: "Parameter_Page",
			args: args{page: new(int)},
			want: []*network.QueryParameter{{Name: "page", Value: new(int)}},
		},
		{
			name: "Parameter_perPage",
			args: args{perPage: new(int)},
			want: []*network.QueryParameter{{Name: "perPage", Value: new(int)}},
		},
		{
			name: "Parameter_Empty",
			want: []*network.QueryParameter(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetQueryParameter(tt.args.page, tt.args.perPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQueryParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}

