provider "google" {
  project = var.project_id
  region  = var.region
}

# Google Cloud Storage Bucket for keys
data "google_storage_bucket" "keys_gcs_bucket" {
  name = var.keys_gcs_bucket
}

# Reference the existing service account key object in the bucket
data "google_storage_bucket_object" "service_account_key" {
  name   = "b-materials-dc73796cdc19.json" #File in bucket
  bucket = data.google_storage_bucket.keys_gcs_bucket.name
}

# Firestore Database
resource "google_firestore_database" "default" {
  name        = "(default)"
  location_id = var.region
  project     = var.project_id
  type        = "FIRESTORE_NATIVE"
}

# VPC connector for API
resource "google_vpc_access_connector" "in-out-goods-app-vpc-connector" {
  name          = "in-out-goods-app-vpc-connector"
  region        = var.region
  network       = "(default)" 
  ip_cidr_range = "10.8.0.0/28" 
}

# Cloud Run Service (API)
resource "google_cloud_run_service" "in-out-goods-app-api" {
  name     = "in-out-goods-app-api"
  location = var.region
  project  = var.project_id

  template {
    metadata {
      annotations = {
        "run.googleapis.com/vpc-access-connector"        = google_vpc_access_connector.in-out-goods-app-vpc-connector.name
        "run.googleapis.com/vpc-egress"                  = "private-ranges-only"
        "run.googleapis.com/ingress"                     = "all"
      }
    }

    spec {
      containers {
        image = "gcr.io/${var.project_id}/in-out-goods-app-api:dev"

        ports {
                container_port = 8080
            }

        env {
          name  = "FIRESTORE_PROJECT_ID"
          value = var.project_id
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

