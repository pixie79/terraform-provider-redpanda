resource "redpanda_schema" "test2" {
  subject     = "test2"
  schema      = file("${path.root}/schemas/onboarding/test2.avsc")
  schema_type = "AVRO"
}
