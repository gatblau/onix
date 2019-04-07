provider "ox" {
  uri  = "http://localhost:8080"
  user = "admin"
  pwd  = "0n1x"
}

resource "ox_model" "aws_ec2" {
  key         = "AWS_EC2"
  name        = "Amazon Elastic Compute Cloud"
  description = "Virtual Machine infrastructure provided by Amazon Web Services."
}

resource "ox_item_type" "aws_customer_gateway" {
  key         = "AWS_CUSTOMER_GATEWAY"
  name        = "AWS Customer Gateway"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_vpn_connection" {
  key         = "AWS_VPN_CONNECTION"
  name        = "AWS VPN Connection"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_transit_gateway" {
  key         = "AWS_TRANSIT_GATEWAY"
  name        = "AWS Transit Gateway"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_transit_gateway_attachment" {
  key         = "AWS_TRANSIT_GATEWAY_ATTACHMENT"
  name        = "AWS Transit Gateway Attachment"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_vpc" {
  key         = "AWS_VPC"
  name        = "AWS VPC"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_route53_record" {
  key         = "AWS_ROUTE53_RECORD"
  name        = "AWS Route 53 Record"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_route53_zone" {
  key         = "AWS_ROUTE53_ZONE"
  name        = "AWS Route 53 Zone"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_route53_zone_association" {
  key         = "AWS_ROUTE53_ZONE_ASSOCIATION"
  name        = "AWS Route 53 Zone Association"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_vpc_ipv4_cidr_block_association" {
  key         = "AWS_VPC_IPV4_CIDR_BLOCK_ASSOCIATION"
  name        = "AWS VPC IPV4 CIDR Block Association"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_subnet" {
  key         = "AWS_SUBNET"
  name        = "AWS Subnet"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_network_interface" {
  key         = "AWS_NETWORK_INTERFACE"
  name        = "AWS Network Interface"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_route_table_association" {
  key         = "AWS_ROUTE_TABLE_ASSOCIATION"
  name        = "AWS Route Table Association"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_route_table" {
  key         = "AWS_ROUTE_TABLE"
  name        = "AWS Route Table"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_route" {
  key         = "AWS_ROUTE"
  name        = "AWS Route"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_instance" {
  key         = "AWS_INSTANCE"
  name        = "AWS Instance"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_eip_association" {
  key         = "AWS_EIP_ASSOCIATION"
  name        = "AWS Elastic IP Association"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_eip" {
  key         = "AWS_EIP"
  name        = "AWS Elastic IP"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_iam_instance_profile" {
  key         = "AWS_IAM_INSTANCE_PROFILE"
  name        = "AWS Id & Access Mgtmt Instance Profile"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_iam_role" {
  key         = "AWS_IAM_ROLE"
  name        = "AWS Id & Access Mgtmt Role"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_iam_role_policy_attachment" {
  key         = "AWS_IAM_ROLE_POLICY_ATTACHMENT"
  name        = "AWS Id & Access Mgtmt Role Policy Attachment"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_volume_attachment" {
  key         = "AWS_VOLUME_ATTACHMENT"
  name        = "AWS Volume Attachment"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_ebs_volume" {
  key         = "AWS_EBS_VOLUME"
  name        = "AWS EBS Volume"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_internet_gateway" {
  key         = "AWS_INTERNET_GATEWAY"
  name        = "AWS Internet Gateway"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_vpc_peering_connection" {
  key         = "AWS_VPC_PEERING_CONNECTION"
  name        = "AWS VPC Peering Connection"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_security_group" {
  key         = "AWS_SECURITY_GROUP"
  name        = "AWS Security Group"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_security_rule" {
  key         = "AWS_SECURITY_RULE"
  name        = "AWS Security Rule"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_lb" {
  key         = "AWS_LB"
  name        = "AWS Load Balancer"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_lb_listener" {
  key         = "AWS_LB_LISTENER"
  name        = "AWS Load Balancer Listener"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_lb_target_group" {
  key         = "AWS_LB_TARGET_GROUP"
  name        = "AWS Load Balancer Target Group"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_lb_target_group_attachment" {
  key         = "AWS_LB_TARGET_GROUP_ATTACHMENT"
  name        = "AWS Load Balancer Target Group Attachment"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_item_type" "aws_key_pair" {
  key         = "AWS_KEY_PAIR"
  name        = "AWS Key Pair"
  description = ""
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_link_type" "aws_ec2_link" {
  key         = "AWS_EC2_LINK"
  name        = "AWS EC2 Link"
  description = "Links AWS EC2 model items."
  model_key   = "${ox_model.aws_ec2.key}"
}

resource "ox_link_rule" "aws_customer_gateway_to_aws_vpn_connection" {
  key                 = "AWS_CUSTOMER_GATEWAY->AWS_VPN_CONNECTION"
  name                = "AWS Customer Gateway to VPN Connection rule"
  description         = "Allow linking Customer Gateways with VPN Connections"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_customer_gateway.key}"
  end_item_type_key   = "${ox_item_type.aws_vpn_connection.key}"
}

resource "ox_link_rule" "aws_vpn_connection_to_aws_transit_gateway" {
  key                 = "AWS_VPN_CONNECTION->AWS_TRANSIT_GATEWAY"
  name                = "AWS VPN Connection to AWS Transit Gateway rule"
  description         = "Allow linking AWS VPN Connections with Transit Gateways"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_vpn_connection.key}"
  end_item_type_key   = "${ox_item_type.aws_transit_gateway.key}"
}

//resource "ox_link_rule" "_To_" {
//  key                 = "->"
//  name                = "... to ... rule"
//  description         = "Allow linking ... with ..."
//  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
//  start_item_type_key = "${ox_item_type..key}"
//  end_item_type_key   = "${ox_item_type..key}"
//}
