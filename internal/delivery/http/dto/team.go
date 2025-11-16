package dto

import "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamDTO struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

func TeamDomainToDTO(team domain.TeamUpsert) TeamDTO {
	members := make([]TeamMember, len(team.Members))

	for i, member := range team.Members {
		members[i] = TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	return TeamDTO{
		TeamName: team.Name,
		Members:  members,
	}
}

func TeamDTOToDomain(team TeamDTO) domain.TeamUpsert {
	members := make([]domain.TeamMember, len(team.Members))

	for i, member := range team.Members {
		members[i] = domain.TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	return domain.TeamUpsert{
		Name:    team.TeamName,
		Members: members,
	}
}
