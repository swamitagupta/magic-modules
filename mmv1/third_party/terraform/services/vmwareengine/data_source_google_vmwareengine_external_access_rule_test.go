package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceVmwareengineExternalAccessRule_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"private_cloud_id":  acctest.BootstrapVmwareenginePrivateCloud(t, "test-tl-pc", acctest.BootstrapVmwareengineNetwork(t, "test-nw"), true),
		"network_policy_id": acctest.BootstrapVmwareengineNetworkPolicy(t, "test-np", acctest.BootstrapVmwareengineNetwork(t, "test-nw")),
		"random_suffix":     acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineExternalAccessRuleDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareengineExternalAccessRule_ds(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores("data.google_vmwareengine_external_access_rule.ds", "google_vmwareengine_external_access_rule.vmw-engine-external-access-rule", map[string]struct{}{}),
				),
			},
		},
	})
}

func testAccVmwareengineExternalAccessRule_ds(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vmwareengine_external_address" "external-access-rule-ea" {
    name = "tf-test-sample-external-address%{random_suffix}"
    parent =  "%{private_cloud_id}"
    internal_ip = "192.168.0.72"
    description = "Sample description."
}

resource "google_vmwareengine_external_access_rule" "vmw-engine-external-access-rule" {
    name = "tf-test-sample-external-access-rule%{random_suffix}"
    parent =  "%{network_policy_id}"
    description = "Sample Description"
    priority = 101
    action = "ALLOW"
    ip_protocol = "tcp"
    source_ip_ranges {
        ip_address = "192.168.50.0"
    }
    source_ports = ["80"]
    destination_ip_ranges {
        external_address = google_vmwareengine_external_address.external-access-rule-ea.id
    }
    destination_ports = ["433"]
}

data "google_vmwareengine_external_access_rule" "ds" {
	name = google_vmwareengine_external_access_rule.vmw-engine-external-access-rule.name
	parent =  "%{network_policy_id}"
	depends_on = [
		google_vmwareengine_external_access_rule.vmw-engine-external-access-rule,
	]
}
`, context)
}
