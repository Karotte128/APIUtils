package db_perm

import (
	"log"

	"github.com/karotte128/apiutils/database"
)

// This is a database based permission provider for simpleauth.
func GetPermissionWrapper(key string) []string {
	if database.ConnPool == nil {
		return nil
	}

	info, err := GetPermission(database.ConnPool, "authentication", key)
	if err != nil {
		log.Println(err.Error())
	}

	return info.Permissions
}
