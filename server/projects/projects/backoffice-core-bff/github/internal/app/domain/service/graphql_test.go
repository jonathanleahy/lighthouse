package service

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/golang/mock/gomock"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/apierror"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/env"
	"github.com/pismo/backoffice-core-bff/internal/app/infrastructure/request"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewGraphqlService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gqlService := NewGraphqlService(&Resolver{})

	assert.NotNil(t, gqlService)
	assert.NotNil(t, gqlService.PlaygroundService)
	assert.NotNil(t, gqlService.GraphServer)
}

func TestHasAnyRole(t *testing.T) {
	ctx := context.WithValue(context.Background(), env.RequestContext, request.RequestContext{
		Roles: []string{"owner", "admin"},
	})
	gqlResolver := func(ctx context.Context) (res interface{}, err error) {
		return true, nil
	}
	type args struct {
		ctx   context.Context
		next  graphql.Resolver
		roles []string
	}
	tests := []struct {
		name    string
		args    args
		wantRes interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "should validate role",
			args:    args{ctx: ctx, next: gqlResolver, roles: []string{"owner"}},
			wantRes: true,
			wantErr: assert.NoError,
		},
		{
			name: "should invalidate role",
			args: args{
				ctx:   ctx,
				next:  gqlResolver,
				roles: []string{"other_role"},
			},
			wantRes: nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := HasAnyRole(tt.args.ctx, nil, tt.args.next, tt.args.roles)
			if !tt.wantErr(t, err, fmt.Sprintf("HasAnyRole(%v, %v, %v, %v)", tt.args.ctx, nil, tt.args.next, tt.args.roles)) {
				return
			}
			if err != nil {
				apiError := err.(apierror.ApiError)
				assert.Equal(t, http.StatusForbidden, apiError.GetHttpStatus())
				assert.Equal(t, "WCPBFF0003", apiError.GetMessage().ErrorCode)
				assert.Equal(t, "Access Denied", apiError.GetMessage().UserMessage)
			}
			assert.Equalf(t, tt.wantRes, gotRes, "HasAnyRole(%v, %v, %v, %v)", tt.args.ctx, nil, tt.args.next, tt.args.roles)
		})
	}
}

