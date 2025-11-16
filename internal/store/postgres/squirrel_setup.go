package postgres

import "github.com/Masterminds/squirrel"

func NewStatementBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
