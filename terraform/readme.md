# Onix Terraform Integration

Onix provides two ways of integrating with Terraform as discussed below.

## Terraform provider for Onix

[Terraform providers](https://www.terraform.io/docs/providers/index.html) are used to create, manage, and update resources.

Onix has a provider that acts as a RESTful client for the Onix Web API.

It can be used to create, update, query and delete configuration data in the Onix configuration database.

It is provided to satisfy a wide range of scenarios where a simple to use CLI client is required.

The documentation for the provider can be found [here](provider/readme.md).

## Terraform backend for Onix

A [Terraform Backend](https://www.terraform.io/docs/backends/index.html) keeps track of Terraform state.

Onix has a backend that acts as a facade for the Onix configuration database.

When you use this backend, all configuration changes are automatically recorded in the Onix database.

If you use Terraform to provision / manage infrastructure resources, then this backend will ensure everything is
securely recorded in Onix and can be audited or used for further automation or management.

The documentation of the backend can be found [here](backend/readme.md).

