variable "subnet_a" {
  type        = string
  description = "Value of the name tag for the app subnet in AZ a"
  default     = "Web_Dev_aza_net"
}

variable "subnet_b" {
  type        = string
  description = "Value of the name tag for the app subnet in AZ b"
  default     = "Web_Dev_azb_net"
}

variable "region" {
    type = string
    default = "ca-central-1"
}

variable "gold_ip" {
  type = string
}

variable "auth_secret" {
  type = string
  sensitive = true
}
