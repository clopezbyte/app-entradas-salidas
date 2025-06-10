variable "project_id" {
    type = string
    description = "GCP Project ID"
    default = "b-materials"
}

variable "region" {
    type = string
    description = "Project Region"
    default = "us-central1"
}

variable "keys_gcs_bucket" {
    type = string
    description = "Google Cloud Storage Bucket for keys"
    default = "b-materials-keys"
}

variable "composer_service_account" {
  type        = string
  description = "Name of the Cloud Composer service account"
  default     = "cloud-run-invoker@b-materials.iam.gserviceaccount.com"
}

variable "composer_dag_bucket" {
  type = string
  description = "Cloud Composer DAG bucket storage"
  default = "us-central1-buho-claw-envir-ca3ff9c9-bucket"
}