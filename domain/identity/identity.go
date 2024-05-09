package identity

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrIncorrectPassword     = errors.New("incorrect password")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidSession        = errors.New("invalid session")
	ErrSessionTooOld         = errors.New("session too old")
	ErrIdentityNotFound      = errors.New("identity not found")
	ErrIdentityInvalidCursor = errors.New("identity invalid cursor")

	ErrIdentityWasDisabled = errors.New("identity was disabled")
)

type Service interface {
	WhoAmI(ctx context.Context, token string) (string, error)
}

type Identity struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	Password          string     `json:"password,omitempty"`
	PasswordChangedAt *time.Time `json:"password_changed_at"`
	Session           *Session   `json:"-"`
}

type Session struct {
	ID        string     `json:"id"`
	Token     *string    `json:"token"`
	ExpiresAt *time.Time `json:"expires_at"`
}
