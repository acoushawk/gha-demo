resource "github_repository" "customer_repo" {
  name        = "testcustomer-config"
  description = "Deployment Manifests For testcustomer"
  visibility  = "internal"
  gitignore_template = "Terraform"
  auto_init = true
  allow_merge_commit = false
  delete_branch_on_merge = true

  lifecycle {
    prevent_destroy = true
  }
}

resource "github_branch_protection" "customer_repo_branch_protection" {
  repository_id = github_repository.customer_repo.name
  pattern          = "main"
  enforce_admins   = false
  allows_deletions = false

  required_status_checks {
    strict  = true
  }

  required_pull_request_reviews {
    dismiss_stale_reviews  = true
    require_code_owner_reviews = true
    required_approving_review_count = 1
  }
}

resource "github_team_repository" "repo_access" {
  team_id    = "cloud-infrastructure"
  repository = github_repository.customer_repo.name
  permission = "push"
}