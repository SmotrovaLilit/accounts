terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.33.0"
    }
  }
}

locals {
  database_dsn = "postgres://${var.database_user_name}:${var.database_user_pass}@${var.database_address}:5432/${var.database_name}?sslmode=disable"
}

resource "kubernetes_config_map" "config" {
  metadata {
    name      = "config"
    namespace = "default"
  }

  data = {
    "config.yml" = <<EOT
service:
  host: "0.0.0.0"
  port: 8080
database_dsn: "${local.database_dsn}"
EOT
  }
}

resource "kubernetes_config_map" "migrations" {
    metadata {
        name      = "migrations"
        namespace = "default"
    }

    data = var.migrations
}

# Deployment
resource "kubernetes_deployment" "backend" {
  metadata {
    name = var.name
    labels = {
      app = var.name
    }
  }

  wait_for_rollout = false

  spec {
    replicas = 1
    selector {
      match_labels = {
        app = var.name
      }
    }

    template {
      metadata {
        labels = {
          app = var.name
        }
      }

      spec {
        container {
          image             = var.image
          name              = var.name
          image_pull_policy = "Always"

          port {
            container_port = 8080
          }

          volume_mount {
            name       = "config-volume"
            mount_path = "/root/config.yml"
            sub_path = "config.yml"
          }
        }

        volume {
          name = "config-volume"

          config_map {
            name = kubernetes_config_map.config.metadata[0].name
          }
        }
      }
    }
  }
}

# Service
resource "kubernetes_service" "backend_service" {
  metadata {
    name = "${var.name}-service"
  }

  spec {
    selector = {
      app = kubernetes_deployment.backend.metadata[0].labels.app
    }

    port {
      port        = 80
      target_port = 8080
    }

    type = "ClusterIP"
  }
}

#Ingress
resource "kubernetes_ingress_v1" "backend_ingress" {
  depends_on = [kubernetes_service.backend_service]

  metadata {
    name = "${var.name}-ingress"
    annotations = {
      "nginx.ingress.kubernetes.io/rewrite-target" = "/"
    }
  }

  spec {
    ingress_class_name = "nginx"

    rule {
      host = var.domain
      http {
        path {
          path = "/"

          backend {
            service {
              name = kubernetes_service.backend_service.metadata[0].name
              port {
                number = 80
              }
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_job" "migration" {
  metadata {
    name = "${var.name}-migrate-job"
  }

  spec {
    template {
      metadata {
        labels = {
          app = "migrate"
        }
      }

      spec {
        container {
          name  = "migrate"
          image = "migrate/migrate:v4.14.1"
          command = [
            "-path=/migrations",
            "-database", local.database_dsn,
            "up"
          ]
          volume_mount {
            mount_path = "/migrations"
            name       = "migrations-volume"
          }
        }

        restart_policy = "Never"

        volume {
          name = "migrations-volume"
          config_map {
            name = kubernetes_config_map.migrations.metadata[0].name
          }
        }
      }
    }

    backoff_limit = 4
  }
}