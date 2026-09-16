package simpleauth

import (
	"errors"
	"net/http"

	"github.com/karotte128/karotteapi/v3/api"
	"github.com/karotte128/karottelib/config"
)

var authMiddleware = api.Middleware{
	Name:        "auth",
	Handler:     authHandler,
	Priority:    3,
	ForceEnable: false,
	Startup:     startup,
	Shutdown:    nil,
}

func authHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authProvider.ReadAuthInfo == nil {
			http.Error(w, "Internal Server Error: Auth provider not set!", http.StatusInternalServerError)
			return
		}

		header := r.Header.Get("X-API-Key")

		var authInfo AuthInfo
		var err error

		authInfo, err = authProvider.ReadAuthInfo(header)

		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		setAuthInfo(r, &authInfo)

		next.ServeHTTP(w, r)
	})
}

var basePermissions []string
var allowUpdatePermissions bool
var allowUpdateInfo bool
var allowUpdateValidUntil bool

func startup() error {
	cfg, ok := api.GetMiddlewareConfig("auth")
	if !ok {
		return errors.New("No middleware config!")
	}

	basePerms, ok := config.GetNestedValue[[]string](cfg, "basePermissions")
	if !ok {
		return errors.New("Internal Server Error: Config value basePermissions not set!")
	}
	basePermissions = basePerms

	aup, ok := config.GetNestedValue[bool](cfg, "allowUpdatePermissions")
	if !ok {
		return errors.New("Config value not set!")
	}
	allowUpdatePermissions = aup

	aui, ok := config.GetNestedValue[bool](cfg, "allowUpdateInfo")
	if !ok {
		return errors.New("Config value not set!")
	}
	allowUpdateInfo = aui

	auvu, ok := config.GetNestedValue[bool](cfg, "allowUpdateValidUntil")
	if !ok {
		return errors.New("Config value not set!")
	}
	allowUpdateValidUntil = auvu

	return nil
}

func init() {
	api.RegisterMiddleware(authMiddleware)
}
