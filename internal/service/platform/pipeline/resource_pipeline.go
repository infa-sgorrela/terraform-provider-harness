package pipeline

import (
	"context"
	"net/http"

	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePipeline() *schema.Resource {
	resource := &schema.Resource{
		Description: "Resource for creating a Harness pipeline.",

		ReadContext:   resourcePipelineRead,
		UpdateContext: resourcePipelineCreateOrUpdate,
		DeleteContext: resourcePipelineDelete,
		CreateContext: resourcePipelineCreateOrUpdate,
		Importer:      helpers.ProjectResourceImporter,

		Schema: map[string]*schema.Schema{
			"yaml": {
				Description: "YAML of the pipeline.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}

	helpers.SetProjectLevelResourceSchema(resource.Schema)

	return resource
}

func resourcePipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	id := d.Id()

	resp, httpResp, err := c.PipelinesApi.GetPipeline(ctx,
		c.AccountId,
		d.Get("org_id").(string),
		d.Get("project_id").(string),
		id,
		&nextgen.PipelinesApiGetPipelineOpts{},
	)

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readPipeline(d, resp.Data)

	return nil
}

func resourcePipelineCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	var err error
	var httpResp *http.Response
	id := d.Id()
	pipeline := buildPipeline(d)

	if id == "" {
		_, httpResp, err = c.PipelinesApi.PostPipeline(ctx, pipeline.Yaml, c.AccountId, pipeline.OrgIdentifier, pipeline.ProjectIdentifier, &nextgen.PipelinesApiPostPipelineOpts{})
	} else {
		_, httpResp, err = c.PipelinesApi.UpdatePipelineV2(ctx, pipeline.Yaml, c.AccountId, pipeline.OrgIdentifier, pipeline.ProjectIdentifier, id, &nextgen.PipelinesApiUpdatePipelineV2Opts{})
	}

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	// The create/update methods don't return the yaml in the response, so we need to query for it again.
	resp, httpResp, err := c.PipelinesApi.GetPipeline(ctx, c.AccountId, pipeline.OrgIdentifier, pipeline.ProjectIdentifier, pipeline.Identifier, &nextgen.PipelinesApiGetPipelineOpts{})
	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readPipeline(d, resp.Data)

	return nil
}

func resourcePipelineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)

	pipeline := buildPipeline(d)

	_, httpResp, err := c.PipelinesApi.DeletePipeline(ctx, c.AccountId, pipeline.OrgIdentifier, pipeline.ProjectIdentifier, pipeline.Identifier, &nextgen.PipelinesApiDeletePipelineOpts{})
	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	return nil
}

// Build PipelineYAML object from stored pipeline yaml
func buildPipeline(d *schema.ResourceData) *nextgen.Pipeline {
	return &nextgen.Pipeline{
		Identifier:        d.Get("identifier").(string),
		Name:              d.Get("name").(string),
		OrgIdentifier:     d.Get("org_id").(string),
		ProjectIdentifier: d.Get("project_id").(string),
		Yaml:              d.Get("yaml").(string),
	}
}

// Read response from API out to the stored identifiers
func readPipeline(d *schema.ResourceData, pipeline *nextgen.PmsPipelineResponse) {
	d.SetId(pipeline.PipelineData.Pipeline.Identifier)
	d.Set("identifier", pipeline.PipelineData.Pipeline.Identifier)
	d.Set("name", pipeline.PipelineData.Pipeline.Name)
	d.Set("org_id", pipeline.PipelineData.Pipeline.OrgIdentifier)
	d.Set("project_id", pipeline.PipelineData.Pipeline.ProjectIdentifier)
	d.Set("yaml", pipeline.YamlPipeline)
}
