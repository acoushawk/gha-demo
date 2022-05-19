locals {
	config_path = "${path.module}/config.yaml"
	config = merge(
		yamldecode(fileexists(local.config_path) ? file(local.config_path) : yamlencode({}))
	)
}

data "terraform_remote_state" "gke_cluster" {
	backend = "gcs"
	config = {
		bucket  = "lw-${local.config.project_id}"
		prefix  = "clusters-tfstate/${local.config.cluster_name}"
	}
}

resource "tls_private_key" "identity" {
  algorithm = "ECDSA"
  ecdsa_curve  = "P521"
}

resource "random_password" "github_webhook_secret" {
  length           = 16
  special          = true
  override_special = "_%@"
}

resource "github_repository_deploy_key" "main" {
  title      = "argocd-${local.config.repo_name}-main"
  repository = "${local.config.repo_name}"
  key        = tls_private_key.identity.public_key_openssh
  read_only  = true
}

resource "github_repository_webhook" "argocd" {
  repository = "${local.config.repo_name}"
  events = ["push"]
  configuration {
    url          = "https://${local.config.argocd_host}/api/webhook"
    content_type = "json"
    secret       = "${random_password.github_webhook_secret.result}"
  }
}

resource "helm_release" "argocd" {
  name  = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  version    = "4.6.0"
  namespace  = "argocd"
  create_namespace = true

  values = [
    templatefile(
      "${path.module}/templates/argocd-values.yaml.tpl",
      {	
        "repoName"                   = "${local.config.repo_name}"
        "repoURL"                    = "git@github.com:lucidworks-managed-fusion/${local.config.repo_name}.git"
        "repoSshPrivateKey"          = "${jsonencode(tls_private_key.identity.private_key_pem)}"
        "repoSshPublicKey"           = "${tls_private_key.identity.public_key_openssh}"
        "clusterType"                = "${local.config.cluster_type}"
        "fusionDevHelmRepoUsername"  = "${var.fusion_dev_helm_repo_username}"
        "fusionDevHelmRepoPassword"  = "${var.fusion_dev_helm_repo_password}"
        "argocdHost"                 = "${local.config.argocd_host}"
        "argocdHostTlsName"          = "${local.config.argocd_host_tls_name}"
        "githubWebhookSecret"        = "${random_password.github_webhook_secret.result}"
        "oktaClientID"               = "${var.okta_client_id}"
        "oktaClientSecret"           = "${var.okta_client_secret}"
        "configSyncRepoURL"          = "git@github.com:lucidworks-managed-fusion/config-sync-config.git"
        "configSyncSshPrivateKey"    = "${jsonencode(var.config_sync_private_key)}"
        "configSyncSshPublicKey"     = "${var.config_sync_public_key}" 
      }
    ),
    templatefile(
      "${path.module}/templates/argocd-applicationset-values.yaml.tpl",
      {	
      }
    ),
    templatefile(
      "${path.module}/templates/argocd-notifications-values.yaml.tpl",
      {	
        "slackNotificationBotToken" = "${var.slack_notification_bot_token}"
      }
    )
  ]
}
