provider "kubernetes" {
	host			= length(data.terraform_remote_state.gke_cluster.outputs) > 0 ? "https://${data.terraform_remote_state.gke_cluster.outputs.endpoint}" : "localhost"
  	config_path		= "~/.kube/config"
	config_context 	= "gke_${local.config.project_id}_${local.config.region}_${local.config.cluster_name}-${local.config.region}"
}

provider "helm" {
	kubernetes {
		host			= length(data.terraform_remote_state.gke_cluster.outputs) > 0 ? "https://${data.terraform_remote_state.gke_cluster.outputs.endpoint}" : "localhost"
  		config_path		= "~/.kube/config"
		config_context 	= "gke_${local.config.project_id}_${local.config.region}_${local.config.cluster_name}-${local.config.region}"
 	}
}

provider "kubectl" {
	host			= length(data.terraform_remote_state.gke_cluster.outputs) > 0 ? "https://${data.terraform_remote_state.gke_cluster.outputs.endpoint}" : "localhost"
  	config_path    	= "~/.kube/config"
	config_context 	= "gke_${local.config.project_id}_${local.config.region}_${local.config.cluster_name}-${local.config.region}"
}