variable "cluster_name" {
  type    = string
}

variable "workers_count" {
  type    = number
  default = 3
}

variable "kubernetes_version" {
  type    = string
  default = "1.31"
}

variable "cluster_location" {
  type    = string
  default = "europe-west1-b"
}
