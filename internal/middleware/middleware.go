package middleware

import (
	"context"
	"net/http"

	"github.com/mastermou8/goProject/internal/store"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type contextKey string

const UserContextKey = contextKey("user")

func setUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func getUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("user not found in context") // bad actor call
	}
	return user
}
