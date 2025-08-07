package utils

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/otterly-id/otterly/backend/internal/api/models"
)

type Claims struct {
	ID   string          `json:"id"`
	Role models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey     []byte
	issuer        string
	audience      string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey []byte, issuer, audience string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		issuer:        issuer,
		audience:      audience,
		tokenDuration: tokenDuration,
	}
}

func (j *JWTManager) GenerateToken(userID, email string, role models.UserRole) (string, time.Duration, error) {
	now := time.Now()

	claims := Claims{
		ID:   userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{j.audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(j.tokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, j.tokenDuration, nil
}

func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if err := j.validateClaims(claims); err != nil {
		return nil, fmt.Errorf("invalid claims: %w", err)
	}

	return claims, nil
}

func (j *JWTManager) validateClaims(claims *Claims) error {
	if claims.Issuer != j.issuer {
		return errors.New("invalid issuer")
	}

	if !j.validateAudience(claims.Audience, j.audience) {
		return errors.New("invalid audience")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	if claims.NotBefore != nil && claims.NotBefore.After(time.Now()) {
		return errors.New("token not yet valid")
	}

	return nil
}

func (j *JWTManager) validateAudience(audience jwt.ClaimStrings, expected string) bool {
	if len(audience) == 0 {
		return false
	}

	return slices.Contains(audience, expected)
}
