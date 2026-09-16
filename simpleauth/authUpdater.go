package simpleauth

import (
	"errors"
	"maps"
	"net/http"
	"slices"
)

func UpdateAuthInfo(r *http.Request, newAuthInfo AuthInfo) error {
	if authProvider.WriteAuthInfo == nil {
		return errors.New("WriteAuthInfo function is not configured")
	}

	oldAuthInfo, ok := getAuthInfo(r)
	if !ok {
		return errors.New("auth info is not set")
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
