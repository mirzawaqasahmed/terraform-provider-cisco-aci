package aci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	cage "github.com/ignw/cisco-aci-go-sdk/src/service"
)

//TODO: Add advanced test cases to cover setting properties and EPGs

func TestAccAciBridgeDomain_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBridgeDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBridgeDomainConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBridgeDomainExists("aci_bridge_domain.basic"),
					resource.TestCheckResourceAttr(
						"aci_bridge_domain.basic", "name", "bd1"),
					resource.TestCheckResourceAttr(
						"aci_bridge_domain.basic", "description", "terraform test bridge domain"),
					resource.TestCheckResourceAttr(
						"aci_bridge_domain.basic", "domain_name", "uni/tn-IGNW-tenant1/BD-bd1"),
				),
			},
		},
	})
}

func testAccCheckAciBridgeDomainExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*cage.Client)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Error -> Resource Not found: %s", n)
		}

		id := rs.Primary.Attributes["id"]

		bd, err := client.BridgeDomains.Get(id)

		if err != nil || bd == nil {
			return fmt.Errorf("Error retreiving bridge domain id: %s", id)
		}

		return nil
	}
}

func testAccCheckAciBridgeDomainDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*cage.Client)

	err := checkDestroy("aci_bridge_domain.basic", state, func(s string) (interface{}, error) {
		return client.BridgeDomains.Get(s)
	})

	if err != nil {
		return err
	}

	err = checkDestroy("aci_tenant.basic", state, func(s string) (interface{}, error) {
		return client.Tenants.Get(s)
	})

	if err != nil {
		return err
	}

	return nil
}

const testAccCheckAciBridgeDomainConfigBasic = `
resource "aci_tenant" "basic" {
    name = "IGNW-tenant1"
    description = "my first tenant"
}

resource "aci_bridge_domain" "basic" {
	name = "bd1"
	description = "terraform test bridge domain"
	tenant_id = "${aci_tenant.basic.id}"
	type = "regular"
}
`
