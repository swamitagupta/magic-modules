package vmwareengine_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccVmwareengineDnsBindPermission_vmwareEngineDnsBindPermissionUpdate(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckVmwareengineDnsBindPermissionCustomDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: vmwareEngineDnsBindPermissionCreateConfigTemplate(),
			},
			{
				ResourceName:            "google_vmwareengine_dns_bind_permission.default-dns",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
			{
				Config: vmwareEngineDnsBindPermissionUpdateConfigTemplate(),
			},
			{
				ResourceName:            "google_vmwareengine_dns_bind_permission.default-dns",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name"},
			},
		},
	})
}

func vmwareEngineDnsBindPermissionCreateConfigTemplate() string {
	return acctest.Nprintf(`
resource "google_vmwareengine_dns_bind_permission" "default-dns" {
  principals {
    service_account = "test-service-account@vmwareengine-terraform-dev.iam.gserviceaccount.com"
  }
  principals {
    user = "devikamittal@google.com"
  }
}
`, nil)
}

func vmwareEngineDnsBindPermissionUpdateConfigTemplate() string {
	return acctest.Nprintf(`
resource "google_vmwareengine_dns_bind_permission" "default-dns" {
  principals {
    service_account = "test-service-account@vmwareengine-terraform-dev.iam.gserviceaccount.com"
  }
}
`, nil)
}

func testAccCheckVmwareengineDnsBindPermissionCustomDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vmwareengine_dns_bind_permission" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}
			log.Printf("[DEBUG] Checking destruction of DNS bind permission %s", name)
			config := acctest.GoogleProviderConfig(t)
			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{VmwareengineBasePath}}projects/{{project}}/locations/global/dnsBindPermission")
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] Generated URL: %s", url)
			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}
			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				log.Printf("[DEBUG] DNS bind permission response: %+v", res)
				_, ok := res["principals"]
				if ok {
					return fmt.Errorf("VmwareengineDnsBindPermission still exists at %s", url)
				}
				log.Printf("[DEBUG] DNS bind permission is deleted")
			}
		}
		return nil
	}
}
