package postgresrepo

import (
	"context"
	databasesql "database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
	pg "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/store/postgres"
)

type TeamRepo struct {
	exec    pg.Execer
	builder squirrel.StatementBuilderType
}

func NewTeamRepo(exec pg.Execer, builder squirrel.StatementBuilderType) *TeamRepo {
	return &TeamRepo{exec: exec, builder: builder}
}

func (r *TeamRepo) InsertTeam(ctx context.Context, teamName string) error {
	query := r.builder.
		Insert("teams").
		Columns("team_name").
		Values(teamName)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error generating sql query: %w", err)
	}

	_, err = r.exec.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.NewError(domain.ErrCodeTeamExists, fmt.Sprintf("team %s already exists", teamName))
		}
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}

func (r *TeamRepo) GetTeamWithMembers(ctx context.Context, teamName string) (domain.TeamUpsert, error) {
	query := r.builder.
		Select("u.user_id", "u.username", "u.is_active").
		From("teams t").
		LeftJoin("users u ON u.team_name = t.team_name").
		Where("t.team_name = ?", teamName)

	sql, args, err := query.ToSql()
	if err != nil {
		return domain.TeamUpsert{}, fmt.Errorf("error generating sql query: %w", err)
	}

	rows, err := r.exec.Query(ctx, sql, args...)
	if err != nil {
		return domain.TeamUpsert{}, fmt.Errorf("error executing query: %w", err)
	}

	defer rows.Close()

	var (
		members    []domain.TeamMember
		seenAnyRow bool
	)

	for rows.Next() {
		seenAnyRow = true

		var (
			memberUserID   databasesql.NullString
			memberUsername databasesql.NullString
			memberIsActive databasesql.NullBool
		)

		err = rows.Scan(&memberUserID, &memberUsername, &memberIsActive)
		if err != nil {
			return domain.TeamUpsert{}, fmt.Errorf("error scanning member: %w", err)
		}

		if memberUserID.Valid {
			members = append(members, domain.TeamMember{
				UserID:   memberUserID.String,
				Username: memberUsername.String,
				IsActive: memberIsActive.Bool,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return domain.TeamUpsert{}, fmt.Errorf("error scanning members: %w", err)
	}

	if !seenAnyRow {
		return domain.TeamUpsert{}, domain.NewError(domain.ErrCodeNotFound, fmt.Sprintf("team %s not found", teamName))
	}

	return domain.TeamUpsert{
		Name:    teamName,
		Members: members,
	}, nil
}
