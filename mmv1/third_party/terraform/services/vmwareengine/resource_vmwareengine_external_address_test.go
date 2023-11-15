package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccVmwareengineExternalAddress_vmwareEngineExternalAddressUpdate(t *testing.T) {
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
				Config: testVmwareEngineExternalAddressConfig(context, "description1", "192.168.0.66"),
			},
			{
				ResourceName:            "google_vmwareengine_external_address.vmw-engine-external-address",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
			{
				Config: testVmwareEngineExternalAddressConfig(context, "description2", "192.168.0.67"),
			},
			{
				ResourceName:            "google_vmwareengine_external_address.vmw-engine-external-address",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "name"},
			},
		},
	})
}

func testVmwareEngineExternalAddressConfig(context map[string]interface{}, description string, internalIp string) string {
	context["internal_ip"] = internalIp
	context["description"] = description
	return acctest.Nprintf(`
resource "google_vmwareengine_external_address" "vmw-engine-external-address" {
	name = "tf-test-sample-external-address%{random_suffix}"
	parent =  "%{private_cloud_id}"
	internal_ip = "%{internal_ip}"
	description = "%{description}"
}
`, context)
}
