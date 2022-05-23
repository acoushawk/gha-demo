terraform {
	backend "gcs" {
		bucket  = "lw-managed-fusion"
		prefix  = "clusters-tfstate/lw-bbqoutlet-non-prod"
	}
}