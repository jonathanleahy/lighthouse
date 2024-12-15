locals {
  project_name         = "console-audit-bff"
  user_name            = "${local.project_name}-${module.aws_static_parameters.short_region_name}"
  role_name            = "${local.project_name}-role-${module.aws_static_parameters.short_region_name}"
  namespace            = "psm-console"
  service_account_name = local.project_name

  sns_policy_name        = "${local.project_name}-sns-policy-${module.aws_static_parameters.short_region_name}"
  console_audit_sns_name = "console-audit"
  secret_policy_name     = "${local.project_name}-secretmanager-policy-${module.aws_static_parameters.short_region_name}"
}

