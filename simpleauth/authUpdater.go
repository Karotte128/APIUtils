package simpleauth

import (
	"errors"
	"maps"
	"net/http"
	"slices"

	"github.com/karotte128/karotteapi/v2/api"
	"github.com/karotte128/karottelib/config"
)

func UpdateAuthInfo(r *http.Request, newAuthInfo AuthInfo) error {
	if authProvider.WriteAuthInfo == nil {
		return errors.New("WriteAuthInfo function is not configured")
	}

	oldAuthInfo, ok := getAuthInfo(r)
	if !ok {
		return errors.New("auth info is not set")
	}

	cfg, ok := api.GetMiddlewareConfig("auth")
	if !ok {
		return errors.New("No middleware config!")
	}

	allowUpdatePermissions, ok := config.GetNestedValue[bool](cfg, "allowUpdatePermissions")
	if !ok {
		return errors.New("Config value not set!")
	}

	allowUpdateInfo, ok := config.GetNestedValue[bool](cfg, "allowUpdateInfo")
	if !ok {
		return errors.New("Config value not set!")
	}

	allowUpdateValidUntil, ok := config.GetNestedValue[bool](cfg, "allowUpdateValidUntil")
	if !ok {
		return errors.New("Config value not set!")
	}

	if !allowUpdatePermissions && !slices.Equal(oldAuthInfo.Permissions, newAuthInfo.Permissions) {
		return errors.New("Updating Permissions is not allowed!")
	}

	if !allowUpdateInfo && !maps.Equal(oldAuthInfo.Info, newAuthInfo.Info) {
		return errors.New("Updating Info is not allowed!")
	}

	if !allowUpdateValidUntil && !(oldAuthInfo.ValidUntil.Equal(newAuthInfo.ValidUntil)) {
		return errors.New("Updating ValidUntil is not allowed!")
	}

	return authProvider.WriteAuthInfo(newAuthInfo)
}
