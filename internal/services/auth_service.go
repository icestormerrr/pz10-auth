package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/icestormerrr/pz10-auth/internal/core"
)

type AuthService struct {
	userRepo     core.UserRepo
	sessionRepo  core.SessionRepo
	tokenManager core.TokenManager
}

func NewAuthService(u core.UserRepo, s core.SessionRepo, t core.TokenManager) *AuthService {
	return &AuthService{
		userRepo:     u,
		sessionRepo:  s,
		tokenManager: t,
	}
}

func (s *AuthService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) Login(email, password string) (accessToken, refreshToken string, userID int64, err error) {
	u, err := s.userRepo.CheckPassword(email, password)
	if err != nil {
		return "", "", 0, errors.New("unauthorized")
	}

	accessToken, err = s.tokenManager.Sign(u.ID, u.Email, u.Role)
	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err = s.generateRefreshToken()
	if err != nil {
		return "", "", 0, err
	}

	if err := s.sessionRepo.SetRefreshToken(u.ID, refreshToken); err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, u.ID, nil
}

func (s *AuthService) RefreshTokens(userID int64, oldRefreshToken string) (newAccessToken, newRefreshToken string, err error) {
	storedToken, err := s.sessionRepo.GetRefreshToken(userID)
	if err != nil || storedToken == "" {
		return "", "", errors.New("refresh_token_not_found")
	}

	if storedToken != oldRefreshToken {
		return "", "", errors.New("invalid_refresh_token")
	}

	u, err := s.userRepo.GetById(userID)
	if err != nil {
		return "", "", errors.New("user_not_found")
	}

	newAccessToken, err = s.tokenManager.Sign(u.ID, u.Email, u.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err = s.generateRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := s.sessionRepo.SetRefreshToken(userID, newRefreshToken); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
