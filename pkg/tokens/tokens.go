package tokens

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"log/slog"
	"time"
)

const (
	ScopeEmailVerification = "email_verification"
	ScopeResetPassword     = "password_reset"
)

type Manager struct {
	hasher hash.Hash
}

func NewManager(tokenSigningKey string) *Manager {
	h := hmac.New(sha256.New, []byte(tokenSigningKey))

	return &Manager{
		h,
	}
}

func (m *Manager) Hash(token string) string {
	m.hasher.Reset()
	m.hasher.Write([]byte(token))
	b := m.hasher.Sum(nil)

	return base64.URLEncoding.EncodeToString(b)
}

type GeneratedTokenDetails struct {
	PlainTextToken string
	HashedToken    string
}

func (m *Manager) GenerateToken() (GeneratedTokenDetails, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		slog.Error("could not hash token", "error", err)
		return GeneratedTokenDetails{}, err
	}

	plainText := base64.URLEncoding.EncodeToString(b)

	hashedToken := m.Hash(plainText)

	return GeneratedTokenDetails{
		plainText,
		hashedToken,
	}, nil
}

type Token struct {
	scope     string
	expiresAt time.Time
	Hash      string
	plainText string
}

func (t *Token) GetPlainText() string {
	return t.plainText
}

func (t *Token) GetExpirationTime() time.Time {
	return t.expiresAt
}

func (t *Token) GetScope() string {
	return t.scope
}

func CreateActivationToken(token, hashedToken string) Token {
	return Token{
		ScopeEmailVerification,
		time.Now().Add(72 * time.Hour),
		hashedToken,
		token,
	}
}

func CreateResetPasswordToken(token, hashedToken string) Token {
	return Token{
		ScopeResetPassword,
		time.Now().Add(24 * time.Hour),
		hashedToken,
		token,
	}
}
