package dto

import "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"

type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func UserDomainToDTO(user domain.User) UserDTO {
	return UserDTO{
		UserID:   user.ID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func UserDTOToDomain(user UserDTO) domain.User {
	return domain.User{
		ID:       user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func UserDomainToDTOs(users []domain.User) []UserDTO {
	userDTOs := make([]UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = UserDomainToDTO(user)
	}
	return userDTOs
}
