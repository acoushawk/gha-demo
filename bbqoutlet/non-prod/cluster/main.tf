locals {
	config_path = "${path.module}/config.yaml"
	config = merge(
		yamldecode(fileexists(local.config_path) ? file(local.config_path) : yamlencode({}))
	)

	cluster_labels = zipmap(
      [for label_key, _ in local.config.cluster_labels : label_key],
      [for _, label_value in local.config.cluster_labels : label_value]
    )
	
	node_pools_labels = merge(
		{all = {customer = element(split("lw-", local.config.name), length(split("lw-", local.config.name))-1)}},
		zipmap(
			[for node_pool in local.config.node_pools : node_pool["name"]],
			[for node_pool in local.config.node_pools : { fusion_node_type = lookup(node_pool,"fusion_node_type",node_pool["name"]) }]
		)
	)

	node_pools_tags = merge(
		{all = [ element(split("lw-", local.config.name), length(split("lw-", local.config.name))-1) ]},
		zipmap(
			[for node_pool in local.config.node_pools : node_pool["name"]],
			[for node_pool in local.config.node_pools : [ node_pool["name"] ]]
		)
	)
}

module "k8s" {
	source = "git@github.com:lucidworks/terraform-google-gke.git?ref=0.3.7"
	
	name                = local.config.name
	region              = local.config.region
	project_id          = local.config.project_id

	release_channel     = "STABLE"
	
	cluster_labels = local.cluster_labels
	
	node_pool_list = local.config.node_pools
	
	node_pools_labels = local.node_pools_labels
	
	node_pools_tags = local.node_pools_tags
	
	resource_usage_export_config = {
		dataset_id : "gke_usage"
	}
	
	node_pools_taints = {}
}

###############################################################################
# output
###############################################################################
output "endpoint" {
	value = module.k8s.endpoint
}

output "get_cluster_config" {
	value     = "gcloud container clusters get-credentials ${local.config.name}-${local.config.region} --region ${local.config.region} --project ${local.config.project_id}"
}
