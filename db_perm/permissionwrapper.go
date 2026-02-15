package db_perm

import (
	"log"

	"github.com/karotte128/apiutils/database"
)

func GetPermissionWrapper(key string) []string {
	info, err := GetPermission(database.ConnPool, "authentication", key)
	if err != nil {
		log.Println(err.Error())
	}

	return info.Permissions
}
