package domain

type Team struct {
	Name    string
	Members []string
}

type TeamMember struct {
	UserID   string
	Username string
	IsActive bool
}

type TeamUpsert struct {
	Name    string
	Members []TeamMember
}
