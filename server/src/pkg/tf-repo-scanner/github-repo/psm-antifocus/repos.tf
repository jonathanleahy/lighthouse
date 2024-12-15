module "module-1" {
  source                 = "github.com/fake-org/tfmod-gh-repo.git?ref=v1.0.4"
  repository_name        = "module-1"
  description            = "Description for module-1"
  required_status_checks = ["module-1/check-1", "module-1/check-2"]
  writers                = ["team-1"]
}

module "module-2" {
  source                 = "github.com/fake-org/tfmod-gh-repo.git?ref=v1.0.4"
  repository_name        = "module-2"
  description            = "Description for module-2"
  required_status_checks = ["module-2/check-1"]
  readers                = ["team-2"]
  writers                = ["team-1", "team-3", "team-4"]
}

module "module-3" {
  source          = "github.com/fake-org/tfmod-gh-repo.git?ref=v1.0.4"
  repository_name = "module-3"
  description     = "Description for module-3"
  readers         = ["team-2"]
  writers         = ["team-1", "team-3", "team-4"]
}

