package core

type UserService interface {
	GetById(userID int64) (User, error)
	GetStats() (map[string]any, error)
}

type AuthService interface {
	Login(email, password string) (accessToken, refreshToken string, userID int64, err error)
	RefreshTokens(userID int64, oldRefreshToken string) (newAccessToken, newRefreshToken string, err error)
}
