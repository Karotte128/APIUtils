package simpleauth

import (
	"strings"
	"time"
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
	return strings.HasPrefix(permission, prefix) &&
		strings.HasSuffix(permission, suffix)
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
func HasPermission(info AuthInfo, perm string) bool {
	if info.ValidUntil.Before(time.Now()) {
		return false
	}

	return checkPermission(info, perm)
}
