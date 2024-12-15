output "aws_region_name" {
  value = data.aws_region.current.name
  description = "Defines AWS region name"
}

output "aws_region_endpoint" {
  value = data.aws_region.current.endpoint
  description = "Defines AWS region endpoint"
}

output "account_id" {
  value = data.aws_caller_identity.current.account_id
  description = "Defines AWS account id"
}

output "caller_arn" {
  value = data.aws_caller_identity.current.arn
  description = "Defines AWS arn caller"
}

output "caller_user" {
  value = data.aws_caller_identity.current.user_id
  description = "Defines AWS arn user"
}
