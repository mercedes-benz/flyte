package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/flyteorg/flyte/flytestdlib/logger"
)

func BlanketAuthorization(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error) {

	identityContext := IdentityContextFromContext(ctx)
	if identityContext.IsEmpty() {
		return handler(ctx, req)
	}

	if !identityContext.Scopes().Has(ScopeAll) {
		logger.Debugf(ctx, "authenticated user doesn't have required scope")
		return nil, status.Errorf(codes.Unauthenticated, "authenticated user doesn't have required scope")
	}

	return handler(ctx, req)
}

// ExecutionUserIdentifierInterceptor injects identityContext.UserID() to identityContext.executionIdentity
func ExecutionUserIdentifierInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error) {
	identityContext := IdentityContextFromContext(ctx)
	identityContext = identityContext.WithExecutionUserIdentifier(identityContext.UserID())
	ctx = identityContext.WithContext(ctx)
	return handler(ctx, req)
}

// AuthorizationInterceptor returns an interceptor that only permits requests if user is permitted for current project
func AuthorizationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	projectID := inferProjectIDFromAdminRequest(info.FullMethod, req)
	if projectID != "" {
		identityContext := IdentityContextFromContext(ctx)
		projects, err := eligibleProjectsByIdentity(identityContext)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to authorize user %s due to error: %v", identityContext.UserID(), err)
			logger.Debugf(ctx, errMsg)
			return ctx, status.Errorf(codes.Unauthenticated, errMsg)
		}
		logger.Debugf(ctx, "Found eligible projects for user %s: %+q", identityContext.UserID(), projects)
		if _, ok := projects[projectID]; ok {
			return handler(ctx, req)
		}
		// User doesn't seem to be authorized to access or create current project
		errMsg := fmt.Sprintf("User %s not permitted to access project %s", identityContext.UserID(), projectID)
		logger.Debugf(ctx, errMsg)
		return ctx, status.Errorf(codes.Unauthenticated, errMsg)
	}
	return handler(ctx, req)
}
