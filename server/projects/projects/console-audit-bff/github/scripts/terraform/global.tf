resource "aws_iam_role" "console_audit_bff_role" {
  name               = local.role_name
  assume_role_policy = data.aws_iam_policy_document.aws_iam_policy_document.json

  tags = {
    project = local.project_name
    env     = var.account
  }
}

# user
resource "aws_iam_user" "user" {
  name = local.user_name
}

# access key
resource "aws_iam_access_key" "access_key" {
  user = aws_iam_user.user.name
}

