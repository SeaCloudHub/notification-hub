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
	Login(ctx context.Context, email string, password string) (*Session, error)
	WhoAmI(ctx context.Context, token string) (*Identity, error)
	ChangePassword(ctx context.Context, id *Identity, oldPassword string, newPassword string) error
	SyncPasswordChangedAt(ctx context.Context, id *Identity) error

	// Admin APIs
	CreateIdentity(ctx context.Context, email string, password string) (*Identity, error)
	ListIdentities(ctx context.Context, pageToken string, pageSize int64) ([]Identity, string, error)
	ChangeIdentityState(ctx context.Context, id string, state string) (*Identity, error)
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
