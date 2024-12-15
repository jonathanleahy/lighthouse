module "my-ecr" {
  source               = "github.com/pismo/tf-module-ecr"
  image_tag_mutability = "MUTABLE"

  image_names = [
    "backoffice-core-bff"
  ]

  principals_full_access = [
    "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
  ]

  principals_readonly_access = [
    # DEV-EXT
    "arn:aws:iam::270036487593:root",
    # CITI-NFT
    "arn:aws:iam::752012924079:root",
    # CITI-PROD
    "arn:aws:iam::023441637782:root",
    # CITI-STAG
    "arn:aws:iam::145741235136:root",
    # INTEGRATION
    "arn:aws:iam::459584242408:root",
    # PROD
    "arn:aws:iam::408082092235:root",
    # PROD-ITAU
    "arn:aws:iam::756778449919:root",
    # IND-PROD
    "arn:aws:iam::056132470094:root",
    # PROD-USA
    "arn:aws:iam::008594187592:root",
    # AUS-PROD
    "arn:aws:iam::905418322201:root",
    # IRL-PROD
    "arn:aws:iam::471112918407:root",
    # GETNET-PROD
    "arn:aws:iam::009160028407:root"
  ]
}
