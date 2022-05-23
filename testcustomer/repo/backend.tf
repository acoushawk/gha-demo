terraform {
	backend "gcs" {
		bucket  = "lw-managed-fusion"
		prefix  = "repo-tfstate/testcustomer-config"
	}
	required_providers {
		github = {
			source  = "integrations/github"
		}
	}
}