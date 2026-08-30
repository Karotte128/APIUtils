package simpleauth

import (
	"net/http"
	"slices"

	"github.com/karotte128/karotteapi/v2/api"
	"github.com/karotte128/karottelib/config"
)

var authMiddleware = api.Middleware{
	Name:        "auth",
	Handler:     authHandler,
	Priority:    3,
	ForceEnable: false,
}

func authHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg, ok := api.GetMiddlewareConfig("auth")
		if !ok {
			http.Error(w, "Internal Server Error: No middleware config!", http.StatusInternalServerError)
			return
		}

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

		basePermissions, ok := config.GetNestedValue[[]string](cfg, "basePermissions")
		if !ok {
			http.Error(w, "Internal Server Error: Config value basePermissions not set!", http.StatusInternalServerError)
			return
		}

		combined := slices.Concat(authInfo.Permissions, basePermissions)
		slices.Sort(combined)
		authInfo.Permissions = slices.Compact(combined)

		setAuthInfo(r, &authInfo)

		next.ServeHTTP(w, r)
	})
}

func init() {
	api.RegisterMiddleware(authMiddleware)
}
