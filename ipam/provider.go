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
	"os"

	//"github.com/google/martian/v3/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	googleoauth "golang.org/x/oauth2/google"
)

// Provider for Simple IPAM
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "URL where to connect with the IPAM Autopilot backend",
			},
			"credentials": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validateCredentials,
				ConflictsWith: []string{"access_token"},
			},

			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"credentials"},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ipam_ip_range":       ResourceIpRange(),
			"ipam_routing_domain": ResourceRoutingDomain(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigure(ctx, d, provider)
	}

	return provider
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, p *schema.Provider) (interface{}, diag.Diagnostics) {
	config := Config{}

	url := d.Get("url").(string)
	if url == "" {
		url = os.Getenv("IPAM_URL")
	}

	if url == "" {
		return nil, diag.Errorf("URL needed to access IPAM Autopilot")
	}

	config.Url = url

	// Check for primary credentials in config. Note that if neither is set, ADCs
	// will be used if available.
	if v, ok := d.GetOk("access_token"); ok {
		config.AccessToken = v.(string)
	}

	if v, ok := d.GetOk("credentials"); ok {
		config.Credentials = v.(string)
	}

	// only check environment variables if neither value was set in config- this
	// means config beats env var in all cases.
	if config.AccessToken == "" && config.Credentials == "" {
		config.Credentials = multiEnvSearch([]string{
			"GOOGLE_CREDENTIALS",
			"GOOGLE_CLOUD_KEYFILE_JSON",
			"GCLOUD_KEYFILE_JSON",
		})

		// Retrieve token if credentials were available in one of ENV variables
		if config.Credentials != "" {

			token, err := config.getToken()
			if err != nil {
				return "", diag.Errorf("unable to retrieve identityToken: %v", err)
			}
			config.AccessToken = token
			return config, nil
		}

		config.AccessToken = multiEnvSearch([]string{
			"GOOGLE_OAUTH_ACCESS_TOKEN",
			"GCP_IDENTITY_TOKEN",
		})

	}

	return config, nil
}

func validateCredentials(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		return
	}
	creds := v.(string)
	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(creds); err == nil {
		return
	}
	if _, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(creds)); err != nil {
		errors = append(errors,
			fmt.Errorf("JSON credentials are not valid: %s", err))
	}

	return
}
