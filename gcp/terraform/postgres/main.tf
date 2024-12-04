resource "google_sql_database_instance" "default" {
  name             = var.name
  database_version = var.database_version
  region           = var.location

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_database" "default" {
  name     = var.name
  instance = google_sql_database_instance.default.name
}

resource "google_sql_user" "default" {
  name     = var.user_name
  password = var.user_pass
  instance = google_sql_database_instance.default.name
}