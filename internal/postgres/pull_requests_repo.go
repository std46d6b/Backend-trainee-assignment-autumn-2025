package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type PullRequestRepo struct {
	exec Execer
}

func NewPullRequestRepo(exec Execer) *PullRequestRepo {
	return &PullRequestRepo{exec: exec}
}

func (r *PullRequestRepo) InsertPullRequest(ctx context.Context, pullRequest domain.PullRequest) error {
	query := psql.
		Insert("pull_requests").
		Columns("pull_request_id", "pull_request_name", "author_id").
		Values(pullRequest.ID, pullRequest.Name, pullRequest.AuthorID)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error generating sql query: %w", err)
	}

	tag, err := r.exec.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.NewDomainError(domain.ErrCodePRExists, fmt.Sprintf("pull request %s already exists", pullRequest.ID))
		}
		return fmt.Errorf("error executing query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (r *PullRequestRepo) getPRBodyData(ctx context.Context, pullRequestID domain.PullRequestID) (domain.PullRequest, error) {
	query := psql.
		Select(
			"pull_request_id",
			"pull_request_name",
			"author_id",
			"status",
			"created_at",
			"merged_at",
		).
		From("pull_requests").
		Where("pull_request_id = ?", pullRequestID)

	sql, args, err := query.ToSql()
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("error generating sql query: %w", err)
	}

	var pr domain.PullRequest

	err = r.exec.QueryRow(ctx, sql, args...).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.PullRequest{}, domain.NewDomainError(domain.ErrCodeNotFound, fmt.Sprintf("pull request %s not found", pullRequestID))
		}

		return domain.PullRequest{}, fmt.Errorf("error executing query: %w", err)
	}

	return pr, nil
}

func (r *PullRequestRepo) addReviewersIDs(ctx context.Context, pullRequest domain.PullRequest) (domain.PullRequest, error) {
	query := psql.
		Select("user_id").
		From("assigned_reviewers").
		Where("pull_request_id = ?", pullRequest.ID)

	sql, args, err := query.ToSql()
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("error generating sql query: %w", err)
	}

	rows, err := r.exec.Query(ctx, sql, args...)
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("error executing query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var reviewerID domain.UserID
		err := rows.Scan(&reviewerID)
		if err != nil {
			return domain.PullRequest{}, fmt.Errorf("error scanning row: %w", err)
		}
		pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, reviewerID)
	}

	return pullRequest, nil
}

func (r *PullRequestRepo) GetByID(ctx context.Context, pullRequestID domain.PullRequestID) (domain.PullRequest, error) {
	pullRequestBase, err := r.getPRBodyData(ctx, pullRequestID)
	if err != nil {
		return domain.PullRequest{}, err
	}

	pullREquest, err := r.addReviewersIDs(ctx, pullRequestBase)
	if err != nil {
		return domain.PullRequest{}, err
	}

	return pullREquest, nil
}

func (r *PullRequestRepo) AddReviewer(ctx context.Context, pullRequestID domain.PullRequestID, reviewerID domain.UserID) error {
	query := psql.
		Insert("assigned_reviewers").
		Columns("pull_request_id", "user_id").
		Values(pullRequestID, reviewerID)

	sql, args, err := query.ToSql()
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

func (r *PullRequestRepo) RemoveReviewer(ctx context.Context, pullRequestID domain.PullRequestID, reviewerID domain.UserID) error {
	query := psql.
		Delete("assigned_reviewers").
		Where("pull_request_id = ?", pullRequestID).
		Where("user_id = ?", reviewerID)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error generating sql query: %w", err)
	}

	tag, err := r.exec.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.NewDomainError(domain.ErrCodeNotAssigned, "incomplete removal")
	}

	return nil
}

func (r *PullRequestRepo) MergePullRequest(ctx context.Context, pullRequest domain.PullRequest) error {
	query := psql.
		Update("pull_requests").
		Set("status", pullRequest.Status).
		Set("merged_at", pullRequest.MergedAt).
		Where("pull_request_id = ?", pullRequest.ID)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error generating sql query: %w", err)
	}

	tag, err := r.exec.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.NewDomainError(domain.ErrCodeNotFound, fmt.Sprintf("pull request %s not found", pullRequest.ID))
	}

	return nil
}
