package postgres

import "github.com/Masterminds/squirrel"

var Psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
