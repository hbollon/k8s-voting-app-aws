variable "aws_region" {
  type = string
  default = "eu-west-3"
  description = "The AWS region to be used for the current infrastructure. Defaults to eu-west-3."
}

variable "default_vpc_subnets" {
  type = list(string)
  description = "The default VPC's subnets to be used for the EKS instance. You must provide two subnets."

  validation {
    condition     = length(var.default_vpc_subnets) == 2
    error_message = "The default_vpc_subnets value must be a a list of two subnets."
  }
}