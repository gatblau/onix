# Provider <img src="../../../docs/pics/ox.png" width="125" height="125" align="right">

The Terraform provider for Onix allows Terraform to manage configuration information in the Onix CMDB as Terraform [resources](https://www.terraform.io/docs/configuration/resources.html) and [data sources](https://www.terraform.io/docs/configuration/data-sources.html).

<a name="toc"></a>
### Table of Contents [(index)](./../readme.md)

- [Terraform Provider for Onix](#terraform-provider-for-onix)
    - [Table of Contents (index)](#table-of-contents-index)
  - [Connection information (up)](#connection-information-up)
  - [Installation (up)](#installation-up)
  - [Resources (up)](#resources-up)
  - [Data Sources (up)](#data-sources-up)

<a name="connection-information"></a>
## Connection information ([up](#toc))

Connection information can be provided by adding the service URI, user and password in-line in the Onix provider as shown below.

__Example Usage:__

```hcl-terraform
provider "ox" {
  uri  = "http://localhost:8080"
  user = "admin"
  pwd  = "0n1x"
}
```
<a name="installation"></a>
## Installation ([up](#toc))

In order to use this provider, it must be manually installed, since terraform init cannot automatically download it.

Install the provider by placing its plugin executable in the user plugins directory. 
The user plugins directory is in one of the following locations, depending on the host operating system:

|Operating system|	User plugins directory|
|---|---|
|Windows	| %APPDATA%\terraform.d\plugins|
|All other systems|	~/.terraform.d/plugins|

Alternatevely, the provider can also be placed under the local folder where Terraform is run from. 

For example: **working_directory/.terraform.d/plugins**

Once the plugin is placed in the correct location, the [terraform init](https://www.terraform.io/docs/commands/init.html) command has to be run to initialise the working directory ready for use.

<a name="resources"></a>
## Resources ([up](#toc))

Resources are the most important element in the Terraform language. Each resource block describes one or more infrastructure objects, such as compute instances, higher-level components such as DNS records or in this case, end points in the [Onix Web API](../../../docs/wapi.md).

The Onix Provider for Terraform, acts as a Restful client allowing to perform requests to the Onix Web API.

The list of resources offered in this provider is as follows:

| Resource | Description |
|---|---|
| [ox_model](./docs/rs_ox_model.md) | Creates, updates or deletes a [modelKey](../../../models/readme.md). |
| [ox_item](./docs/rs_ox_item.md) | Creates, updates or deletes a configuration item. |
| [ox_item_type](./docs/rs_ox_item_type.md) | Creates, updates or deletes a configuration item type. |
| [ox_link](./docs/rs_ox_link.md) | Creates, updates or deletes a configuration item link. |
| [ox_link_type](./docs/rs_ox_link_type.md) | Creates, updates or deletes a configuration item link type. |
| [ox_link_rule](./docs/rs_ox_link_rule.md) | Creates, updates or deletes a configuration item link rule. |

<a name="data-sources"></a>
## Data Sources ([up](#toc))

Data sources allow data to be fetched or computed for use elsewhere in Terraform configuration. Use of Onix data sources allows a Terraform configuration to make use of information defined in the CMDB.

The list of the data sources offered in this provider is as follows:

| Source | Description |
|---|---|
| [ox_model_data](./docs/rs_ox_model_data.md) | Queries a [modelKey](../../../models/readme.md). |
| [ox_item_data](./docs/rs_ox_item_data.md) | Queries a configuration item. |
| [ox_item_type_data](./docs/rs_ox_item_type_data.md) | Queries a configuration item type. |
| [ox_link_data](./docs/rs_ox_link_data.md) | Queries a configuration item link. |
| [ox_link_type_data](./docs/rs_ox_link_type_data.md) | Queries a configuration item link type. |
| [ox_link_rule_data](./docs/rs_ox_link_rule_data.md) | Queries a configuration item link rule. |

