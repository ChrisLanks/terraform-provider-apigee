package apigee

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/zambien/go-apigee-edge"
	"log"
	"strings"
	"testing"
)

func TestAccTargetServer_Updated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTargetServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckTargetServerConfigRequired,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetServerExists("apigee_target_server.foo", "foo_target_server"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "name", "foo_target_server"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "host", "http://some.api.com"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "env", "test"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "port", "80"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "ssl_info.0.ssl_enabled", "0"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "ssl_info.0.client_auth_enabled", "0"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "ssl_info.0.ignore_validation_errors", "false"),
				),
			},
			resource.TestStep{
				Config: testAccCheckTargetServerConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetServerExists("apigee_target_server.foo", "foo_target_server_updated"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "name", "foo_target_server_updated"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "host", "https://some.updatedapi.com"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "env", "test"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "enabled", "false"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "port", "443"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "ssl_info.0.ssl_enabled", "1"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "ssl_info.0.client_auth_enabled", "1"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "ssl_info.0.key_store", "freetrial"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "ssl_info.0.trust_store", "freetrial"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "ssl_info.0.key_alias", "freetrial"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "ssl_info.0.ignore_validation_errors", "true"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "ssl_info.0.ciphers.0", "AES256"),
					resource.TestCheckResourceAttr(
						"apigee_target_server.foo", "ssl_info.0.protocols.0", "https"),
				),
			},
		},
	})
}

func testAccCheckTargetServerDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*apigee.EdgeClient)

	if err := targetServerDestroyHelper(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckTargetServerExists(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*apigee.EdgeClient)
		if err := targetServerExistsHelper(s, client, name); err != nil {
			log.Print("Error in testAccCheckTargetServerExists: %s", err)
			return err
		}
		return nil
	}
}

const testAccCheckTargetServerConfigRequired = `
resource "apigee_target_server" "foo" {
  name = "foo_target_server"
  host = "http://some.api.com"
  env = "test"
  enabled = true
  port = 80

  ssl_info {
    ssl_enabled = false
    client_auth_enabled = false
    ignore_validation_errors = false
  }
}
`

const testAccCheckTargetServerConfigUpdated = `
resource "apigee_target_server" "foo" {
  name = "foo_target_server_updated"
  host = "https://some.updatedapi.com"
  env = "test"
  enabled = false
  port = 443

  ssl_info {
    ssl_enabled = true
    client_auth_enabled = true
    key_store = "freetrial"
    trust_store = "freetrial"
    key_alias = "freetrial"
    ignore_validation_errors = true # don't really do this...
    ciphers = ["AES256"]
    protocols = ["https"]
  }
}
`

func targetServerDestroyHelper(s *terraform.State, client *apigee.EdgeClient) error {

	for _, r := range s.RootModule().Resources {
		id := r.Primary.ID

		if id == "" {
			return fmt.Errorf("No target server ID is set")
		}

		_, _, err := client.TargetServers.Get("foo_target_server", "test")

		if err != nil {
			if strings.Contains(err.Error(), "404 ") {
				return nil
			}
			return fmt.Errorf("Received an error retrieving target server  %+v\n", err)
		}
	}

	return fmt.Errorf("Target server still exists")
}

func targetServerExistsHelper(s *terraform.State, client *apigee.EdgeClient, name string) error {

	for _, r := range s.RootModule().Resources {
		id := r.Primary.ID

		if id == "" {
			return fmt.Errorf("No target server ID is set")
		}

		if targetServerData, _, err := client.TargetServers.Get(name, "test"); err != nil {
			return fmt.Errorf("Received an error retrieving target server  %+v\n", targetServerData)
		} else {
			log.Print("Created target server name: %s", targetServerData.Name)
		}

	}
	return nil
}
