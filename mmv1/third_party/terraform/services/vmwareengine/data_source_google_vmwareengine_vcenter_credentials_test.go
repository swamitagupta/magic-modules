package vmwareengine_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceVmwareengineVcenterCredentials_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"private_cloud_id": acctest.BootstrapVmwareenginePrivateCloud(t, "test-tl-pc", acctest.BootstrapVmwareengineNetwork(t, "test-nw"), true),
		"random_suffix":    acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleVmwareengineVcenterCredentialsConfig(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleVmwareengineVcenterCredentialsMeta("data.google_vmwareengine_vcenter_credentials.ds"),
				),
			},
		},
	})
}

func testAccCheckGoogleVmwareengineVcenterCredentialsConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_vmwareengine_vcenter_credentials" "ds" {
	parent =  "%{private_cloud_id}"
}
`, context)
}

func testAccCheckGoogleVmwareengineVcenterCredentialsMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Vcenter credentials data source: %s", n)
		}

		_, ok = rs.Primary.Attributes["username"]
		if !ok {
			return errors.New("can't find 'username' attribute")
		}

		_, ok = rs.Primary.Attributes["password"]
		if !ok {
			return errors.New("can't find 'password' attribute")
		}

		return nil
	}
}
