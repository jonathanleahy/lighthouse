# secrets manager
resource "aws_secretsmanager_secret" "app_secret" {
  name = "/v1/${local.project_name}/global/app-secrets"
}

# secrets value
resource "aws_secretsmanager_secret_version" "app_secret_version" {
  secret_id     = aws_secretsmanager_secret.app_secret.id
  secret_string = jsonencode({"key_id" = aws_iam_access_key.access_key.id, "secret" = aws_iam_access_key.access_key.secret})
}

data "aws_iam_policy_document" "app_secretmanager_policy" {
  statement {
    actions = [
      "secretsmanager:GetResourcePolicy",
      "secretsmanager:GetSecretValue",
      "secretsmanager:DescribeSecret",
      "secretsmanager:ListSecretVersionIds"
    ]
    effect = "Allow"
    resources = [
      aws_secretsmanager_secret.app_secret.arn,
      aws_secretsmanager_secret_version.app_secret_version.arn,
    ]
  }
}

resource "aws_iam_policy" "app_secretmanager_policy" {
  name   = local.secret_policy_name
  policy = data.aws_iam_policy_document.app_secretmanager_policy.json
}

resource "aws_iam_role_policy_attachment" "app_secretmanager_policy_attach" {
  role       = local.role_name
  policy_arn = aws_iam_policy.app_secretmanager_policy.arn
}

