package dblite

import (
	"fmt"
	ref "github.com/intdxdt/goreflect"
	"regexp"
	"strings"
)

var reCreateTable = regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([^\s(]+)`)

type Pair[T, U any] struct {
	A T
	B U
}

type KeyValPair[T, U any] struct {
	Key T
	Val U
}

type Map[K comparable, V any] map[K]V

func (dict Map[K, V]) Keys() []K {
	var keys = make([]K, 0, len(dict))
	for k := range dict {
		keys = append(keys, k)
	}
	return keys
}

func (dict Map[K, V]) Values() []V {
	var vals = make([]V, 0, len(dict))
	for _, v := range dict {
		vals = append(vals, v)
	}
	return vals
}

func (dict Map[K, V]) Flatten() []KeyValPair[K, V] {
	var pairs = make([]KeyValPair[K, V], 0, len(dict))
	for k, v := range dict {
		pairs = append(pairs, KeyValPair[K, V]{k, v})
	}
	return pairs
}

func MapFn[T any, U any](slice []T, fn func(T) U) []U {
	var result = make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func KeysToMap[K comparable, V any](keys []K, val V) Map[K, V] {
	var dict = make(map[K]V, len(keys))
	for _, k := range keys {
		dict[k] = val
	}
	return dict
}

func ColumnNames(cols []string) string {
	return strings.Join(cols, ",")
}

func ColumnEqualPlaceholders(cols []string, dbType string) string {
	var columns = make([]string, len(cols))
	for i, col := range cols {
		switch dbType {
		case "postgres":
			columns[i] = fmt.Sprintf("%s = $%d", col, i+1)
		default:
			columns[i] = fmt.Sprintf("%s = ?", col)
		}
	}
	return strings.Join(columns, ", ")
}

func ColumnEqualParamAttributes(cols []string, dbType string) string {
	startIndex := len(cols) + 2
	updates := make([]string, len(cols))
	for i, col := range cols {
		switch dbType {
		case "postgres":
			updates[i] = fmt.Sprintf("%s = $%d", col, startIndex+i)
		default:
			updates[i] = fmt.Sprintf("%s = ?", col)
		}
	}
	return strings.Join(updates, ", ")
}

func ColumnPlaceholders(cols []string, dbType string) string {
	switch dbType {
	case "postgres":
		placeholders := make([]string, len(cols))
		for i := range cols {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}
		return strings.Join(placeholders, ",")
	default:
		return strings.TrimRight(strings.Repeat("?,", len(cols)), ",")
	}
}

func MapFnWithIndex[T any](in []T, fn func(int, T) string) []string {
	out := make([]string, len(in))
	for i, val := range in {
		out[i] = fn(i, val)
	}
	return out
}

func UpdatePlaceholders(cols []string, dbType string) string {
	return strings.Join(MapFnWithIndex(cols, func(i int, col string) string {
		switch dbType {
		case "postgres":
			return fmt.Sprintf(`%v=$%d`, col, i+1)
		default: // defaults to sqlite3
			return fmt.Sprintf(`%v=?`, col)
		}
	}), `,`)
}

func ColumnsByExclusion[T ITable[T]](model T, excludeColumns []string) ([]string, error) {
	var fields, err = ref.Fields(model)
	if err != nil {
		return nil, err
	}

	fields, _, err = ref.FilterFieldReferences(fields, model)
	if err != nil {
		return nil, err
	}

	var cols = make([]string, 0, len(fields))
	var dict = KeysToMap(excludeColumns, true)
	for _, field := range fields {
		if dict[field] {
			continue
		}
		cols = append(cols, field)
	}
	return cols, nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func TableNameFromCreateSql(sql string) (string, error) {
	var matches = reCreateTable.FindStringSubmatch(sql)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("table name not found")
}
