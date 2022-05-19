terraform {
	backend "gcs" {
		bucket  = "lw-[[.ProjectID]]"
		prefix  = "clusters-tfstate/[[.ClusterName]]"
	}
}