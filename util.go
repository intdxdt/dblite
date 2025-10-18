package dblite

import (
	"fmt"
	"regexp"

	ref "github.com/intdxdt/goreflect"
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

func KeysToMap[K comparable, V any](keys []K, val V) Map[K, V] {
	var dict = make(map[K]V, len(keys))
	for _, k := range keys {
		dict[k] = val
	}
	return dict
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

func TableNameFromCreateSql(sql string) (string, error) {
	var matches = reCreateTable.FindStringSubmatch(sql)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("table name not found")
}
