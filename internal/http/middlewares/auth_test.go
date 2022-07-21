package middlewares

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/g6834/team41/analytics/internal/models"
)

const (
	loginCookieName   = "Login"
	accessCookieName  = "AccessToken"
	refreshCookieName = "RefreshToken"
)

type AuthServiceStub struct {
	inLogin   string
	inTokens  models.TokenPair
	outTokens models.TokenPair
	outErr    error
}

func (a *AuthServiceStub) Validate(ctx context.Context, login string, tokens models.TokenPair) (models.TokenPair, error) {
	a.inLogin = login
	a.inTokens = tokens
	return a.outTokens, a.outErr
}

func TestCheckAuth(t *testing.T) {
	var nextHandlerWasCalled int
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//tests may also goes here..

		nextHandlerWasCalled++
	})

	authStub := &AuthServiceStub{
		outTokens: models.TokenPair{
			AccessToken:  "123",
			RefreshToken: "456",
		},
	}
	handlerCheckAuth := GetCheckAuthFunc(authStub)(nextHandler)

	//prepare inbound cookies
	inRecorder := httptest.NewRecorder()
	http.SetCookie(inRecorder, &http.Cookie{Name: loginCookieName, Value: "test@example.org"})
	http.SetCookie(inRecorder, &http.Cookie{Name: accessCookieName, Value: "7-7"})
	http.SetCookie(inRecorder, &http.Cookie{Name: refreshCookieName, Value: "8-8"})

	inRequest := &http.Request{Header: http.Header{"Cookie": inRecorder.Header()["Set-Cookie"]}}

	//execute middleware
	middlewareRecorder := httptest.NewRecorder()
	handlerCheckAuth.ServeHTTP(middlewareRecorder, inRequest)

	//check http chain
	require.Greater(t, nextHandlerWasCalled, 0, "next handler wasn't called")

	//check auth.Validate call params
	expectLogin := "test@example.org"
	expectTokens := models.TokenPair{
		AccessToken:  "7-7",
		RefreshToken: "8-8",
	}
	require.Equal(t, expectLogin, authStub.inLogin)
	require.Equal(t, expectTokens, authStub.inTokens)

	//check outbound cookie
	helperRequest := &http.Request{Header: http.Header{"Cookie": middlewareRecorder.Header()["Set-Cookie"]}}

	cookieLogin, err := helperRequest.Cookie(loginCookieName)
	require.NoError(t, err, "missing login cookie")
	require.NotNil(t, cookieLogin, "missing login cookie")
	require.Equal(t, expectLogin, cookieLogin.Value, "wrong login cookei")

	cookieAccess, err := helperRequest.Cookie(accessCookieName)
	require.NoError(t, err, "missing access token")
	require.NotNil(t, cookieAccess, "missing access token")
	require.Equal(t, "123", cookieAccess.Value, "wrong access token")

	cookieRefresh, err := helperRequest.Cookie(refreshCookieName)
	require.NoError(t, err, "missing refresh token")
	require.NotNil(t, cookieRefresh, "missing refresh token")
	require.Equal(t, "456", cookieRefresh.Value, "wrong refresh token")
}

func TestCheckAuthErrorCookie(t *testing.T) {
	var nextHandlerWasCalled int
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//tests may also goes here..

		nextHandlerWasCalled++
	})

	authStub := &AuthServiceStub{}
	handlerCheckAuth := GetCheckAuthFunc(authStub)(nextHandler)

	inRequest, _ := http.NewRequest("GET", "http://testing", nil)

	//execute middleware
	middlewareRecorder := httptest.NewRecorder()
	handlerCheckAuth.ServeHTTP(middlewareRecorder, inRequest)

	//check http chain
	require.Equal(t, nextHandlerWasCalled, 0, "next handler was called")

	//check status code
	require.Equal(t, http.StatusForbidden, middlewareRecorder.Code)
}

func TestCheckAuthErrorValidate(t *testing.T) {
	var nextHandlerWasCalled int
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//tests may also goes here..

		nextHandlerWasCalled++
	})

	authStub := &AuthServiceStub{
		outErr: errors.New("validate error"),
	}
	handlerCheckAuth := GetCheckAuthFunc(authStub)(nextHandler)

	//prepare inbound cookies
	inRecorder := httptest.NewRecorder()
	http.SetCookie(inRecorder, &http.Cookie{Name: loginCookieName, Value: "test@example.org"})
	http.SetCookie(inRecorder, &http.Cookie{Name: accessCookieName, Value: "7-7"})
	http.SetCookie(inRecorder, &http.Cookie{Name: refreshCookieName, Value: "8-8"})

	inRequest := &http.Request{Header: http.Header{"Cookie": inRecorder.Header()["Set-Cookie"]}}

	//execute middleware
	middlewareRecorder := httptest.NewRecorder()
	handlerCheckAuth.ServeHTTP(middlewareRecorder, inRequest)

	//check http chain
	require.Equal(t, nextHandlerWasCalled, 0, "next handler was called")

	//check status code
	require.Equal(t, http.StatusForbidden, middlewareRecorder.Code)
}
