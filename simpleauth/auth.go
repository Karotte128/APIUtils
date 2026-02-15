package simpleauth

import (
	"context"
)

type PermissionProvider func(string) []string

var permProvider PermissionProvider

// This function checks if the AuthInfo of a request has the given permission.
func HasPermission(ctx context.Context, perm string) bool {
	info := GetAuthInfo(ctx)

	if info == nil {
		return false
	}

	return checkPermission(*info, perm)
}

// This sets up simpleauth using the permission provider.
func Setup(provider PermissionProvider) {
	permProvider = provider
}
