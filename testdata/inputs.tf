variable "aws_region" {
  default     = "us-east-1"
  description = "The AWS region in which all resources will be created."
  type        = "string"
}

variable "vpc_id" {
  description = "The id of the specific VPC to retrieve."
  type        = "string"
  default = "1"
}

variable "instance_count" {
  description = "The count of desired instances of EC2."
  type        = "number"
  default     = "2"
}

variable "zones" {
  type = "list"
  description = "The selected zones"
  default = [
    "us-east-1", "us-east-2"
  ]
}