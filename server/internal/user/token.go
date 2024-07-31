package user

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateTokens(jwtUser JWTUser) (*Tokens, error) {
	accessClaims := jwt.MapClaims{
		"id":       jwtUser.ID,
		"username": jwtUser.Username,
		"email":    jwtUser.Email,
		"exp":      time.Now().Add(15 * time.Minute).Unix(), // 15 minutes
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")

	signedAccessToken, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"id":       jwtUser.ID,
		"username": jwtUser.Username,
		"email":    jwtUser.Email,
		"exp":      time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	signedRefreshToken, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  AccessToken(signedAccessToken),
		RefreshToken: RefreshToken(signedRefreshToken),
	}, nil
}

// TODO: Move sessions types to session package

func ValidateAccessToken(token AccessToken) (*JWTUser, error) {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	parsedToken, err := jwt.ParseWithClaims(string(token), jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	user := &JWTUser{
		ID:       claims["id"].(int),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}

	return user, nil
}

func ValidateRefreshToken(token RefreshToken) (*JWTUser, error) {
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	parsedToken, err := jwt.ParseWithClaims(string(token), jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	user := &JWTUser{
		ID:       claims["id"].(int),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}

	return user, nil
}
