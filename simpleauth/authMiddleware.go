package simpleauth

import (
	"net/http"

	"github.com/karotte128/karotteapi/v2/api"
)

var authMiddleware = api.Middleware{
	Name:        "auth",
	Handler:     authHandler,
	Priority:    3,
	ForceEnable: false,
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

func init() {
	api.RegisterMiddleware(authMiddleware)
}
