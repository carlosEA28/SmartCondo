variable "aws-region" {
  type        = string
  default     = "sa-east-1"
  description = "Região do AWS"
}

variable "aws_profile" {
  type        = string
  default     = "default"
  description = "Profile do AWS CLI"
}