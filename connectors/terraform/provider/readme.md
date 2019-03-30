# Terraform Provider for Onix

The Terraform provider for Onix allows Terraform to manage configuration information in the Onix CMDB as Terraform [resources](https://www.terraform.io/docs/configuration/resources.html) and [data sources](https://www.terraform.io/docs/configuration/data-sources.html).

## Connection information

Connection information can be provided by adding the service URI, user and password in-line in the Onix provider as shown below.

Usage:

```hcl-terraform
provider "onix" {
  uri  = "onix-uri"
  user = "my-user"
  pwd  = "my-pwd"
}
```
## Installing the provider

In order to use this provider, it must be manually installed, since terraform init cannot automatically download it.

Install the provider by placing its plugin executable in the user plugins directory. 
The user plugins directory is in one of the following locations, depending on the host operating system:

|Operating system|	User plugins directory|
|---|---|
|Windows	| %APPDATA%\terraform.d\plugins|
All other systems|	~/.terraform.d/plugins|

Once the plugin is installed, terraform init can initialise it normally.


## Resources

### ox_item

Provides a resource to create/update or delete configuration item data in the CMDB.

#### Example Usage
```hcl-terraform
resource "ox_item" "ITEM_01" {
  "key": "ITEM_01"
  "name": "Item 01"
  "description": "Item 01 description."
  "type": "NODE"
  "modelKey": "MODEL_01"
  "tag": [
    "tag1", "tag2", "tag3"
  ]
  "meta": {
    "key1" : {
      "key2": "value2",
      "key3": "value3"
    }
  }
  "attribute": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

## Data Sources

