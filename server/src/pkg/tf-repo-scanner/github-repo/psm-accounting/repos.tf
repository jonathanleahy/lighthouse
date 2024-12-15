module "module-11" {
  source                 = "github.com/fake-org/tfmod-gh-repo.git?ref=v1.0.4"
  repository_name        = "module-11"
  description            = "Description for module-1"
  required_status_checks = ["module-11/check-1", "module-1/check-2"]
  writers                = ["team-1"]
}

module "module-22" {
  source                 = "github.com/fake-org/tfmod-gh-repo.git?ref=v1.0.4"
  repository_name        = "module-22"
  description            = "Description for module-22"
  required_status_checks = ["module-22/check-1"]
  readers                = ["team-2"]
  writers                = ["team-1", "team-3", "team-4"]
}

