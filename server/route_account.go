package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
)

type getAccountPageRender interface {
	RenderAccountPage(w io.Writer, signedIn bool, characters []domain.CharacterDTO) error
	RenderAccountFragment(w io.Writer, characters []domain.CharacterDTO) error
}

type getAccountPageIdentity interface {
	IsSignedIn(r *http.Request) bool
	GetUserID(r *http.Request) *uuid.UUID
}

type getAccountPageCharacters interface {
	GetCharactersByUserID(userID uuid.UUID) ([]domain.CharacterDTO, error)
}

func getAccountHandler(render getAccountPageRender, identity getAccountPageIdentity, characters getAccountPageCharacters) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)
		signedIn := identity.IsSignedIn(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if !signedIn {
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusSeeOther)
		}

		userID := identity.GetUserID(r)
		if userID == nil {
			slog.Error("user id nil")
		}

		characters, err := characters.GetCharactersByUserID(*userID)
		if err != nil {
			slog.Error("get characters by user id", "err", err)
		}

		if isHtmx {
			err := render.RenderAccountFragment(w, characters)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err := render.RenderAccountPage(w, signedIn, characters)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})
}
