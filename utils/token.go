package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type TokenManager interface {
	GenerateAccessToken(userId string) string
	GenerateRefreshToken() string
	VerifyAccessToken(token string) (string, error)
	VerifyAccessTokenSignature(token string) (string, error)
}

type tokenManagerImpl struct {
	secret        string
	tokenDuration time.Duration
}

func (t *tokenManagerImpl) GenerateAccessToken(userId string) string {
	expiry := time.Now().Add(t.tokenDuration).Unix()
	textualExpiry := strconv.FormatInt(expiry, 10)
	message := userId + "." + textualExpiry
	signature := HmacSha256(message, t.secret)
	return message + "." + signature
}

func (t *tokenManagerImpl) VerifyAccessToken(token string) (string, error) {
	userId, expiry, err := t.getTokenPayload(token)
	if err != nil {
		return "", nil
	}

	if time.Now().Unix() > expiry {
		return "", errors.New("token expired")
	}

	return userId, nil
}

func (t *tokenManagerImpl) GenerateRefreshToken() string {
	return generateUserRefreshToken()
}

func (t *tokenManagerImpl) VerifyAccessTokenSignature(token string) (string, error) {
	userId, _, err := t.getTokenPayload(token)
	if err != nil {
		return "", nil
	}
	return userId, nil
}

func (t *tokenManagerImpl) getTokenPayload(token string) (string, int64, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", 0, errors.New("invalid token")
	}

	userId := parts[0]
	expiry, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return "", 0, errors.New("invalid token")
	}

	message := userId + "." + parts[1]
	expectedSignature := HmacSha256(message, t.secret)
	if parts[2] != expectedSignature {
		return "", 0, errors.New("wrong signature")
	}

	return userId, expiry, nil
}

func generateUserRefreshToken() string {
	return GenerateSecret(32)
}

func NewTokenManager(secret string, tokenDuration time.Duration) TokenManager {
	return &tokenManagerImpl{
		secret:        secret,
		tokenDuration: tokenDuration,
	}
}
