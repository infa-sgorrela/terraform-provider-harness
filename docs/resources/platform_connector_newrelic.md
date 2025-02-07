---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "harness_platform_connector_newrelic Resource - terraform-provider-harness"
subcategory: "Next Gen"
description: |-
  Resource for creating a New Relic connector.
---

# harness_platform_connector_newrelic (Resource)

Resource for creating a New Relic connector.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_id` (String) Account ID of the NewRelic account.
- `api_key_ref` (String) Reference to the Harness secret containing the api key.
- `identifier` (String) Unique identifier of the resource.
- `name` (String) Name of the resource.
- `url` (String) Url of the NewRelic server.

### Optional

- `delegate_selectors` (Set of String) Connect using only the delegates which have these tags.
- `description` (String) Description of the resource.
- `org_id` (String) Unique identifier of the Organization.
- `project_id` (String) Unique identifier of the Project.
- `tags` (Set of String) Tags to associate with the resource. Tags should be in the form `name:value`.

### Read-Only

- `id` (String) The ID of this resource.


