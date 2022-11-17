// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ipam

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

// Config is used as a general provider config throughout the provider
type Config struct {
	Url         string
	AccessToken string
	Credentials string
	context     context.Context
}

// staticTokenSource is used to be able to identify static token sources without reflection.
type staticTokenSource struct {
	oauth2.TokenSource
}

const (
	userInfoScope = "https://www.googleapis.com/auth/userinfo.email"
)

func (c *Config) getToken() (string, error) {

	creds, err := c.GetCredentials([]string{userInfoScope})
	if err != nil {
		return "", fmt.Errorf("error calling getCredentials(): %v", err)
	}

	ctx := context.Background()

	targetAudience := c.Url // "http://ipam-autopilot.com"

	co := []option.ClientOption{}
	co = append(co, idtoken.WithCredentialsJSON(creds.JSON))

	idTokenSource, err := idtoken.NewTokenSource(ctx, targetAudience, co...)

	if err != nil {
		return "", fmt.Errorf("unable to retrieve TokenSource: %v", err)
	}
	idToken, err := idTokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve Token: %v", err)
	}

	return idToken.AccessToken, nil
}

// Get a set of credentials with a given scope (clientScopes) based on the Config object.
func (c *Config) GetCredentials(clientScopes []string) (googleoauth.Credentials, error) {
	if c.AccessToken != "" {
		contents, _, err := pathOrContents(c.AccessToken)
		if err != nil {
			return googleoauth.Credentials{}, fmt.Errorf("Error loading access token: %s", err)
		}

		token := &oauth2.Token{AccessToken: contents}

		log.Printf("[INFO] Authenticating using configured Google JSON 'access_token'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return googleoauth.Credentials{
			TokenSource: staticTokenSource{oauth2.StaticTokenSource(token)},
		}, nil
	}

	if c.Credentials != "" {
		contents, _, err := pathOrContents(c.Credentials)
		if err != nil {
			return googleoauth.Credentials{}, fmt.Errorf("error loading credentials: %s", err)
		}

		creds, err := googleoauth.CredentialsFromJSON(c.context, []byte(contents), clientScopes...)
		if err != nil {
			return googleoauth.Credentials{}, fmt.Errorf("unable to parse credentials from '%s': %s", contents, err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'credentials'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return *creds, nil
	}

	log.Printf("[INFO] Authenticating using DefaultClient...")
	log.Printf("[INFO]   -- Scopes: %s", clientScopes)
	defaultTS, err := googleoauth.DefaultTokenSource(context.Background(), clientScopes...)
	if err != nil {
		return googleoauth.Credentials{}, fmt.Errorf("Attempted to load application default credentials since neither `credentials` nor `access_token` was set in the provider block.  No credentials loaded. To use your gcloud credentials, run 'gcloud auth application-default login'.  Original error: %w", err)
	}

	return googleoauth.Credentials{
		TokenSource: defaultTS,
	}, err
}
