module "my-ecr" {
  source                 = "github.com/pismo/tf-module-ecr?ref=1.0.0"
  image_names            = ["console-audit-bff"]
  image_tag_mutability   = "MUTABLE"
  protected_tags         = ["production"]
  principals_full_access = [
    "arn:aws:iam::${var.account_id}:root",
    "arn:aws:iam::231406070346:root"

  ]
  principals_readonly_access = [
    // pismolabs account
    "arn:aws:iam::${var.account_id}:root",
    // prod account
    "arn:aws:iam::408082092235:root",
    // itau account
    "arn:aws:iam::756778449919:root",
    // integration account
    "arn:aws:iam::459584242408:root",
    //citi stag
    "arn:aws:iam::145741235136:root",
    // shared services account
    "arn:aws:iam::231406070346:root",
    // ind-prod services account
    "arn:aws:iam::056132470094:root",
    //prod-usa
    "arn:aws:iam::008594187592:root",
    //aus-prod
    "arn:aws:iam::905418322201:root",
    //getnet-prod
    "arn:aws:iam::009160028407:root",
    //irl-prod
    "arn:aws:iam::471112918407:root"
  ]
}
