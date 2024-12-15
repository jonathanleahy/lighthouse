resource "aws_iam_role" "backoffice_core_bff_role" {
  name               = local.role_name
  assume_role_policy = data.aws_iam_policy_document.aws_iam_policy_document.json
}

# user
resource "aws_iam_user" "user" {
  name = local.user_name
}

# access key
resource "aws_iam_access_key" "access_key" {
  user = aws_iam_user.user.name
}

