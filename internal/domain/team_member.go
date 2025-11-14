package domain

type TeamMember struct {
	UserID   UserID
	Username string
	IsActive bool
}

type TeamUpsert struct {
	Name    TeamName
	Members []TeamMember
}
