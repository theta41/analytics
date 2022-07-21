package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/analytics/internal/http/util"
	"gitlab.com/g6834/team41/analytics/internal/ports"
)

func GetCheckAuthFunc(auth ports.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//login := "test@example.org" //stub login
			login := util.GetLoginFromCookie(r)
			tokens := util.GetTokensFromCookie(r)

			if login == "" || tokens.AccessToken == "" || tokens.RefreshToken == "" {
				http.Error(w, "{}", http.StatusForbidden)
				logrus.Error("empty login or one of tokens")
				return
			}

			// Validate access rights
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			newTokens, err := auth.Validate(ctx, login, tokens)
			if err != nil {
				http.Error(w, "{}", http.StatusForbidden)
				logrus.Error(err)
				return
			}

			util.PutLoginToCookie(w, login)
			util.PutTokensToCookie(w, newTokens)
			next.ServeHTTP(w, r)
		})
	}
}
