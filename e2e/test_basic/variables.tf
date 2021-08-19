variable "todo_text" {
  type = string
}

variable "todo_user_id" {
  type = string
}

variable "compute_mutation_keys" {
  type = map(string)
}

variable "compute_from_create" {
  type = bool
}

variable "force_replace" {
  type = bool
  default = false
}