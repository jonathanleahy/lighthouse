locals {
  project_name         = "backoffice-core-bff"
  user_name            = "${local.project_name}-${module.aws_static_parameters.short_region_name}"
  role_name            = "${local.project_name}-role-${module.aws_static_parameters.short_region_name}"
  namespace            = "psm-crm"
  service_account_name = local.project_name

  sns_policy_name        = "${local.project_name}-sns-policy-${module.aws_static_parameters.short_region_name}"
  console_audit_sns_name = "console-audit"
}

