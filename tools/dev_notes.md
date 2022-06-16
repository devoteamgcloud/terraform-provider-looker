# Development notes

## Local development

For local development, I'm using the following configuration.
It's very inconvenient to run the provider using local repo mirror,
as the version number doesn't match the provider binary hash.
For that purpose, you have to use `dev_overrides`.
However, this doesn't allow to test the code completion with `terraform-ls`.
For that reason, I configure both a `filesystem_mirror` and `dev_overrides`.

$ `~/.terraformrc`  
```terraformrc
provider_installation {
  dev_overrides {
    "devoteamgcloud/looker" = "<omit>/git/terraform-provider-looker"
  }

  filesystem_mirror {
    path    = "/Users/vsix/.terraform.d/providers"
    include = ["example.com/*/*", "registry.terraform.io/devoteamgcloud/*"]
  }
  direct {
    exclude = ["example.com/*/*"]
  }
}
```

```shell
mkdir -p ~/.terraform.d/providers/registry.terraform.io/devoteamgcloud/looker/0.0.0-dev/darwin_amd64/

ln -s \
  ./terraform-provider-looker \
  ~/.terraform.d/providers/registry.terraform.io/devoteamgcloud/looker/0.0.0-dev/darwin_amd64/terraform-provider-looker_v0.0.0-dev
  
# TF IaC
trash .terraform .terraform.lock.hcl; terraform init
```

## Building the provider

To build the provider for debug, and remove the overly verbose package paths,
you can use `-trimpath=DIR;DIR2`.

```shell
go build \
  -gcflags=all='-N -l -trimpath={{THIS_DIR}};{{GOMODCACHE}}/github.com/hashicorp' \
  -o terraform-provider-looker
```