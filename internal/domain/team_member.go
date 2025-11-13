package domain

type TeamMember struct {
	UserID   UserID
	Username string
	IsActive bool
}

type TeamUpset struct {
	Name    TeamName
	Members []TeamMember
}
