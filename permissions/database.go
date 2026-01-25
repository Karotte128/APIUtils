package permissions

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karotte128/apiutils/database"
)

type DbAuthInfo struct {
	ApiKey      string   `db:"apikey"`
	Permissions []string `db:"permissions"`
}

func UpdatePermission(connpool *pgxpool.Pool, table string, data DbAuthInfo) error {
	err := database.UpdateStruct(connpool, table, data, "apikey = $3", data.ApiKey)

	return err
}

func SetPermission(connpool *pgxpool.Pool, table string, data DbAuthInfo) error {
	err := database.InsertStruct(connpool, table, data)

	return err
}

func GetPermission(connpool *pgxpool.Pool, table string, key string) (DbAuthInfo, error) {
	auth, err := database.SelectStruct[DbAuthInfo](connpool, table, "apikey = $1", key)

	return auth, err
}
