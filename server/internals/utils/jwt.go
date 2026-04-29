package utils

import (
	"errors"
	"os"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/golang-jwt/jwt/v5"
)

var jwt_secret_key = []byte(os.Getenv("JWT_SECRET"))
var ErrInvalidToken = errors.New("Invalid Token")
var ErrUnexpectedSigningMethod = errors.New("Unexpected Signing Method")

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokenPair(cfg *config.JWTConfig, userID uint, email, role string) (accessToken, refreshToken string, err error) {
	// Access Token
	accessClaims := &Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWTTokenExpiration)),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	accessTokenString, err := at.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := &Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWTRefreshTokenExpiration)),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := rt.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, err

}

func ValidateToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
