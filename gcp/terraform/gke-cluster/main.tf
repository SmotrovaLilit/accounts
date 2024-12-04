terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.12.0"
    }
  }
}

locals {
  location_name   = var.cluster_location
}

data "google_container_engine_versions" "supported" {
  location       = local.location_name
  version_prefix = var.kubernetes_version
}

resource "google_container_cluster" "default" {
  name               = var.cluster_name
  location           = local.location_name
  initial_node_count = var.workers_count
  min_master_version = data.google_container_engine_versions.supported.latest_master_version
  node_version       = data.google_container_engine_versions.supported.latest_master_version

  node_config {
    machine_type = "n1-standard-4"

    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }

  deletion_protection = false
}
