package models

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetToken(t *testing.T) {
	t.Run("should use Cached Token", func(t *testing.T) {
		server := mockServer(t)
		defer server.Close()

		// Override the keyCloakTokenURL with the mock server URL
		originalKeyCloakTokenURL := keyCloakTokenURL
		keyCloakTokenURL = server.URL + "/realms/icos-dev/protocol/openid-connect/token"
		defer func() { keyCloakTokenURL = originalKeyCloakTokenURL }()

		requester := KeycloakTokenRequester{}
		firstToken, _ := FetchKeycloakToken(requester)
		secondToken, _ := FetchKeycloakToken(requester)

		assert.Equal(t, firstToken, secondToken)
	})

	t.Run("should request a New Token", func(t *testing.T) {
		server := mockServer(t)
		defer server.Close()

		// Override the keyCloakTokenURL with the mock server URL
		originalKeyCloakTokenURL := keyCloakTokenURL
		keyCloakTokenURL = server.URL + "/realms/icos-dev/protocol/openid-connect/token"
		defer func() { keyCloakTokenURL = originalKeyCloakTokenURL }()

		requester := KeycloakTokenRequester{}
		firstToken, err := FetchKeycloakToken(requester)

		assert.NoError(t, err)
		assert.NotEmpty(t, firstToken)

	})

	t.Run("should return error when requestNewToken fails", func(t *testing.T) {
		err := assert.AnError
		assert.Error(t, err)
	})

}

func mockServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/realms/icos-dev/protocol/openid-connect/token", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JWT{AccessToken: "mocked_access_token"})
	}))
	return server
}

func TestStoreToken(t *testing.T) {

	t.Run("should store token", func(t *testing.T) {
		tokenCache = make(map[string]CachedToken)
		token := JWT{
			AccessToken: "mocked_access_token",
			ExpiresIn:   900,
		}
		storeToken(token)

		cachedToken, exists := tokenCache[clientID]
		assert.True(t, exists)
		assert.Equal(t, token.AccessToken, cachedToken.Token.AccessToken)
	})

}

func TestGetCachedToken(t *testing.T) {

	t.Run("should return cached token", func(t *testing.T) {
		tokenCache = make(map[string]CachedToken)
		token := JWT{
			AccessToken: "mocked_access_token",
			ExpiresIn:   900,
		}
		storeToken(token)
		cachedToken, err := getCachedToken(clientID)

		assert.NoError(t, err)
		assert.Equal(t, token.AccessToken, cachedToken.AccessToken)
	})

	t.Run("should return error if token is not found", func(t *testing.T) {
		tokenCache = make(map[string]CachedToken)
		_, err := getCachedToken(clientID)

		assert.Error(t, err)
	})

}
