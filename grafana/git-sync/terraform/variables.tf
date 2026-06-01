variable "grafana_url" {
  type    = string
  default = "https://sso-grafana-sandbox.apps.gold.devops.gov.bc.ca"
}

variable "grafana_service_account_token" {
  type    = string
  default = "<replace me>"
}

variable "repo_url" {
  type    = string
  default = ""
}

variable "repo_branch" {
  type    = string
  default = "main"
}

variable "repo_path" {
  type    = string
  default = "/"
}

variable "repo_sync_interval_seconds" {
  type    = number
  default = 3600
}

variable "github_app_id" {
  type    = string
  default = "<replace me>"
}

variable "github_installation_id" {
  type    = string
  default = "<replace me>"
}

variable "github_private_key_b64_encoded" {
  type    = string
  default = ""
}
