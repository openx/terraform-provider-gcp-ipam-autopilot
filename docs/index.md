---
page_title: "Provider: GCP IPAM Autopilot"
subcategory: "Networking"
description: |-
  Terraform provider for interacting with IPAM Autopilot API developed by Google Cloud PSO.
---

# GCP IPAM Autopilot Provider

-> Visit the [original IPAM Autopilot page](https://github.com/GoogleCloudPlatform/professional-services/tree/main/tools/ipam-autopilot) for detailed information on API and how to host service used by provider on Google Cloud Run.

The IPAM Autopilot provider is used to interact with a IPAM Autopilot application hosted on Cloud Run. This version of provider is clone of original provider supplied and hosted by with IPAM Autopilot and addressing the cumbersome nature of compilation and distribution of provider binaries by Terraform and IPAM Autopilot service itself.
By publishing this provider using standard HashiCorp process, provider distribution functionality in original version of IPAM Autopilot can be disabled thus simplifying significantly the amount of configuration needed to start using the service.

Use the navigation to the left to read about the available resources supported by provider.

## Example Usage

Do not keep your authentication password in HCL for production environments, use Terraform environment variables.

```terraform
terraform {
  required_providers {
    ipam = {
      version = "~> 0.4"
      source = "openx/gcp-ipam-autopilot"
    }
  }
}

provider "ipam" {
  url = "https://<cloud run hostname>"
}

resource "ipam_routing_domain" "test" {
  name = "Test Domain"
}

resource "ipam_ip_range" "main" {
  range_size = 8
  name = "main range"
  domain = ipam_routing_domain.test.id
  cidr = "10.0.0.0/8"
}

resource "ipam_ip_range" "sub1" {
  range_size = 24
  name = "sub range 1"
  domain = ipam_routing_domain.test.id
  parent = ipam_ip_range.main.cidr
  cidr = "10.0.1.0/24"
}

```

### Authentication

The provider uses application default credentials to authenticate against the backend. Alternatively you can set GOOGLE_CREDENTIALS variable to point or contain JSON account key or directly provide an identity token via the GCP_IDENTITY_TOKEN environment variable to access the Cloud Run instance, the audience for the identity token should be the domain of the Cloud Run service.

