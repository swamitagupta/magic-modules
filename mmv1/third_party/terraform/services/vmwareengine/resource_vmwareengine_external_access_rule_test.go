package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVmwareengineExternalAccessRule_vmwareEngineExternalAccessRuleUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"private_cloud_id":  acctest.BootstrapVmwareenginePrivateCloud(t, "test-tl-pc", acctest.BootstrapVmwareengineNetwork(t, "test-nw"), true),
		"network_policy_id": acctest.BootstrapVmwareengineNetworkPolicy(t, "test-np", acctest.BootstrapVmwareengineNetwork(t, "test-nw")),
		"random_suffix":     acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testVmwareEngineExternalAccessRuleConfig(context, "description1", "ALLOW", 101, "192.168.40.0"),
			},
			{
				ResourceName:            "google_vmwareengine_external_access_rule.vmw-engine-external-access-rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
			{
				Config: testVmwareEngineExternalAccessRuleConfig(context, "description2", "DENY", 105, "192.168.50.0"),
			},
			{
				ResourceName:            "google_vmwareengine_external_access_rule.vmw-engine-external-access-rule",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
		},
	})
}

func testVmwareEngineExternalAccessRuleConfig(context map[string]interface{}, description string, action string, priority int, sourceIpAddress string) string {
	context["description"] = description
	context["action"] = action
	context["priority"] = priority
	context["source_ip_address"] = sourceIpAddress
	return acctest.Nprintf(`
resource "google_vmwareengine_external_address" "external-access-rule-ea" {
	name = "tf-test-sample-external-address%{random_suffix}"
	parent =  "%{private_cloud_id}"
	internal_ip = "192.168.0.73"
	description = "Sample description."
}

resource "google_vmwareengine_external_access_rule" "vmw-engine-external-access-rule" {
	name = "tf-test-sample-external-access-rule%{random_suffix}"
	parent =  "%{network_policy_id}"
	description = "%{description}"
	priority = "%{priority}"
	action = "%{action}"
	ip_protocol = "tcp"
	source_ip_ranges {
		ip_address = "%{source_ip_address}"
	}
	source_ports = ["80"]
	destination_ip_ranges {
		external_address = google_vmwareengine_external_address.external-access-rule-ea.id
	}
	destination_ports = ["433"]
}
`, context)
}
