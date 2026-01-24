package database

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func buildUpdate[T any](value T, table string, where string) (string, []any, error) {
	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)

	var args []any
	var assignments []string

	if t.Kind() == reflect.Pointer {
		v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("value must be a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		tag := field.Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}

		args = append(args, v.Field(i).Interface())
		assignments = append(assignments, fmt.Sprintf("%s = $%d", tag, len(args)))
	}

	sql := fmt.Sprintf(`UPDATE %s SET %s WHERE %s`, table, strings.Join(assignments, ", "), where)

	return sql, args, nil
}

func UpdateStruct[T any](pool *pgxpool.Pool, table string, value T, where string, whereArgs ...any) error {
	sql, args, err := buildUpdate(value, table, where)
	if err != nil {
		return err
	}

	args = append(args, whereArgs...)

	_, err = pool.Exec(context.Background(), sql, args...)
	return err
}
