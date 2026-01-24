package database

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func buildInsert[T any](value T, table string) (string, []any, error) {
	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)

	var columns []string
	var args []any
	var placeholders []string

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
		if tag == "-" || tag == "" {
			continue
		}

		columns = append(columns, tag)
		args = append(args, v.Field(i).Interface())
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
	}

	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	return sql, args, nil
}

func InsertStruct[T any](pool *pgxpool.Pool, table string, value T) error {
	sql, args, err := buildInsert(value, table)
	if err != nil {
		return err
	}

	_, err = pool.Exec(context.Background(), sql, args...)
	return err
}
