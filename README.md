# Terraform Red Panda Provider

This is a custom Terraform provider for managing topics and schemas in Red Panda.

## Prerequisites

- [Go](https://golang.org/dl/) >= 1.16
- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Red Panda](https://vectorized.io/redpanda/) >= 21.11.1

## Building the provider

1. Clone the repository:

```bash
git clone https://github.com/pixie79/terraform-provider-redpanda.git
cd terraform-provider-internal
```


2. Build the provider:

```bash
go build -o terraform-provider-internal
```


3. Install the provider:

```bash
mkdir -p ~/.terraform.d/plugins/example.com/local/internal/0.1.0/linux_amd64
mv terraform-provider-internal ~/.terraform.d/plugins/example.com/local/internal/0.1.0/linux_amd64/
```

Replace `linux_amd64` with the appropriate directory for your platform if needed.

## Using the provider

1. Create a `main.tf` file:

```hcl
terraform {
  required_providers {
    redpanda = {
      source  = "example.com/local/internal"
      version = "0.1.0"
    }
  }
}

provider "redpanda" {
  api_url = "http://localhost:8081"
}

resource "redpanda_topic" "example_topic" {
  name               = "example-topic"
  partitions         = 3
  replication_factor = 1
}

resource "schema" "example_schema" {
  subject    = "example-topic-value"
  schema     = "{\"type\": \"record\", \"name\": \"example\", \"fields\": [{\"name\": \"field1\", \"type\": \"string\"}]}"
  schema_type = "AVRO"
}
```

Replace the api_url with the URL to your Red Panda schema registry.

2. Initialize Terraform:

```bash
terraform init
```

3. Apply the configuration:

```bash
terraform apply
```

4. To destroy the created resources:

```bash
terraform destroy
```


# License

This project is licensed under the MIT License.


This `README.md` file provides an overview of the Terraform Red Panda provider, including prerequisites, build instructions, installation, and usage.
