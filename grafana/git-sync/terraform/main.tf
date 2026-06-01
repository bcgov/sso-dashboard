resource "grafana_apps_provisioning_connection_v0alpha1" "github_app" {
  metadata {
    uid = "sso-grafana-dashboards-github-app-connection"
  }

  spec {
    title       = "GitHub App Connection"
    description = "GitHub connection authenticated with a GitHub App"
    type        = "github"
    url         = "https://github.com"

    github {
      app_id          = var.github_app_id
      installation_id = var.github_installation_id
    }
  }

  secure {
    private_key = {
      create = var.github_private_key_b64_encoded
    }
  }
  secure_version = 1
}


resource "grafana_apps_provisioning_repository_v0alpha1" "github_repo" {
  depends_on = [grafana_apps_provisioning_connection_v0alpha1.github_app]

  metadata {
    uid = "sso-grafana-dashboards-github-app-repo"
  }

  spec {
    title       = "GitHub App Repository"
    description = "Folder-scoped GitHub repository authenticated via a referenced GitHub App connection"
    type        = "github"

    workflows = ["branch"]

    sync {
      enabled          = true
      target           = "folder"
      interval_seconds = var.repo_sync_interval_seconds
    }

    github {
      url    = var.repo_url
      branch = var.repo_branch
      path   = var.repo_path
    }

    connection {
      name = grafana_apps_provisioning_connection_v0alpha1.github_app.metadata.uid
    }
  }
}
