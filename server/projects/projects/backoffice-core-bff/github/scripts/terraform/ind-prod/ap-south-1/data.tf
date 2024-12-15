data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "aws_iam_policy_document" {
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    effect  = "Allow"

    condition {
      test     = "ForAnyValue:StringLike"
      variable = "${replace(module.aws_static_parameters.aws_iam_openid_connect_provider.issuer_url, "https://", "")}:sub"
      values   = ["system:serviceaccount:${local.namespace}:${local.service_account_name}"]
    }

    principals {
      identifiers = [
        module.aws_static_parameters.aws_iam_openid_connect_provider.arn
      ]
      type = "Federated"
    }
  }

  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"

    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/role-${data.aws_caller_identity.current.account_id}-external-secrets-${data.aws_region.current.name}"
      ]
    }
  }
}
