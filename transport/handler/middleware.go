package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"test/model"
)

type CtxKey string

const sessionCookie = "forum_session"

const UserContextKey = CtxKey("user")

func (h *Manager) SessMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookie)
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Println("sessmidllwareHandler:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		id, err := h.service.User.Authorizate(cookie.Value)
		if errors.Is(err, model.ErrNoCookie) {
			log.Println("sessmidllwareHandlerErrNoCookie:", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if errors.Is(err, model.ErrSessionIsExpired) {
			log.Println("sessmidllwareHandlerErrSessionIsExpired:", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Println("sessmidllwareHandler:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(context.Background(), UserContextKey, id)

		r = r.WithContext(ctx)

		next(w, r)
	})
}
