package repository

import (
	"context"

	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type TeamRepository interface {
	InsertTeam(ctx context.Context, teamName string) error
	GetTeamWithMembers(ctx context.Context, teamName string) (domain.TeamUpsert, error)
}
