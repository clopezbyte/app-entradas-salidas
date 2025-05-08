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