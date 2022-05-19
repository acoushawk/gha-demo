terraform {
	backend "gcs" {
		bucket  = "lw-managed-fusion"
		prefix  = "repo-tfstate/[[.CustomerName]]-config"
	}
	required_providers {
		github = {
			source  = "integrations/github"
		}
	}
}