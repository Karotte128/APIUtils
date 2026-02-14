package simpleauth

import (
	"context"

	"github.com/karotte128/karotteapi"
	"github.com/karotte128/karotteapi/core"
)

// AuthInfo is created by the auth middleware.
// It contains the authentication status and permissions of the request.
type AuthInfo struct {
	// ApiKey is the raw key sent by the user. Do not use this.
	ApiKey string

	// Permissions is the list of permissions the user has.
	Permissions []string
}

func setAuthInfo(ctx context.Context, authInfo *AuthInfo) context.Context {
	reqCtx := karotteapi.RequestContext{
		Info:       authInfo,
		ContextKey: "auth",
	}

	return core.SetRequestContext(ctx, &reqCtx)
}

func getAuthInfo(ctx context.Context) *AuthInfo {
	reqCtx := core.GetRequestContext(ctx, "auth")

	return reqCtx.Info.(*AuthInfo)
}
