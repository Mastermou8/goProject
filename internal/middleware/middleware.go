package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mastermou8/goProject/internal/store"
	"github.com/mastermou8/goProject/internal/tokens"
	"github.com/mastermou8/goProject/internal/utils"
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

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("user not found in context") // bad actor call
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = setUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r) // go to next function in the chain
			return
		}
		headerParts := strings.Split(authHeader, " ") //Bearer <TOKEN>
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}
		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "token expired or invalid"})

			return
		}
		r = setUser(r, user)
		next.ServeHTTP(w, r)

	})

}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be authenticated to access this resource"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
