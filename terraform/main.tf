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

# Cloud Run Job for dbt transformation from Bronze to Silver
resource "google_cloud_run_v2_job" "in-out-analytics-dbt-job" {
  name = "in-out-analytics-dbt-job"
  location = var.region
  project = var.project_id
  deletion_protection = "false"

  template {
    template {
      containers {
        image = "gcr.io/b-materials/in-out-analytics-dbt-job:prodv2"  #Docker image

        env {
          name  = "DBT_KEY_URL"
          value = "gs://adhoc_data_buho_picks/b-materials-dc73796cdc19.json"
        }
      }
    }
  }
}

#Cloud Composer
provider "google-beta" {
  alias   = "composer"
  project = var.project_id
  region  = var.region
}

# Cloud Composer Environment
resource "google_composer_environment" "buho_claw_environment" {
  provider = google-beta.composer
  name     = "buho-claw-environment"
  region   = var.region

  labels = {
    environment = "production"
    team        = "dmca"
    project     = var.project_id
  }

  config {
    software_config {
      image_version = "composer-3-airflow-2.10.5-build.3"

      airflow_config_overrides = {
        core-dag_serialization = "True"
        core-default_timezone  = "America/Mexico_City"
      }
    }

    node_config {
      service_account = var.composer_service_account
    }

    resilience_mode = "STANDARD_RESILIENCE"
    environment_size = "ENVIRONMENT_SIZE_SMALL"

    maintenance_window {
      start_time = "2023-01-01T00:00:00Z"
      end_time   = "2023-01-01T07:00:00Z"
      recurrence = "FREQ=WEEKLY;BYDAY=SA,SU"
    }

    database_config {
      zone = "us-central1-a"  
    }

    recovery_config {
      scheduled_snapshots_config {
        enabled = false
      }
    }

    #Network config
    # network_config {
    #   network                    = var.network_id
    #   subnetwork                 = var.subnet_id
    #   enable_private_environment = true
    #   master_ipv4_cidr_block     = "10.0.0.0/28"
    #   worker_ipv4_cidr_block     = "10.0.1.0/24"
    # }

    #Dataplex Lineage config

    web_server_network_access_control {
      allowed_ip_range {
        value       = "0.0.0.0/0"
        description = "Allow access from all IP addresses"
      }
    }
  }
}

# In Out Analytics Pipeline DAG for Cloud Composer
resource "google_storage_bucket_object" "in_out_analytics_pipeline_dag" {
  name   = "dags/in_out_analytics_pipeline_dag.py"  
  bucket = var.composer_dag_bucket  
  source = "C:/Users/Carlos Lopez/app-entradas-salidas/analytics/dags/in_out_analytics_pipeline_dag.py"
}