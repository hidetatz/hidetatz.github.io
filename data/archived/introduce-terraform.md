Introduce Terraform with simple template---2018-04-30 06:18:53

I have created simple terraform template for AWS.

### At first: code is

[here](https://github.com/ygnmhdtt/minimal-terraform) .

### Terraform

I love AWS, and use it every day at my work, or as my hobby.
But, I always create aws resources with aws console or aws-cli or aws-shell.
This is my first experience of terraform and I want to log it.

### Install terraform

As Linux, very simple.

```
$ wget https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_linux_amd64.zip
$ unzip unzip terraform_0.11.7_linux_amd64.zip
$ mv terraform /usr/local/bin # or to any your $PATH
$ which terraform
```

### Makefile

I wrote Makefile like this. (I alse love Makefile)

```
TERRAFORM       = terraform

init:
        $(TERRAFORM) init

keygen:
        ssh-keygen -t rsa -f $(F)
        chmod 400 $F*
        mv $(F)* ~/.ssh/

plan:
        $(TERRAFORM) plan

show:
        $(TERRAFORM) show

apply:
        $(TERRAFORM) apply

destroy:
        $(TERRAFORM) destroy
```

### write .tf files

If you want to try terraform with my sample,
only you have to do is `clone` .

```
$ git clone git@github.com:ygnmhdtt/minimal-terraform.git
```

and follow README.

### init

```
$ make init  
terraform init
Initializing modules...  
- module.vpc
  Getting source "./modules/vpc"  
- module.subnet
  Getting source "./modules/subnet"  
- module.route_table  
  Getting source "./modules/route_table"
- module.security_group
  Getting source "./modules/security_group"  
- module.ec2  
  Getting source "./modules/ec2"  

Initializing provider plugins...
- Checking for available provider plugins on https://releases.hashicorp.com...
- Downloading plugin for provider "aws" (1.16.0)...
  
The following providers do not have any version constraints in configuration,
so the latest version was installed.  
  
To prevent automatic upgrades to new major versions that may contain breaking
changes, it is recommended to add version = "..." constraints to the
corresponding provider blocks in configuration, with the constraint strings
suggested below.
  
* provider.aws: version = "~> 1.16"

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

### plan

```
$ make plan  
terraform plan  
Refreshing Terraform state in-memory prior to plan...  
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.  
  
  
------------------------------------------------------------------------
  
An execution plan has been generated and is shown below.  
Resource actions are indicated with the following symbols:  
  + create  
  
Terraform will perform the following actions:  

  + module.ec2.aws_instance.bastion  
      id:                                          <computed>
      ami:                                         "ami-0c11b26d"
      associate_public_ip_address:                 "true"  
      availability_zone:                           <computed>  
      ebs_block_device.#:                          <computed>
      ebs_optimized:                               "false"  
      ephemeral_block_device.#:                    <computed>  
      get_password_data:                           "false"  
      instance_state:                              <computed>
      instance_type:                               "t2.micro"  
      ipv6_address_count:                          <computed>
      ipv6_addresses.#:                            <computed>  
      key_name:                                    "tf-key"  
      monitoring:                                  "false"  
      network_interface.#:                         <computed>
      network_interface_id:                        <computed>
      password_data:                               <computed>
      placement_group:                             <computed>  
      primary_network_interface_id:                <computed>  
      private_dns:                                 <computed>  
      private_ip:                                  <computed>  
      public_dns:                                  <computed>  
      public_ip:                                   <computed>  
      root_block_device.#:                         "1"  
      root_block_device.0.delete_on_termination:   "false"  
      root_block_device.0.volume_id:               <computed>  
      root_block_device.0.volume_size:             "50"  
      root_block_device.0.volume_type:             "standard"
      security_groups.#:                           <computed>
      source_dest_check:                           "true"
      subnet_id:                                   "${var.subnet_id}"
      tags.%:                                      "1"
      tags.Name:                                   "main"
      tenancy:                                     <computed>
      volume_tags.%:                               <computed>
      vpc_security_group_ids.#:                    <computed>

  + module.ec2.aws_key_pair.tf-key  
      id:                                          <computed>
      fingerprint:                                 <computed>
      key_name:                                    "tf-key"  
      public_key:                                  "my ssh key"
  
  + module.route_table.aws_internet_gateway.public  
      id:                                          <computed>  
      tags.%:                                      "1"  
      tags.Name:                                   "main"  
      vpc_id:                                      "${var.vpc_id}"
  
  + module.route_table.aws_route_table.public  
      id:                                          <computed>
      propagating_vgws.#:                          <computed>
      route.#:                                     "1"  
      route.~2593658367.cidr_block:                "0.0.0.0/0"
      route.~2593658367.egress_only_gateway_id:    ""  
      route.~2593658367.gateway_id:                "${aws_internet_gateway.public.id}"
      route.~2593658367.instance_id:               ""  
      route.~2593658367.ipv6_cidr_block:           ""  
      route.~2593658367.nat_gateway_id:            ""  
      route.~2593658367.network_interface_id:      ""  
      route.~2593658367.vpc_peering_connection_id: ""  
      tags.%:                                      "1"  
      tags.Name:                                   "main"  
      vpc_id:                                      "${var.vpc_id}"  
  
  + module.route_table.aws_route_table_association.public  
      id:                                          <computed>  
      route_table_id:                              "${aws_route_table.public.id}"
      subnet_id:                                   "${var.subnet_id}"

  + module.security_group.aws_security_group.22_pxy  
      id:                                          <computed>
      arn:                                         <computed>
      description:                                 "Permit ssh from proxy server."
      egress.#:                                    "1"
      egress.482069346.cidr_blocks.#:              "1"
      egress.482069346.cidr_blocks.0:              "0.0.0.0/0"
      egress.482069346.description:                ""
      egress.482069346.from_port:                  "0"
      egress.482069346.ipv6_cidr_blocks.#:         "0"
      egress.482069346.prefix_list_ids.#:          "0"
      egress.482069346.protocol:                   "-1"
      egress.482069346.security_groups.#:          "0"
      egress.482069346.self:                       "false"
      egress.482069346.to_port:                    "0"
      ingress.#:                                   "1"
      ingress.2541437006.cidr_blocks.#:            "1"
      ingress.2541437006.cidr_blocks.0:            "0.0.0.0/0"
      ingress.2541437006.description:              ""
      ingress.2541437006.from_port:                "22"
      ingress.2541437006.ipv6_cidr_blocks.#:       "0"
      ingress.2541437006.protocol:                 "tcp"
      ingress.2541437006.security_groups.#:        "0"
      ingress.2541437006.self:                     "false"
      ingress.2541437006.to_port:                  "22"
      name:                                        "22_pxy"
      owner_id:                                    <computed>
      revoke_rules_on_delete:                      "false"
      tags.%:                                      "1"
      tags.Name:                                   "22_pxy"
      vpc_id:                                      "${var.vpc_id}"

  + module.subnet.aws_subnet.public
      id:                                          <computed>
      assign_ipv6_address_on_creation:             "false"
      availability_zone:                           "ap-northeast-1a"
      cidr_block:                                  "10.0.1.0/24"
      ipv6_cidr_block:                             <computed>
      ipv6_cidr_block_association_id:              <computed>
      map_public_ip_on_launch:                     "false"
      tags.%:                                      "1"
      tags.Name:                                   "main"

  + module.vpc.aws_vpc.vpc  
      id:                                          <computed>  
      assign_generated_ipv6_cidr_block:            "false"  
      cidr_block:                                  "10.0.0.0/16"
      default_network_acl_id:                      <computed>
      default_route_table_id:                      <computed>
      default_security_group_id:                   <computed>
      dhcp_options_id:                             <computed>  
      enable_classiclink:                          <computed>
      enable_classiclink_dns_support:              <computed>
      enable_dns_hostnames:                        "true"  
      enable_dns_support:                          "true"
      instance_tenancy:                            <computed>
      ipv6_association_id:                         <computed>
      ipv6_cidr_block:                             <computed>
      main_route_table_id:                         <computed>  
      tags.%:                                      "1"  
      tags.Name:                                   "main"  
  
  
Plan: 8 to add, 0 to change, 0 to destroy.  
  
------------------------------------------------------------------------
  
Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
```

### create resources

```
$ make apply

(...)

Plan: 8 to add, 0 to change, 0 to destroy.                                                                                                                               [96/1918]
  
Do you want to perform these actions?  
  Terraform will perform the actions described above.  
  Only 'yes' will be accepted to approve.  
  
  Enter a value: yes  
  
module.ec2.aws_key_pair.tf-key: Creating...  
  fingerprint: "" => "<computed>"  
  key_name:    "" => "tf-key"  
  public_key:  "" => "my ssh key"
module.vpc.aws_vpc.vpc: Creating...  
  assign_generated_ipv6_cidr_block: "" => "false"  
  cidr_block:                       "" => "10.0.0.0/16"  
  default_network_acl_id:           "" => "<computed>"  
  default_route_table_id:           "" => "<computed>"  
  default_security_group_id:        "" => "<computed>"  
  dhcp_options_id:                  "" => "<computed>"  
  enable_classiclink:               "" => "<computed>"  
  enable_classiclink_dns_support:   "" => "<computed>"  
  enable_dns_hostnames:             "" => "true"  
  enable_dns_support:               "" => "true"
  instance_tenancy:                 "" => "<computed>"
  ipv6_association_id:              "" => "<computed>"
  ipv6_cidr_block:                  "" => "<computed>"
  main_route_table_id:              "" => "<computed>"
  tags.%:                           "" => "1"
  tags.Name:                        "" => "main"
module.ec2.aws_key_pair.tf-key: Creation complete after 0s (ID: tf-key)
module.vpc.aws_vpc.vpc: Creation complete after 5s (ID: vpc-d724bfb0)
module.route_table.aws_internet_gateway.public: Creating...
  tags.%:    "0" => "1"
  tags.Name: "" => "main"
  vpc_id:    "" => "vpc-d724bfb0"
module.subnet.aws_subnet.public: Creating...
  assign_ipv6_address_on_creation: "" => "false"
  availability_zone:               "" => "ap-northeast-1a"
  cidr_block:                      "" => "10.0.1.0/24"
  ipv6_cidr_block:                 "" => "<computed>"
  ipv6_cidr_block_association_id:  "" => "<computed>"
  map_public_ip_on_launch:         "" => "false"
  tags.%:                          "" => "1"
  tags.Name:                       "" => "main"
  vpc_id:                          "" => "vpc-d724bfb0"
module.security_group.aws_security_group.22_pxy: Creating...  
  arn:                                   "" => "<computed>"
  description:                           "" => "Permit ssh from proxy server."
  egress.#:                              "" => "1"
  egress.482069346.cidr_blocks.#:        "" => "1"
  egress.482069346.cidr_blocks.0:        "" => "0.0.0.0/0"
  egress.482069346.description:          "" => ""
  egress.482069346.from_port:            "" => "0"
  egress.482069346.ipv6_cidr_blocks.#:   "" => "0"
  egress.482069346.prefix_list_ids.#:    "" => "0"
  egress.482069346.protocol:             "" => "-1"
  egress.482069346.security_groups.#:    "" => "0"
  egress.482069346.self:                 "" => "false"
  egress.482069346.to_port:              "" => "0"
  ingress.#:                             "" => "1"
  ingress.2541437006.cidr_blocks.#:      "" => "1"
  ingress.2541437006.cidr_blocks.0:      "" => "0.0.0.0/0"
  ingress.2541437006.description:        "" => ""
  ingress.2541437006.from_port:          "" => "22"
  ingress.2541437006.ipv6_cidr_blocks.#: "" => "0"
  ingress.2541437006.protocol:           "" => "tcp"
  ingress.2541437006.security_groups.#:  "" => "0"
  ingress.2541437006.self:               "" => "false"
  ingress.2541437006.to_port:            "" => "22"
  name:                                  "" => "22_pxy"
  owner_id:                              "" => "<computed>"
  revoke_rules_on_delete:                "" => "false"
  tags.%:                                "" => "1"
  tags.Name:                             "" => "22_pxy"
  vpc_id:                                "" => "vpc-d724bfb0"
module.route_table.aws_internet_gateway.public: Creation complete after 1s (ID: igw-e8a5e28c)
module.route_table.aws_route_table.public: Creating...                                                                                                                   [18/1918]
  propagating_vgws.#:                         "" => "<computed>"
  route.#:                                    "" => "1"  
  route.4271899941.cidr_block:                "" => "0.0.0.0/0"
  route.4271899941.egress_only_gateway_id:    "" => ""  
  route.4271899941.gateway_id:                "" => "igw-e8a5e28c"
  route.4271899941.instance_id:               "" => ""  
  route.4271899941.ipv6_cidr_block:           "" => ""  
  route.4271899941.nat_gateway_id:            "" => ""  
  route.4271899941.network_interface_id:      "" => ""  
  route.4271899941.vpc_peering_connection_id: "" => ""  
  tags.%:                                     "" => "1"  
  tags.Name:                                  "" => "main"  
  vpc_id:                                     "" => "vpc-d724bfb0"
module.subnet.aws_subnet.public: Creation complete after 1s (ID: subnet-2f83b666)
module.route_table.aws_route_table.public: Creation complete after 1s (ID: rtb-7b7f961d)
module.route_table.aws_route_table_association.public: Creating...
  route_table_id: "" => "rtb-7b7f961d"  
  subnet_id:      "" => "subnet-2f83b666"  
module.route_table.aws_route_table_association.public: Creation complete after 0s (ID: rtbassoc-eb05198d)
module.security_group.aws_security_group.22_pxy: Creation complete after 3s (ID: sg-6de11815)
module.ec2.aws_instance.bastion: Creating...
  ami:                                       "" => "ami-0c11b26d"
  associate_public_ip_address:               "" => "true"
  availability_zone:                         "" => "<computed>"
  ebs_block_device.#:                        "" => "<computed>"
  ebs_optimized:                             "" => "false"
  ephemeral_block_device.#:                  "" => "<computed>"
  get_password_data:                         "" => "false"
  instance_state:                            "" => "<computed>"
  instance_type:                             "" => "t2.micro"
  ipv6_address_count:                        "" => "<computed>"
  ipv6_addresses.#:                          "" => "<computed>"
  key_name:                                  "" => "tf-key"
  monitoring:                                "" => "false"
  network_interface.#:                       "" => "<computed>"
  network_interface_id:                      "" => "<computed>"
  password_data:                             "" => "<computed>"
  placement_group:                           "" => "<computed>"
  primary_network_interface_id:              "" => "<computed>"
  private_dns:                               "" => "<computed>"
  private_ip:                                "" => "<computed>"
  public_dns:                                "" => "<computed>"
  public_ip:                                 "" => "<computed>"
  root_block_device.#:                       "" => "1"
  root_block_device.0.delete_on_termination: "" => "false"
  root_block_device.0.volume_id:             "" => "<computed>"
  root_block_device.0.volume_size:           "" => "50"
  root_block_device.0.volume_type:           "" => "standard"
  security_groups.#:                         "" => "<computed>"
  source_dest_check:                         "" => "true"
  subnet_id:                                 "" => "subnet-2f83b666"
  tags.%:                                    "" => "1"
  tags.Name:                                 "" => "main"
  tenancy:                                   "" => "<computed>"
  volume_tags.%:                             "" => "<computed>"
  vpc_security_group_ids.#:                  "" => "1"
  vpc_security_group_ids.1189096013:         "" => "sg-6de11815"
module.ec2.aws_instance.bastion: Still creating... (10s elapsed)
module.ec2.aws_instance.bastion: Still creating... (20s elapsed)
module.ec2.aws_instance.bastion: Creation complete after 23s (ID: i-00ebded5f73c9a78b)

Apply complete! Resources: 8 added, 0 changed, 0 destroyed.
```

### created resources

![](/images/introduce-terraform/1.png)
![](/images/introduce-terraform/2.png)
![](/images/introduce-terraform/3.png)
![](/images/introduce-terraform/4.png)
![](/images/introduce-terraform/5.png)
![](/images/introduce-terraform/6.png)

### login

```
$ ssh ec2-user@18.182.42.233 -i ~/.ssh/tf-key
The authenticity of host '18.182.42.233 (18.182.42.233)' can't be established.
ECDSA key fingerprint is SHA256:PehqGCGZUEkBSHWVH5LNhUVc5rzbh9Xma+UTiz64oGk.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '18.182.42.233' (ECDSA) to the list of known hosts.

       __|  __|_  )
       _|  (     /   Amazon Linux AMI
      ___|\___|___|

https://aws.amazon.com/amazon-linux-ami/2016.09-release-notes/
25 package(s) needed for security, out of 74 available
Run "sudo yum update" to apply all updates.
Amazon Linux version 2018.03 is available.
[ec2-user@ip-10-0-1-137 ~]$
[ec2-user@ip-10-0-1-137 ~]$
[ec2-user@ip-10-0-1-137 ~]$ whoami
ec2-user
```

succeed login!!

### destroy resources

```
make destroy
terraform destroy
aws_key_pair.tf-key: Refreshing state... (ID: tf-key)
aws_vpc.vpc: Refreshing state... (ID: vpc-d724bfb0)
aws_security_group.22_pxy: Refreshing state... (ID: sg-6de11815)
aws_subnet.public: Refreshing state... (ID: subnet-2f83b666)
aws_internet_gateway.public: Refreshing state... (ID: igw-e8a5e28c)
aws_route_table.public: Refreshing state... (ID: rtb-7b7f961d)
aws_instance.bastion: Refreshing state... (ID: i-00ebded5f73c9a78b)
aws_route_table_association.public: Refreshing state... (ID: rtbassoc-eb05198d)

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  - module.ec2.aws_instance.bastion

  - module.ec2.aws_key_pair.tf-key

  - module.route_table.aws_internet_gateway.public

  - module.route_table.aws_route_table.public

  - module.route_table.aws_route_table_association.public

  - module.security_group.aws_security_group.22_pxy

  - module.subnet.aws_subnet.public

  - module.vpc.aws_vpc.vpc


Plan: 0 to add, 0 to change, 8 to destroy.

Do you really want to destroy?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

module.ec2.aws_instance.bastion: Destroying... (ID: i-00ebded5f73c9a78b)
module.route_table.aws_route_table_association.public: Destroying... (ID: rtbassoc-eb05198d)
module.route_table.aws_route_table_association.public: Destruction complete after 1s
module.route_table.aws_route_table.public: Destroying... (ID: rtb-7b7f961d)
module.route_table.aws_route_table.public: Destruction complete after 8s
module.route_table.aws_internet_gateway.public: Destroying... (ID: igw-e8a5e28c)
module.ec2.aws_instance.bastion: Still destroying... (ID: i-00ebded5f73c9a78b, 10s elapsed)
module.route_table.aws_internet_gateway.public: Still destroying... (ID: igw-e8a5e28c, 10s elapsed)
module.ec2.aws_instance.bastion: Still destroying... (ID: i-00ebded5f73c9a78b, 20s elapsed)
module.route_table.aws_internet_gateway.public: Still destroying... (ID: igw-e8a5e28c, 20s elapsed)
module.ec2.aws_instance.bastion: Still destroying... (ID: i-00ebded5f73c9a78b, 30s elapsed)
module.route_table.aws_internet_gateway.public: Still destroying... (ID: igw-e8a5e28c, 30s elapsed)
module.ec2.aws_instance.bastion: Still destroying... (ID: i-00ebded5f73c9a78b, 40s elapsed)
module.route_table.aws_internet_gateway.public: Still destroying... (ID: igw-e8a5e28c, 40s elapsed)
module.ec2.aws_instance.bastion: Still destroying... (ID: i-00ebded5f73c9a78b, 50s elapsed)
module.route_table.aws_internet_gateway.public: Destruction complete after 47s
module.ec2.aws_instance.bastion: Still destroying... (ID: i-00ebded5f73c9a78b, 1m0s elapsed)
module.ec2.aws_instance.bastion: Destruction complete after 1m2s
module.subnet.aws_subnet.public: Destroying... (ID: subnet-2f83b666)
module.ec2.aws_key_pair.tf-key: Destroying... (ID: tf-key)
module.security_group.aws_security_group.22_pxy: Destroying... (ID: sg-6de11815)
module.ec2.aws_key_pair.tf-key: Destruction complete after 0s
module.security_group.aws_security_group.22_pxy: Destruction complete after 1s
module.subnet.aws_subnet.public: Destruction complete after 1s
module.vpc.aws_vpc.vpc: Destroying... (ID: vpc-d724bfb0)
module.vpc.aws_vpc.vpc: Destruction complete after 0s

Destroy complete! Resources: 8 destroyed.
```

* Simple template will help creating aws resources with terraform
* template must be DRY, and small
