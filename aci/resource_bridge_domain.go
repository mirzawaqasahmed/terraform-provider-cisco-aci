package aci

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	cage "github.com/ignw/cisco-aci-go-sdk/src/service"
)

func resourceAciBridgeDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciBridgeDomainCreate,
		Update: resourceAciBridgeDomainUpdate,
		Read:   resourceAciBridgeDomainRead,
		Delete: resourceAciBridgeDomainDelete,

		SchemaVersion: 1,

		Schema: MergeSchemaMaps(
			GetBaseSchema(),
			map[string]*schema.Schema{
				"tenant_id": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"type": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"arp_flood": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"optimize_wan": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"mode_detect_mode": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"allow_intersite_bum_traffic": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"intersite_l2_stretch": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"ip_learning": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"limit_ip_to_subnets": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"ll_ip_address": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Default:  "::",
				},
				"mac": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Default:  "00:22:BD:F8:19:FF",
				},
				"multi_dest_forwarding": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Default:  "bd-flood",
				},
				"multicast": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"unicast_route": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"unknown_unicast_mac": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Default:  "proxy",
				},
				"unknown_multicast_mac": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Default:  "flood",
				},
				"virtual_mac": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Default:  "not-applicable",
				},
				"subnets": &schema.Schema{
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: GetBaseSchema(),
					},
				},
				"endpoint_groups": &schema.Schema{
					Type:     schema.TypeList,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		),
	}
}

func resourceAciBridgeDomainFieldMap() map[string]string {

	return MergeStringMaps(GetBaseFieldMap(),
		map[string]string{
			"Type":                     "type",
			"ArpFlood":                 "arp_flood",
			"OptimizeWan":              "optimize_wan",
			"MoveDetectMode":           "mode_detect_mode",
			"AllowIntersiteBumTraffic": "allow_intersite_bum_traffic",
			"IntersiteL2Stretch":       "intersite_l2_stretch",
			"IpLearning":               "ip_learning",
			"LimitIpToSubnets":         "limit_ip_to_subnets",
			"LLIpAddress":              "ll_ip_address",
			"MAC":                      "mac",
			"MultiDestForwarding":      "multi_dest_forwarding",
			"Multicast":                "multicast",
			"UnicastRoute":             "unicast_route",
			"UnknownUnicastMAC":        "unknown_unicast_mac",
			"UnknownMulticastMAC":      "unknown_multicast_mac",
			"VirtualMAC":               "virtual_mac",
		})
}

func resourceAciBridgeDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cage.Client)
	resource := &AciResource{d}

	tenant, err := ValidateAndFetchTenant(d, meta)

	if err != nil {
		return fmt.Errorf("Error creating bridge domain: %s", err.Error())
	}

	bridgeDomain := client.BridgeDomains.New(resource.Get("name").(string), resource.Get("description").(string))

	resource.MapFieldsToAci(resourceAciBridgeDomainFieldMap(), bridgeDomain)

	tenant.AddBridgeDomain(bridgeDomain)

	dn, err := client.BridgeDomains.Save(bridgeDomain)

	if err != nil {
		return fmt.Errorf("Error saving bridge domain: %s", err.Error())
	}

	d.Set("domain_name", dn)
	d.SetId(dn)

	return nil
}

func resourceAciBridgeDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cage.Client)
	resource := &AciResource{d}

	if resource.Id() == "" {
		return fmt.Errorf("Error missing resource identifier")
	}

	bridgeDomain, err := client.BridgeDomains.Get(resource.Id())

	if err != nil {
		return fmt.Errorf("Error updating application profile id: %s", resource.Id())
	}

	resource.MapFields(resourceAciBridgeDomainFieldMap(), bridgeDomain)

	return nil
}

func resourceAciBridgeDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	// HACK: currently same implementation as create
	return resourceAciBridgeDomainCreate(d, meta)
}

func resourceAciBridgeDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cage.Client)
	return DeleteAciResource(d, client.BridgeDomains.Delete)
}

/*
func (d *AciResource) SetSubnets(items []*cage.Subnet) {
	subnets := make([]map[string]interface{}, len(items))

	for i, item := range items {
		subnets[i] = d.ConvertToBaseMap(&item.ResourceAttributes)
	}
	d.Set("subnets", subnets)
}
*/
