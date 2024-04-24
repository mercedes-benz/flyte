package auth

import (
	"context"
	"testing"

	authConfig "github.com/flyteorg/flyte/flyteadmin/auth/config"
	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyte/flyteidl/gen/pb-go/flyteidl/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/util/sets"
)

func TestBlanketAuthorization(t *testing.T) {
	t.Run("authenticated and authorized", func(t *testing.T) {
		allScopes := sets.NewString(ScopeAll)
		identityCtx := IdentityContext{
			audience: "aud",
			userID:   "uid",
			appID:    "appid",
			scopes:   &allScopes,
		}
		handlerCalled := false
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return nil, nil
		}
		ctx := context.WithValue(context.TODO(), ContextKeyIdentityContext, identityCtx)
		_, err := BlanketAuthorization(ctx, nil, nil, handler)
		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})
	t.Run("unauthenticated", func(t *testing.T) {
		handlerCalled := false
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return nil, nil
		}
		ctx := context.TODO()
		_, err := BlanketAuthorization(ctx, nil, nil, handler)
		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})
	t.Run("authenticated and not authorized", func(t *testing.T) {
		identityCtx := IdentityContext{
			audience: "aud",
			userID:   "uid",
			appID:    "appid",
		}
		handlerCalled := false
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return nil, nil
		}
		ctx := context.WithValue(context.TODO(), ContextKeyIdentityContext, identityCtx)
		_, err := BlanketAuthorization(ctx, nil, nil, handler)
		asStatus, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, asStatus.Code(), codes.Unauthenticated)
		assert.False(t, handlerCalled)
	})
}

func TestGetUserIdentityFromContext(t *testing.T) {
	identityContext := IdentityContext{
		userID: "yeee",
	}

	ctx := identityContext.WithContext(context.Background())

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		identityContext := IdentityContextFromContext(ctx)
		euid := identityContext.ExecutionIdentity()
		assert.Equal(t, euid, "yeee")
		return nil, nil
	}

	_, err := ExecutionUserIdentifierInterceptor(ctx, nil, nil, handler)
	assert.NoError(t, err)
}

func setAuthzConfig(authz authConfig.ProjectAuthorizationConfig) {
	getAuthConfigFunc = func() *authConfig.Config {
		return &authConfig.Config{ProjectAuthorization: authz}
	}
}

func TestAuthorizationInterceptor(t *testing.T) {
	claims := map[string]interface{}{"entitlement_group": []string{"r1", "r2"}}
	idCtx := IdentityContext{
		userID: "yeee",
		claims: &claims,
	}
	ctx := idCtx.WithContext(context.Background())
	request := admin.ExecutionCreateRequest{Project: "p1"}
	info := grpc.UnaryServerInfo{FullMethod: service.AdminService_CreateExecution_FullMethodName}

	t.Run("authorized for current project", func(t *testing.T) {
		handlerCalled := false
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return nil, nil
		}
		authz := authConfig.ProjectAuthorizationConfig{
			ProjectSets: map[string][]string{
				"r1": {"p1"},
			},
		}
		setAuthzConfig(authz)
		_, err := AuthorizationInterceptor(ctx, &request, &info, handler)
		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})

	t.Run("not authorized for current project", func(t *testing.T) {
		handlerCalled := false
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return nil, nil
		}
		request := admin.ExecutionCreateRequest{Project: "p1"}
		info := grpc.UnaryServerInfo{FullMethod: service.AdminService_CreateExecution_FullMethodName}

		authz := authConfig.ProjectAuthorizationConfig{ProjectSets: map[string][]string{
			"r3": {"p3"},
		}}
		setAuthzConfig(authz)
		_, err := AuthorizationInterceptor(ctx, &request, &info, handler)
		assert.ErrorContains(t, err, "rpc error: code = Unauthenticated desc = User yeee not permitted to access project p1")
		assert.False(t, handlerCalled)
	})

	t.Run("authorized by empty projectId", func(t *testing.T) {
		request := admin.ExecutionCreateRequest{}
		handlerCalled := false
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			handlerCalled = true
			return nil, nil
		}

		_, err := AuthorizationInterceptor(ctx, &request, &info, handler)
		assert.NoError(t, err)
		assert.True(t, handlerCalled)
	})
}
