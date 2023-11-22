package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVmwareengineSubnet_vmwareEngineUserDefinedSubnetUpdate(t *testing.T) {
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
				Config: testVmwareEngineSubnetConfig(context, "192.168.80.0/26"),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_subnet.ds", "google_vmwareengine_subnet.vmw-engine-subnet", map[string]struct{}{}),
				),
			},
			{
				ResourceName:            "google_vmwareengine_subnet.vmw-engine-subnet",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
			{
				Config: testVmwareEngineSubnetConfig(context, "192.168.90.0/26"),
			},
			{
				ResourceName:            "google_vmwareengine_subnet.vmw-engine-subnet",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
		},
	})
}

func testVmwareEngineSubnetConfig(context map[string]interface{}, ipCidrRange string) string {
	context["ip_cidr_range"] = ipCidrRange
	return acctest.Nprintf(`

resource "google_vmwareengine_subnet" "vmw-engine-subnet" {
	name = "service-2"
	parent =  "%{private_cloud_id}"
	ip_cidr_range = "%{ip_cidr_range}"
}

data "google_vmwareengine_subnet" ds {
	name = "service-2"
	parent = "%{private_cloud_id}"
	depends_on = [
    google_vmwareengine_subnet.vmw-engine-subnet,
  ]
}

`, context)
}
