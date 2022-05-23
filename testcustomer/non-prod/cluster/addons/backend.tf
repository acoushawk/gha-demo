terraform {
	backend "gcs" {
		bucket  = "lw-managed-fusion"
		prefix  = "clusters-tfstate/lw-testcustomer-non-prod-addons"
	}
	required_providers {
		kubectl = {
			source  = "gavinbunney/kubectl"
		}
	}
}