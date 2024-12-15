resource "aws_iam_policy" "backoffice_core_bff_sns_policy" {
  name   = local.sns_policy_name
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "sns",
            "Effect": "Allow",
            "Action": [
                "sns:Publish",
                "sns:SetTopicAttributes",
                "sns:ListSubscriptionsByTopic",
                "sns:GetTopicAttributes",
                "sns:Receive",
                "sns:Subscribe"
            ],
            "Resource": [
                "arn:aws:sns:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${local.console_audit_sns_name}"
            ]
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "app_sns_policy_attach" {
  role       = aws_iam_role.backoffice_core_bff_role.name
  policy_arn = aws_iam_policy.backoffice_core_bff_sns_policy.arn
}

resource "aws_iam_user_policy_attachment" "sns_policy" {
  user       = aws_iam_user.user.name
  policy_arn = aws_iam_policy.backoffice_core_bff_sns_policy.arn
}

