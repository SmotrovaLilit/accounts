variable "location" {
  type    = string
  default = "europe-west1-b"
}

variable "name" {
  type    = string
  default = "postgres"
}

variable "database_version" {
  type    = string
  default = "POSTGRES_14"
}

variable "user_name" {
  type    = string
  default = "root"
}

variable "user_pass" {
  type = string
}
