variable "name" {
  type = string
}

variable "image" {
  type = string
}

variable "domain" {
  type = string
}

variable "database_address" {
  type = string
}

variable "database_user_pass" {
  type = string
}

variable "database_user_name" {
  type = string
}

variable "database_name" {
  type = string
}

variable "migrations" {
  type = map(string)
  default = {}
}