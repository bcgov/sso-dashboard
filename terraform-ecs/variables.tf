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
  type    = string
  default = "ca-central-1"
}

variable "auth_secret" {
  type        = string
  description = "Authentication secret to use loki API"
  sensitive   = true
}

variable "loki_read_cpu" {
  type        = number
  description = "CPU as vCPU, e.g. 1000 = 1cpu"
  default     = 256
}

variable "loki_write_cpu" {
  type        = number
  description = "CPU as vCPU, e.g. 1000 = 1cpu"
  default     = 256
}

variable "loki_read_memory" {
  type        = number
  description = "Memory in Mb"
  default     = 512
}

variable "loki_write_memory" {
  type        = number
  description = "Memory in Mb"
  default     = 512
}

variable "bucket_name" {
  type    = string
}

variable "loki_tag" {
  type = string
  default = "dev"
}
