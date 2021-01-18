variable "aws_region" {
  description = "The AWS region in which all resources will be created."
  type        = string
  default     = "us-east-1"
}

variable "vpc_id" {
  description = "The id of the specific VPC to retrieve."
  type        = string
  default     = "1"
}

variable "instance_count" {
  description = "The count of desired instances of EC2."
  type        = number
  default     = "2"
}

variable "zones" {
  description = "The selected zones for deployment."
  type        = list(string)
  default = [
    "us-east-1", "us-east-2"
  ]
}