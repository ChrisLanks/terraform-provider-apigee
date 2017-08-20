// File : provider.go
package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/zambien/go-apigee-edge"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_uri": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_BASE_URI", nil),
				Description: "Apigee Edge Base URI",
			},
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_USER", nil),
				Description: "Apigee User",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_PASSWORD", nil),
				Description: "Apigee User Password",
			},
			"org": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APIGEE_ORG", nil),
				Description: "Apigee Organization",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"apigee_api_proxy":            resourceApiProxy(),
			"apigee_api_proxy_deployment": resourceApiProxyDeployment(),
			"apigee_target_servers":       resourceTargetServers(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {

	auth := apigee.EdgeAuth{
		Username: d.Get("user").(string),
		Password: d.Get("password").(string),
	}

	opts := &apigee.EdgeClientOptions{
		MgmtUrl: d.Get("base_uri").(string),
		Org:     d.Get("org").(string),
		Auth:    &auth, Debug: false,
	}

	client, err := apigee.NewEdgeClient(opts)
	if err != nil {
		return client, err
	}

	return client, nil
}
