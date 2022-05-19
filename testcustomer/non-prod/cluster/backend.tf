terraform {
	backend "gcs" {
		bucket  = "lw-managed-fusion-dev"
		prefix  = "clusters-tfstate/lw-testcustomer-non-prod"
	}
}