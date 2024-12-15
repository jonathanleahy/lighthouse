# secrets manager
resource "aws_secretsmanager_secret" "app_secret" {
  name = "/v1/${local.project_name}/global/app-secrets"
}

# secrets value
resource "aws_secretsmanager_secret_version" "app_secret_version" {
  secret_id     = aws_secretsmanager_secret.app_secret.id
  secret_string = jsonencode({"key_id" = aws_iam_access_key.access_key.id, "secret" = aws_iam_access_key.access_key.secret})
}
