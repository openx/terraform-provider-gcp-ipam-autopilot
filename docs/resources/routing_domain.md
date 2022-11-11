---
page_title: "ipam_routing_domain Resource - terraform-provider-gcp-ipam-autopilot"
subcategory: ""
description: |-
  The resource ipam_routing_domain allows you to create routing domain and optionally linking it to VPCs in your infrastructure.
---

# Resource `ipam_routing_domain`

The ipam_routing_domain resource allows you to create routing domain. Routing domain can be identified as collection of interconnected networks that exchange routes between each other..

## Example Usage

```terraform
resource "ipam_routing_domain" "prod" {
  name = "Production Network Domain"
  vpcs = [ google_compute_network.prod.self_link ]
}
```

## Argument Reference

- `name` - (Required) Unique name of network domain.
- `vpcs` - (Optional) List of VPC network references.

## Attributes Reference

No additional attributes exposed.

