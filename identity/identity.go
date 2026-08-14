package identity

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const userContextKey contextKey = "userID"

type Identity struct{}

func (i *Identity) GetUserID(r *http.Request) *uuid.UUID {
	ctx := r.Context()
	userID, ok := ctx.Value(userContextKey).(uuid.UUID)
	if !ok {
		return nil
	}
	return &userID
}

func (i *Identity) IsSignedIn(r *http.Request) bool {
	ctx := r.Context()
	_, ok := ctx.Value(userContextKey).(uuid.UUID)
	if !ok {
		return false
	}
	return true
}

func EnrichContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userContextKey, userID)
}
