package database

import (
	_ "embed"
	"strings"
)

//go:embed schema.sql
var schemaSQL string

var schemaStatement = splitSQLStatement(schemaSQL)

func splitSQLStatement(src string) []string {
	parts := strings.Split(src, ";")
	out := make([]string, 0, len(parts))

	for _, part := range parts {
		stmt := strings.TrimSpace(part)
		if stmt == "" {continue}
		out = append(out, stmt+";")
	}

	return out
}
