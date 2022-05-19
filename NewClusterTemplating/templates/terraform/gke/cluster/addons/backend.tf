terraform {
	backend "gcs" {
		bucket  = "lw-[[.ProjectID]]"
		prefix  = "clusters-tfstate/[[.ClusterName]]-addons"
	}
	required_providers {
		kubectl = {
			source  = "gavinbunney/kubectl"
		}
	}
}