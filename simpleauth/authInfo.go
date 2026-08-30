package simpleauth

import (
	"net/http"
	"time"

	"github.com/karotte128/karotteapi/v2/api"
)

// AuthInfo is created by the auth middleware.
// It contains the authentication status and permissions of the request.
type AuthInfo struct {
	// ApiKey is the raw key sent by the user. Do not use this.
	ApiKey string

	// Permissions is the list of permissions the user has.
	Permissions []string

	// ValidUntil is a timestamp used to invalidate API keys.
	ValidUntil time.Time

	// Info can contain additional information about the key (like metadata)
	Info map[string]any
}

func setAuthInfo(r *http.Request, authInfo *AuthInfo) {
	api.SetRequestContext(r, "auth.info", authInfo)
}

func getAuthInfo(r *http.Request) (*AuthInfo, bool) {
	info, ok := api.GetRequestContext[*AuthInfo](r, "auth.info")
	if !ok {
		return nil, false
	}

	if info == nil {
		return nil, false
	}

	return info, true
}

func GetAuthInfo(r *http.Request) (AuthInfo, bool) {
	info, ok := api.GetRequestContext[*AuthInfo](r, "auth.info")
	if !ok {
		return AuthInfo{}, false
	}

	return *info, true
}
