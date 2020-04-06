# Data Format for Terraform State

When the Terraform state is imported in Onix, Terraform resources are separated into their own configuration items.
These items are linked using Terraform State Links as it is shown in the meta model below:

![Terraform Meta Model](./tf_model.png)

The Terraform resource "instances" information is stored in the "meta" field. Other data items suchs as mode, type, name and provider are stored in the configuration item "attributes".

For an example of the Terraform state see [here](terraform.tfstate.json). The equivalent state in Onix format is [here](state_items.json).