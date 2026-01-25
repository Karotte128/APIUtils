package database

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func buildSelect[T any](table string, where string) (string, error) {
	var zero T

	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("T must be a struct")
	}

	cols := make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		tag := field.Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}

		cols = append(cols, tag)
	}

	if len(cols) == 0 {
		return "", fmt.Errorf("no selectable columns found")
	}

	sql := fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(cols, ", "), table)

	if where != "" {
		sql += " WHERE " + where
	}

	return sql, nil
}

func SelectStructs[T any](pool *pgxpool.Pool, table string, where string, whereArgs ...any) ([]T, error) {
	sql, err := buildSelect[T](table, where)
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(context.Background(), sql, whereArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

func SelectStruct[T any](pool *pgxpool.Pool, table string, where string, whereArgs ...any) (T, error) {
	values, err := SelectStructs[T](pool, table, where, whereArgs...)
	if err != nil {
		var zero T
		return zero, err
	}

	if len(values) == 0 {
		var zero T
		return zero, nil
	}

	return values[0], nil
}
