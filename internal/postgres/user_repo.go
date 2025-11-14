package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type UserRepo struct {
	exec Execer
}

func NewUserRepo(exec Execer) *UserRepo {
	return &UserRepo{exec: exec}
}

func (r *UserRepo) UpsertUser(ctx context.Context, user domain.User) error {
	query := psql.
		Insert("users").
		Columns("user_id", "username", "team_name", "is_active").
		Values(user.ID, user.Username, user.Team, user.IsActive)

	withUpdate := query.Suffix("ON CONFLICT (user_id) DO UPDATE SET username = EXCLUDED.username, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active")

	sql, args, err := withUpdate.ToSql()
	if err != nil {
		return fmt.Errorf("error generating sql query: %w", err)
	}

	tag, err := r.exec.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, userID domain.UserID) (domain.User, error) {
	query := psql.
		Select("user_id", "username", "team_name", "is_active").
		From("users").
		Where("user_id = ?", userID)

	sql, args, err := query.ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("error generating sql query: %w", err)
	}

	var user domain.User

	err = r.exec.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Username, &user.Team, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.NewDomainError(domain.ErrCodeNotFound, fmt.Sprintf("user %s not found", userID))
		}

		return domain.User{}, fmt.Errorf("error scanning user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) SetIsActive(ctx context.Context, userID domain.UserID, isActive bool) error {
	query := psql.
		Update("users").
		Set("is_active", isActive).
		Where("user_id = ?", string(userID))

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error generating sql query: %w", err)
	}

	tag, err := r.exec.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.NewDomainError(domain.ErrCodeNotFound, fmt.Sprintf("user %s not found", userID))
	}

	return nil
}

func (r *UserRepo) ListReviewPRs(ctx context.Context, userID domain.UserID) ([]domain.PullRequest, error) {
	query := psql.
		Select(
			"pr.pull_request_id",
			"pr.pull_request_name",
			"pr.author_id",
			"pr.status",
			"pr.created_at",
			"pr.merged_at",
		).
		From("assigned_reviewers ar").
		Where("ar.user_id = ?", userID).
		Join("pull_requests pr ON pr.pull_request_id = ar.pull_request_id").
		OrderBy("pr.created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error generating sql query: %w", err)
	}

	rows, err := r.exec.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	defer rows.Close()

	var pullRequests []domain.PullRequest

	for rows.Next() {
		var pr domain.PullRequest
		err = rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning pull request: %w", err)
		}

		pullRequests = append(pullRequests, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning pull requests: %w", err)
	}

	return pullRequests, nil
}
