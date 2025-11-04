package core

type UserRepo interface {
	GetById(id int64) (User, error)
	CheckPassword(email, pass string) (User, error)
}

type SessionRepo interface {
	SetRefreshToken(userID int64, refreshToken string) error
	GetRefreshToken(userID int64) (string, error)
}

type TokenValidator interface {
	Sign(userID int64, email, role string) (string, error)
	Parse(tokenStr string) (map[string]any, error)
}
