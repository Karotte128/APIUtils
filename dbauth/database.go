package dbauth

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karotte128/apiutils/v2/database"
)

type DbAuthInfo struct {
	ApiKey      string         `db:"apikey"`
	Permissions []string       `db:"permissions"`
	ValidUntil  time.Time      `db:"valid_until"`
	Info        map[string]any `db:"info"`
}

func UpdateAuth(connpool *pgxpool.Pool, table string, data DbAuthInfo) error {
	err := database.UpdateStruct(connpool, table, data, "apikey = $5", data.ApiKey)

	return err
}

func SetAuth(connpool *pgxpool.Pool, table string, data DbAuthInfo) error {
	err := database.InsertStruct(connpool, table, data)

	return err
}

func GetAuth(connpool *pgxpool.Pool, table string, key string) (DbAuthInfo, error) {
	auth, err := database.SelectStruct[DbAuthInfo](connpool, table, "apikey = $1", key)

	return auth, err
}
