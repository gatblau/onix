package client

/*
   Onix Configuration Manager - HTTP Client
   Copyright (c) 2018-2021 by www.gatblau.org

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software distributed under
   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
   either express or implied.
   See the License for the specific language governing permissions and limitations under the License.

   Contributors to this project, hereby assign copyright in this code to the project,
   to be licensed under the same terms as the rest of the code.
*/
import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

// the authentication mode type
type AuthenticationMode int

// the different authentication modes available to the client
const (
	// no authentication used
	None AuthenticationMode = iota
	// use basic access authentication
	Basic
	// use OpenId Connect token
	OIDC
)

// client configuration information
type ClientConf struct {
	// the base URI for the service
	BaseURI string
	// disables TLS certificate verification
	InsecureSkipVerify bool
	// how to authenticate with the Web API
	AuthMode AuthenticationMode
	// the user username for Basic and OpenId user authentication
	Username string
	// the user password for Basic and OpenId user authentication
	Password string
	// the URI of the OpenId server token endpoint
	// used by the client to retrieve an OpenId token
	TokenURI string
	// the username to authenticate with the token service
	ClientId string
	// the password to authenticate with the token service
	AppSecret string
	// time out
	Timeout time.Duration
}

// sets the AuthMode from a passed-in string
func (cfg *ClientConf) SetAuthMode(authMode string) {
	switch strings.ToLower(authMode) {
	case "none":
		cfg.AuthMode = None
	case "basic":
		cfg.AuthMode = Basic
	case "oidc":
		cfg.AuthMode = OIDC
	default:
		log.Warn().Msgf("authMode value '%s' not recognised, defaulting to 'basic' authentication", authMode)
		cfg.AuthMode = Basic
	}
}

// gets the authentication token based on the authentication mode selected
func (cfg *ClientConf) getAuthToken() (string, error) {
	switch cfg.AuthMode {
	case Basic:
		return cfg.basicToken(cfg.Username, cfg.Password), nil
	case OIDC:
		token, err := cfg.bearerToken(cfg.TokenURI, cfg.ClientId, cfg.AppSecret, cfg.Username, cfg.Password)
		return token, err
	case None:
		return "", nil
	default:
		log.Warn().Msg("no authentication mode identified, defaulting to none")
		return "", nil
	}
}

// validates the client configuration
func checkConf(cfg *ClientConf) error {
	if len(cfg.BaseURI) == 0 {
		return errors.New("BaseURI is not defined")
	}
	// if the protocol is not specified, the add http as default
	// this is to avoid the server producing empty responses if no protocol is specified in the URI
	if !strings.HasPrefix(strings.ToLower(cfg.BaseURI), "http") {
		log.Warn().Msgf("no protocol defined for Onix URI '%s', 'http://' will be added to it", cfg.BaseURI)
		cfg.BaseURI = fmt.Sprintf("http://%s", cfg.BaseURI)
	}
	if cfg.AuthMode == Basic {
		if len(cfg.Username) == 0 {
			return errors.New("username is not defined")
		}
		if len(cfg.Password) == 0 {
			return errors.New("password is not defined")
		}
	}
	if cfg.AuthMode == OIDC {
		if len(cfg.Username) == 0 {
			return errors.New("username is not defined")
		}
		if len(cfg.Password) == 0 {
			return errors.New("password is not defined")
		}
		if len(cfg.TokenURI) == 0 {
			return errors.New("token URI is not defined")
		}
		if len(cfg.ClientId) == 0 {
			return errors.New("client Id is not defined")
		}
		if len(cfg.AppSecret) == 0 {
			return errors.New("app secret is not defined")
		}
	}
	// if timeout is zero, it never timeout so is not good
	if cfg.Timeout == 0*time.Second {
		// set a default timeout of 5 secs
		cfg.Timeout = 5 * time.Second
	}
	return nil
}

// creates a new Basic Authentication Token
func (cfg *ClientConf) basicToken(user string, pwd string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, pwd))))
}

// gets an OAuth 2 bearer token
func (cfg *ClientConf) bearerToken(tokenURI string, clientId string, secret string, user string, pwd string) (string, error) {
	// constructs a payload for the form POST to the authorisation server token URI
	// passing the type of grant, the username, password and scopes
	payload := strings.NewReader(
		fmt.Sprintf("grant_type=password&username=%s&password=%s&scope=openid%%20onix", user, pwd))

	// creates the http request
	req, err := http.NewRequest(POST, tokenURI, payload)

	// if any errors then return
	if err != nil {
		return "", errors.New("Failed to create request: " + err.Error())
	}

	// adds the relevant http headers
	req.Header.Add("accept", "application/json")                        // need a response in json format
	req.Header.Add("authorization", cfg.basicToken(clientId, secret))   // authenticates with client id and secret
	req.Header.Add("cache-control", "no-cache")                         // forces caches to submit the request to the origin server for validation before releasing a cached copy
	req.Header.Add("content-type", "application/x-www-form-urlencoded") // posting an http form

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.InsecureSkipVerify,
			},
		},
	}

	// submits the request to the authorisation server
	response, err := client.Do(req)

	// if any errors then return
	if err != nil {
		return "", errors.New("Failed when submitting request: " + err.Error())
	}
	if response.StatusCode != 200 {
		return "", errors.New("Failed to obtain access token: " + response.Status + " Hint: the client might be unauthorised.")
	}

	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	result := new(OAuthTokenResponse)

	// decodes the response
	err = json.NewDecoder(response.Body).Decode(result)

	// if any errors then return
	if err != nil {
		return "", err
	}

	// constructs and returns a bearer token
	return fmt.Sprintf("Bearer %s", result.AccessToken), nil
}
