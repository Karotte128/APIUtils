package simpleauth

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/karotte128/karotteapi/v2/api"
	"github.com/karotte128/karottelib/config"
)

func permissionMatch(pattern string, permission string) bool {
	// Fast path
	if pattern == "*" {
		return true
	}

	// Enforce at most one wildcard
	if strings.Count(pattern, "*") > 1 {
		return false
	}

	// No wildcard → exact match
	if !strings.Contains(pattern, "*") {
		return pattern == permission
	}

	// Exactly one wildcard
	prefix, suffix, _ := strings.Cut(pattern, "*")
	return strings.HasPrefix(permission, prefix) && strings.HasSuffix(permission, suffix)
}

func checkPermission(info AuthInfo, requiredPerm string) bool {
	for _, perm := range info.Permissions {
		if permissionMatch(perm, requiredPerm) {
			return true
		}
	}

	return false
}

// This function checks if the AuthInfo of a request has the given permission.
func HasPermission(info AuthInfo, perm string) (bool, error) {
	if info.ValidUntil.Before(time.Now()) {
		return false, nil
	}

	cfg, ok := api.GetMiddlewareConfig("auth")
	if !ok {
		return false, errors.New("No middleware config!")
	}

	basePermissions, ok := config.GetNestedValue[[]string](cfg, "basePermissions")
	if !ok {
		return false, errors.New("Internal Server Error: Config value basePermissions not set!")
	}

	combined := slices.Concat(info.Permissions, basePermissions)
	slices.Sort(combined)
	info.Permissions = slices.Compact(combined)

	return checkPermission(info, perm), nil
}
