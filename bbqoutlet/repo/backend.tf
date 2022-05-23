terraform {
	backend "gcs" {
		bucket  = "lw-managed-fusion"
		prefix  = "repo-tfstate/bbqoutlet-config"
	}
	required_providers {
		github = {
			source  = "integrations/github"
		}
	}
}