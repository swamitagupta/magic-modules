package vmwareengine_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceVmwareEngineSubnet_managementSubnet(t *testing.T) {
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
				Config: testAccDataSourceVmwareEngineSubnetConfig(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleVmwareengineSubnetMeta("data.google_vmwareengine_subnet.ds"),
				),
			},
		},
	})
}

func testAccDataSourceVmwareEngineSubnetConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_vmwareengine_subnet" "ds" {
	name = "vsan"
	parent =  "%{private_cloud_id}"
}
`, context)
}

func testAccCheckGoogleVmwareengineSubnetMeta(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find subnet data source: %s", n)
		}

		_, ok = rs.Primary.Attributes["name"]
		if !ok {
			return errors.New("can't find 'name' attribute")
		}

		_, ok = rs.Primary.Attributes["state"]
		if !ok {
			return errors.New("can't find 'state' attribute")
		}

		_, ok = rs.Primary.Attributes["type"]
		if !ok {
			return errors.New("can't find 'type' attribute")
		}

		return nil
	}
}
