variable "slack_notification_bot_token" {
  type        = string
  description = "The oauth token for the slack notification bot"
}

variable "fusion_dev_helm_repo_username" {
  type        = string
  description = "The username for the fusion-dev-helm artifactory repo"
}

variable "fusion_dev_helm_repo_password" {
  type        = string
  description = "The password for the fusion-dev-helm artifactory repo"
}

variable "okta_client_id" {
  type        = string
  description = "The okta app client id to use for oauth"
}

variable "okta_client_secret" {
  type        = string
  description = "The okta app client secret to use for oauth"
}

variable "config_sync_private_key" {
  type        = string
  description = "The confyg-sync repo private key"
}

variable "config_sync_public_key" {
  type        = string
  description = "The confyg-sync repo public key"
}