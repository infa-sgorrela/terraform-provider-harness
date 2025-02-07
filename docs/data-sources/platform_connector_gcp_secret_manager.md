---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "harness_platform_connector_gcp_secret_manager Data Source - terraform-provider-harness"
subcategory: "Next Gen"
description: |-
  Datasource for looking up GCP Secret Manager connector.
---

# harness_platform_connector_gcp_secret_manager (Data Source)

Datasource for looking up GCP Secret Manager connector.

## Example Usage

```terraform
data "harness_platform_connector_gcp_secret_manager" "example" {
  identifier = "identifier"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `identifier` (String) Unique identifier of the resource.
- `name` (String) Name of the resource.
- `org_id` (String) Unique identifier of the Organization.
- `project_id` (String) Unique identifier of the Project.

### Read-Only

- `credentials_ref` (String) Reference to the secret containing credentials of IAM service account for Google Secret Manager.
- `delegate_selectors` (Set of String) Connect using only the delegates which have these tags.
- `description` (String) Description of the resource.
- `id` (String) The ID of this resource.
- `is_default` (Boolean) Indicative if this is default Secret manager for secrets.
- `tags` (Set of String) Tags to associate with the resource. Tags should be in the form `name:value`.


