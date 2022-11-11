---
page_title: "ipam_ip_range Resource - terraform-provider-gcp-ipam-autopilot"
subcategory: ""
description: |-
  The resource ipam_ip_range allows you to dynamically allocate subnet range within routing domain.
---

# Resource `ipam_ip_range`

The resource ipam_ip_range allows you to dynamically allocate subnet range within routing domain ensuring there is no overlap with subnets in existing VPCs and previosly allocated IP ranges.

## Example Usage

```terraform
resource "ipam_ip_range" "frontend" {
  range_size = 8
  name = "frontend range"
  domain = ipam_routing_domain.prod.id
  cidr = "10.0.0.0/8"
}

resource "ipam_ip_range" "pod_range" {
  range_size = 22
  name = "gke pod range"
  domain = ipam_routing_domain.prod.id
  parent = ipam_ip_range.frontend.cidr
}
```

## Argument Reference

- `name` - (Required) Name of the IP range.
- `domain` - (Optional) Associate IP range with specific routing domain.
- `parent` - (Optional) Allocate subnet from within parent IP range.
- `cidr` - (Optional) Specify requested CIDR explicitly.
- `range_size` - (Optional) Specify the size of requested CIDR.

## Attributes Reference

No additional attributes exposed.

