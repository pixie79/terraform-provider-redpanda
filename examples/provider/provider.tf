provider "redpanda" {
  api_url = "http://localhost:18081"
}


terraform {
  required_version = ">= 1.3.6"
  required_providers {
    redpanda = {
      source  = "local.com/pixie79/redpanda"
      version = "0.0.1"
    }
  }
  #   backend "s3" {
  #     encrypt = true
  #   }
}
