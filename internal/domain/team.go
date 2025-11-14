package domain

type TeamName string

type Team struct {
	Name    TeamName
	Members []UserID
}

type TeamMember struct {
	UserID   UserID
	Username string
	IsActive bool
}

type TeamUpsert struct {
	Name    TeamName
	Members []TeamMember
}
