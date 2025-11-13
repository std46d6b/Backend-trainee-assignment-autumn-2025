package domain

type UserID string

type User struct {
	ID       UserID
	Username string
	Team     TeamName
	IsActive bool
}
