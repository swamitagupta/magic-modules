package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceVmwareengineExternalAddress_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"private_cloud_id":  acctest.BootstrapVmwareenginePrivateCloud(t, "test-tl-pc", acctest.BootstrapVmwareengineNetwork(t, "test-nw"), true),
		"network_policy_id": acctest.BootstrapVmwareengineNetworkPolicy(t, "test-np", acctest.BootstrapVmwareengineNetwork(t, "test-nw")),
		"random_suffix":     acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineExternalAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareengineExternalAddress_ds(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_external_address.ds", "google_vmwareengine_external_address.vmw-engine-external-address", map[string]struct{}{}),
				),
			},
		},
	})
}

func testAccVmwareengineExternalAddress_ds(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_external_address" "vmw-engine-external-address" {
    name = "tf-test-sample-external-address%{random_suffix}"
    parent =  "%{private_cloud_id}"
    internal_ip = "192.168.0.66"
    description = "Sample description."
}

data "google_vmwareengine_external_address" "ds" {
	name = google_vmwareengine_external_address.vmw-engine-external-address.name
	parent =  "%{private_cloud_id}"
	depends_on = [
		google_vmwareengine_external_address.vmw-engine-external-address,
	]
}
`, context)
}
