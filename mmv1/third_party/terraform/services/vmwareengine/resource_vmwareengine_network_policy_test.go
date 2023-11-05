package vmwareengine_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccVmwareengineNetworkPolicy_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVmwareengineNetworkPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVmwareengineNetworkPolicy_config(context, false, false),
			},
			{
				ResourceName:            "google_vmwareengine_network_policy.vmw-engine-network-policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
			{
				Config: testAccVmwareengineNetworkPolicy_config(context, true, true),
			},
			{
				ResourceName:            "google_vmwareengine_network_policy.vmw-engine-network-policy",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name"},
			},
		},
	})
}

func testAccVmwareengineNetworkPolicy_config(context map[string]interface{}, internetAccess bool, externalIp bool) string {
	context["internet_access"] = internetAccess
	context["external_ip"] = externalIp	
	
	return acctest.Nprintf(`
resource "google_vmwareengine_network" "network-policy-nw" {
    name              = "tf-test-sample-nw%{random_suffix}"
    location          = "global" 
    type              = "STANDARD"
    description       = "VMwareEngine standard network sample"
}

resource "google_vmwareengine_network_policy" "vmw-engine-network-policy" {
    location = "%{region}"
    name = "tf-test-sample-network-policy%{random_suffix}"

    internet_access {
        enabled = "%{internet_access}"
    }

    external_ip {
        enabled = "%{external_ip}"
    }

    edge_services_cidr = "192.168.30.0/26"
    vmware_engine_network = google_vmwareengine_network.network-policy-nw.id
}
`, context)
}
