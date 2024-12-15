output "aws_region_name" {
  value = data.aws_region.current.name
  description = "aws region name"
}

output "aws_region_endpoint" {
  value = data.aws_region.current.endpoint
  description = "aws region endpoint"
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
  description = "account id"
}

output "caller_arn" {
  value = data.aws_caller_identity.current.arn
  description = "aws caller identity arn"
}

output "caller_user" {
  value = data.aws_caller_identity.current.user_id
  description = "aws caller identity user"
}
