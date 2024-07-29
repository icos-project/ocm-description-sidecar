/*
  OCM-DESCRIPTOR-SIDECAR
  Copyright © 2022-2024 EVIDEN

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

  This work has received funding from the European Union's HORIZON research
  and innovation programme under grant agreement No. 101070177.
*/

package models

import (
	"encoding/json"
	"fmt"
	"icos/server/ocm-descriptor-sidecar/utils/logs"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

// This will allow us to easily mock the token request logic in tests
type TokenRequester interface {
	RequestNewToken() (JWT, error)
}

// KeycloakTokenRequester is a concrete implementation of the TokenRequester interface
type KeycloakTokenRequester struct{}
type JWT struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type CachedToken struct {
	Token      JWT
	ExpiryTime time.Time
}

var (
	keyCloakURL      = os.Getenv("KEYCLOAK_BASE_URL") // "https://keycloak.dev.icos.91.109.56.214.sslip.io"
	keyCloakRealm    = os.Getenv("KEYCLOAK_REALM")    // "icos-dev"
	keyCloakTokenURL = keyCloakURL + "/realms/" + keyCloakRealm + "/protocol/openid-connect/token"
	clientID         = os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret     = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	tokenCache       = make(map[string]CachedToken)
)

// FetchKeycloakToken fetches a token from the Keycloak server
func FetchKeycloakToken(requester TokenRequester) (JWT, error) {
	if cachedToken, err := getCachedToken(clientID); err == nil {
		logs.Logger.Println("Using Cached Token")
		return cachedToken, nil
	}

	logs.Logger.Println("Requesting New Token")
	token, err := requester.RequestNewToken()
	if err != nil {
		return JWT{}, err
	}

	storeToken(token)
	return token, nil
}

// RequestNewToken requests a new token from the Keycloak server
func (k KeycloakTokenRequester) RequestNewToken() (JWT, error) {
	reqToken, err := createTokenRequest()
	if err != nil {
		return JWT{}, err
	}

	resToken, err := sendTokenRequest(reqToken)
	if err != nil {
		return JWT{}, err
	}
	defer resToken.Body.Close()

	token, err := parseTokenResponse(resToken)
	if err != nil {
		return JWT{}, err
	}

	return token, nil
}

// createTokenRequest creates the token request
func createTokenRequest() (*http.Request, error) {
	reqTokenBody := url.Values{}
	reqTokenBody.Set("client_id", clientID)
	reqTokenBody.Set("grant_type", "client_credentials")
	reqTokenBody.Set("client_secret", clientSecret)

	logs.Logger.Println("Request Token Body is: ")
	logs.Logger.Println(reqTokenBody)

	reqToken, err := http.NewRequest("POST", keyCloakTokenURL, strings.NewReader(reqTokenBody.Encode()))
	if err != nil {
		logs.Logger.Println("ERROR " + err.Error())
		return nil, err
	}

	reqToken.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return reqToken, nil
}

// sendTokenRequest sends the token request to the server
func sendTokenRequest(reqToken *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resToken, err := client.Do(reqToken)
	if err != nil {
		logs.Logger.Println("ERROR " + err.Error())
		return nil, err
	}

	logs.Logger.Println("New Token Response Info: ")
	b, err := httputil.DumpResponse(resToken, true)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))

	return resToken, nil
}

// parseTokenResponse parses the token response from the server
func parseTokenResponse(resToken *http.Response) (JWT, error) {
	var token JWT

	tokenBody, err := io.ReadAll(resToken.Body)
	if err != nil {
		logs.Logger.Println("ERROR " + err.Error())
		return JWT{}, err
	}

	err = json.Unmarshal(tokenBody, &token)
	if err != nil {
		logs.Logger.Println("ERROR " + err.Error())
		return JWT{}, err
	}

	return token, nil
}

// getCachedToken fetches the token from the cache
func getCachedToken(clientID string) (JWT, error) {
	if cachedToken, ok := tokenCache[clientID]; ok {
		if time.Now().Before(cachedToken.ExpiryTime) {
			return cachedToken.Token, nil
		}
		delete(tokenCache, clientID)
	}
	return JWT{}, fmt.Errorf("token not found in cache")
}

// StoreToken stores the token in the cache
func storeToken(token JWT) {
	cachedToken := CachedToken{
		Token:      token,
		ExpiryTime: time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
	}
	tokenCache[clientID] = cachedToken
}
