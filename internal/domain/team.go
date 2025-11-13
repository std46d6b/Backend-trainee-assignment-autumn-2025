package domain

type TeamName string

type Team struct {
	Name    TeamName
	Members []UserID
}
