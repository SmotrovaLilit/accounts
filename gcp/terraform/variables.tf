resource "random_password" "database_user_pass" {
  length           = 16
  special          = true
  override_special = "!"
}

resource "random_id" "cluster_name" {
  byte_length = 5
}

locals {
  cluster_name       = "accounts-k8s-${random_id.cluster_name.hex}"
  database_user_pass = random_password.database_user_pass.result
  database_user_name      = "root"
  database_name      = "accounts"
}
