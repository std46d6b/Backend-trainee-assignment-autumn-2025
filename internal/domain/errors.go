package domain

type ErrorCode string

const (
	ErrCodeTeamExists  ErrorCode = "TEAM_EXISTS"
	ErrCodePRExists    ErrorCode = "PR_EXISTS"
	ErrCodePRMerged    ErrorCode = "PR_MERGED"
	ErrCodeNotAssigned ErrorCode = "NOT_ASSIGNED"
	ErrCodeNoCandidate ErrorCode = "NO_CANDIDATE"
	ErrCodeNotFound    ErrorCode = "NOT_FOUND"
)

type Error struct {
	Code    ErrorCode
	Message string
}

func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return string(e.Code) + ": " + e.Message
}
