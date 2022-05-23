terraform {
	backend "gcs" {
		bucket  = "lw-managed-fusion"
		prefix  = "clusters-tfstate/lw-testcustomer-non-prod"
	}
}