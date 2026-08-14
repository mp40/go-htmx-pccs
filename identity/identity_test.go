package identity

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestEnrichContextWithUserID(t *testing.T) {
	identity := &Identity{}
	id := uuid.New()
	ctx := EnrichContextWithUserID(context.Background(), id)
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)

	got := identity.GetUserID(r)
	if got == nil {
		t.Errorf("got nil, want uuid")
	}
	if *got != id {
		t.Errorf("got %v, want %v", *got, id)
	}
}
