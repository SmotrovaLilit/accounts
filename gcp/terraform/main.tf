terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.12.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.33.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.16.1"
    }
  }
}

locals {
  region              = "europe-west1"
  zone                = "europe-west1-b"
  network_name        = "default-network"
  accounts_migrations = {
    for file in fileset("../../migrations", "*.sql") :
    basename(file) => file("../../migrations/${file}")
  }
}

# Provider is configured using environment variables: GOOGLE_REGION, GOOGLE_PROJECT, GOOGLE_CREDENTIALS.
# This can be set statically, if preferred. See docs for details.
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#full-reference
provider "google" {
  project = "dsxack"
  region  = local.zone
  zone    = local.zone
}

# Configure kubernetes provider with Oauth2 access token.
# https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/client_config
# This fetches a new token, which will expire in 1 hour.
data "google_client_config" "default" {
  depends_on = [module.gke-cluster]
}

# Defer reading the cluster data until the GKE cluster exists.
data "google_container_cluster" "default" {
  name = local.cluster_name
  depends_on = [module.gke-cluster]
}

# Defer reading the database instance data until the database instance exists.
data "google_sql_database_instance" "default" {
  name = local.database_name
  depends_on = [module.postgres]
}

# Kubernetes config for GKE.
provider "kubernetes" {
  host  = "https://${data.google_container_cluster.default.endpoint}"
  token = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(
    data.google_container_cluster.default.master_auth[0].cluster_ca_certificate,
  )
}

# Helm config for GKE.
provider "helm" {
  kubernetes {
    host  = "https://${data.google_container_cluster.default.endpoint}"
    token = data.google_client_config.default.access_token
    cluster_ca_certificate = base64decode(
      data.google_container_cluster.default.master_auth[0].cluster_ca_certificate,
    )
  }
}

module "gke-cluster" {
  source           = "./gke-cluster"
  cluster_name     = local.cluster_name
  cluster_location = local.zone
}

module "postgres" {
  source       = "./postgres"
  user_pass    = local.database_user_pass
  user_name    = local.database_user_name
  name         = local.database_name
  location     = local.region
}

module "kubernetes-config" {
  depends_on = [module.gke-cluster]
  source = "./kubernetes-config"
}

module "accounts" {
  depends_on = [module.kubernetes-config, module.postgres]
  source             = "./backend"
  domain             = "accounts.test"
  image              = "gcr.io/dsxack/accounts:v2"
  name               = "accounts"
  migrations         = local.accounts_migrations
  database_address   = data.google_sql_database_instance.default.public_ip_address
  database_user_pass = local.database_user_pass
  database_user_name = local.database_user_name
  database_name      = local.database_name
}

output "google_sql_database_ip_address" {
  value = data.google_sql_database_instance.default.public_ip_address
}

output "accounts_migrations" {
  value = local.accounts_migrations
}