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
  name        = "default"
  location_id = var.region
  project     = var.project_id
  type        = "FIRESTORE_NATIVE"
}

# Retrieve the self-link of the default network
data "google_compute_network" "default" {
  name    = "default"
  project = var.project_id
}

# VPC connector for API
resource "google_vpc_access_connector" "in-out-goods-app-vpc-connector" {
  name          = "in-out-vpc-conn"
  region        = var.region
  network       = data.google_compute_network.default.self_link
  ip_cidr_range = "10.8.0.0/28" 
  min_instances = 2
  max_instances = 10
  machine_type = "f1-micro"
}

# Cloud Run Service (API)
resource "google_cloud_run_service" "in-out-goods-app-api" {
  name     = "in-out-goods-app-api"
  location = var.region
  project  = var.project_id

  metadata {
    annotations = {
        "run.googleapis.com/ingress"                     = "all"
    }
  }

  template {
    metadata {
      annotations = {
        "run.googleapis.com/vpc-access-connector"        = google_vpc_access_connector.in-out-goods-app-vpc-connector.name
        "run.googleapis.com/vpc-egress"                  = "private-ranges-only"
      }
    }

    spec {
      containers {
        image = "gcr.io/${var.project_id}/in-out-goods-app-api:prodv7"

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
  
  depends_on = [google_vpc_access_connector.in-out-goods-app-vpc-connector]
}

# Allow outside invocations to API
resource "google_cloud_run_service_iam_member" "unauthenticated_invoker" {
  service    = google_cloud_run_service.in-out-goods-app-api.name
  location   = var.region
  project    = var.project_id
  role       = "roles/run.invoker"
  member     = "allUsers"
}

# Firestore to BQ pipeline Cloud Run Job
resource "google_cloud_run_v2_job" "in-out-analytics-pipeline" {
  name = "in-out-analytics-pipeline"
  location = var.region
  project = var.project_id
  deletion_protection = "false"
  template {
    template {
      containers {
        image = "gcr.io/b-materials/in-out-analytics-pipeline:dev"  #Docker image

        env {
          name = "PROJECT_ID"
          value = "b-materials"
        }

      }
    }
  }
}