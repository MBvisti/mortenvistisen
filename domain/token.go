package domain

import "time"

type Token struct {
	scope     string
	expiresAt time.Time
	Hash      string
	plain     string
}

func NewToken(
	scope string,
	expiresAt time.Time,
	Hash string,
	plain string,
) Token {
	return Token{
		scope,
		expiresAt,
		Hash,
		plain,
	}
}

func (t *Token) GetPlainText() string {
	return t.plain
}

func (t *Token) GetExpirationTime() time.Time {
	return t.expiresAt
}

func (t *Token) GetScope() string {
	return t.scope
}
