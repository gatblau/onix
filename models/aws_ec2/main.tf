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

resource "ox_link_rule" "aws_transit_gateway_to_aws_vpc" {
  key                 = "AWS_TRANSIT_GATEWAY->AWS_VPC"
  name                = "AWS Transit Gateway to VPC rule"
  description         = "Allow linking AWS Transit Gateways with VPCs"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_transit_gateway.key}"
  end_item_type_key   = "${ox_item_type.aws_vpc.key}"
}

resource "ox_link_rule" "aws_vpc_to_aws_route53_zone" {
  key                 = "AWS_VPC->AWS_ROUTE53_ZONE"
  name                = "AWS VPC to Route 53 Zone rule"
  description         = "Allow AWS VPCs with Route 53 Zones"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_vpc.key}"
  end_item_type_key   = "${ox_item_type.aws_route53_zone.key}"
}

resource "ox_link_rule" "aws_route53_zone_to_aws_route53_record" {
  key                 = "AWS_ROUTE53_ZONE->AWS_ROUTE53_RECORD"
  name                = "AWS Route 53 Zone to Route 53 Record rule"
  description         = "Allow linking Route 53 Zones with Route 53 Records"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_route53_zone.key}"
  end_item_type_key   = "${ox_item_type.aws_route53_record.key}"
}

resource "ox_link_rule" "aws_vpc_to_aws_vpc_peering_connection" {
  key                 = "AWS_VPC->AWS_VPC_PEERING_CONNECTION"
  name                = "AWS VPC to VPC Peering Connection rule"
  description         = "Allow linking VPCs with VPC Peering Connections"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_vpc.key}"
  end_item_type_key   = "${ox_item_type.aws_vpc_peering_connection.key}"
}

resource "ox_link_rule" "aws_vpc_to_aws_internet_gateway" {
  key                 = "AWS_VPC->AWS_VPC_TO_AWS_INTERNET_GATEWAY"
  name                = "AWS VPC to Internet Gateway rule"
  description         = "Allow linking AWS VPCs with Internet Gateways"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_vpc.key}"
  end_item_type_key   = "${ox_item_type.aws_internet_gateway.key}"
}

resource "ox_link_rule" "aws_vpc_to_aws_security_group" {
  key                 = "AWS_VPC->AWS_SECURITY_GROUP"
  name                = "AWS VPC to Security Group rule"
  description         = "Allow linking AWS VPCs with Security Groups"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_vpc.key}"
  end_item_type_key   = "${ox_item_type.aws_security_group.key}"
}

resource "ox_link_rule" "aws_security_group_to_aws_security_rule" {
  key                 = "AWS_SECURITY_GROUP->AWS_SECURITY_RULE"
  name                = "AWS Security Group to Security Rule rule"
  description         = "Allow linking Security Groups with Security Rules"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_security_group.key}"
  end_item_type_key   = "${ox_item_type.aws_security_rule.key}"
}

resource "ox_link_rule" "aws_security_group_to_aws_network_interface" {
  key                 = "AWS_SECURITY_GROUP->AWS_NETWORK_INTERFACE"
  name                = "AWS Security Group to Network Interface rule"
  description         = "Allow linking AWS Security Groups with Network Interfaces"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_security_group.key}"
  end_item_type_key   = "${ox_item_type.aws_network_interface.key}"
}

resource "ox_link_rule" "aws_security_group_to_aws_instance" {
  key                 = "AWS_SECURITY_GROUP->AWS_INSTANCE"
  name                = "AWS Security Group to AWS Instance rule"
  description         = "Allow linking AWS Security Groups with Instances"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_security_group.key}"
  end_item_type_key   = "${ox_item_type.aws_instance.key}"
}

resource "ox_link_rule" "aws_vpc_to_aws_instance" {
  key                 = "AWS_VPC->AWS_INSTANCE"
  name                = "AWS VPC to Instance rule"
  description         = "Allow linking AWS VPCs with Instances"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_vpc.key}"
  end_item_type_key   = "${ox_item_type.aws_instance.key}"
}

resource "ox_link_rule" "aws_vpc_to_aws_vpc_ipv4_cidr_block_association" {
  key                 = "AWS_VPC->AWS_VPC_IPV4_CIDR_BLOCK_ASSOCIATION"
  name                = "AWS VPC to VPC IPV4 CIDR Block Association rule"
  description         = "Allow linking AWS VPC with IPV4 CIDR Block Associations"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_vpc.key}"
  end_item_type_key   = "${ox_item_type.aws_vpc_ipv4_cidr_block_association.key}"
}

resource "ox_link_rule" "aws_vpc_ipv4_cidr_block_association_to_aws_subnet" {
  key                 = "AWS_VPC_IPV4_CIDR_BLOCK_ASSOCIATION->AWS_SUBNET"
  name                = "AWS VPC IPV4 CIDR Block Association to Subnet rule"
  description         = "Allow linking ... with ..."
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_vpc_ipv4_cidr_block_association.key}"
  end_item_type_key   = "${ox_item_type.aws_subnet.key}"
}

resource "ox_link_rule" "aws_subnet_to_aws_instance" {
  key                 = "AWS_SUBNET->AWS_INSTANCE"
  name                = "AWS Subnet to Instance rule"
  description         = "Allow linking AWS Subnets with Instances"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_subnet.key}"
  end_item_type_key   = "${ox_item_type.aws_instance.key}"
}

resource "ox_link_rule" "aws_route_table_to_aws_instance" {
  key                 = "AWS_ROUTE_TABLE->AWS_INSTANCE"
  name                = "AWS Route Table to Instance rule"
  description         = "Allow linking AWS Route Tables with Instances"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_route_table.key}"
  end_item_type_key   = "${ox_item_type.aws_instance.key}"
}

resource "ox_link_rule" "aws_route_table_to_aws_route" {
  key                 = "AWS_ROUTE_TABLE->AWS_ROUTE"
  name                = "AWS Route Table to Route rule"
  description         = "Allow linking AWS Route Tables with Routes"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_route_table.key}"
  end_item_type_key   = "${ox_item_type.aws_route.key}"
}

resource "ox_link_rule" "aws_instance_to_aws_network_interface" {
  key                 = "AWS_INSTANCE->AWS_NETWORK_INTERFACE"
  name                = "AWS Instance to Network Interface rule"
  description         = "Allow linking AWS Instances with Network Interfaces"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_instance.key}"
  end_item_type_key   = "${ox_item_type.aws_network_interface.key}"
}

resource "ox_link_rule" "aws_instance_to_aws_ebs_volume" {
  key                 = "AWS_INSTANCE->AWS_EBS_VOLUME"
  name                = "AWS Instance to EBS Volume rule"
  description         = "Allow linking AWS Instances with EBS Volumes"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_instance.key}"
  end_item_type_key   = "${ox_item_type.aws_ebs_volume.key}"
}

resource "ox_link_rule" "aws_instance_to_aws_key_pair" {
  key                 = "AWS_INSTANCE->AWS_KEY_PAIR"
  name                = "AWS Instance to Key Pair rule"
  description         = "Allow linking AWS Instances with Key Pairs"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_instance.key}"
  end_item_type_key   = "${ox_item_type.aws_key_pair.key}"
}

resource "ox_link_rule" "aws_instance_to_aws_eip" {
  key                 = "AWS_INSTANCE->AWS_EIP"
  name                = "AWS Instance to Elastic IP rule"
  description         = "Allow linking AWS Instances with Elastic IPs"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_instance.key}"
  end_item_type_key   = "${ox_item_type.aws_eip.key}"
}

resource "ox_link_rule" "aws_security_group_to_aws_lb" {
  key                 = "AWS_SECURITY_GROUP->AWS_LB"
  name                = "AWS Security Group to AWS Load Balancer rule"
  description         = "Allow linking AWS Security Groups with Load Balancers"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_security_group.key}"
  end_item_type_key   = "${ox_item_type.aws_lb.key}"
}

resource "ox_link_rule" "aws_lb_to_aws_lb_listener" {
  key                 = "AWS_LB->AWS_LB_LISTENER"
  name                = "AWS Load Balancer to Load Balancer Listener rule"
  description         = "Allow linking Load Balancers with Load Balancer Listeners"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_lb.key}"
  end_item_type_key   = "${ox_item_type.aws_lb_listener.key}"
}

resource "ox_link_rule" "aws_lb_to_aws_lb_target_group" {
  key                 = "AWS_LB->AWS_LB_TARGET_GROUP"
  name                = "AWS Load Balancer to Load Balancer Target Group rule"
  description         = "Allow linking AWS Load Balancers with Load Balancer Target Groups"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_lb.key}"
  end_item_type_key   = "${ox_item_type.aws_lb_target_group.key}"
}

resource "ox_link_rule" "aws_lb_target_group_to_aws_lb_target_group_attachment" {
  key                 = "AWS_LB_TARGET_GROUP->AWS_LB_TARGET_GROUP_ATTACHMENT"
  name                = "AWS Load Balancer Target Group to Load Balancer Target Group attachment rule"
  description         = "Allow linking Load Balancer Target Groups with Group Attachments"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_lb_target_group.key}"
  end_item_type_key   = "${ox_item_type.aws_lb_target_group_attachment.key}"
}

resource "ox_link_rule" "aws_instance_to_aws_instance_profile" {
  key                 = "AWS_INSTANCE->AWS_INSTANCE_PROFILE"
  name                = "AWS Instance to Instance Profile rule"
  description         = "Allow linking AWS Instances with Instance Profiles"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_instance.key}"
  end_item_type_key   = "${ox_item_type.aws_iam_instance_profile.key}"
}

resource "ox_link_rule" "aws_iam_instance_profile_to_aws_iam_role" {
  key                 = "AWS_IAM_INSTANCE_PROFILE->AWS_IAM_ROLE"
  name                = "AWS Id & Access Mgmt Instance Profile to Role rule"
  description         = "Allow linking AWS Id & Access Mgmt Instance Profile with Roles"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_iam_instance_profile.key}"
  end_item_type_key   = "${ox_item_type.aws_iam_role.key}"
}

resource "ox_link_rule" "aws_iam_role_to_aws_iam_role_policy_attachment" {
  key                 = "AWS_IAM_ROLE->AWS_IAM_ROLE_POLICY_ATTACHMENT"
  name                = "AWS IAM Role to Role Policy Attachment rule"
  description         = "Allow linking AWS IAM Roles with Role Policy Attachments"
  link_type_key       = "${ox_link_type.aws_ec2_link.key}"
  start_item_type_key = "${ox_item_type.aws_iam_role.key}"
  end_item_type_key   = "${ox_item_type.aws_iam_role_policy_attachment.key}"
}



