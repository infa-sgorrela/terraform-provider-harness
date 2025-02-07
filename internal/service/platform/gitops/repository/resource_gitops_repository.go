package repository

import (
	"context"
	"github.com/antihax/optional"
	hh "github.com/harness/harness-go-sdk/harness/helpers"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGitopsRepositories() *schema.Resource {
	resource := &schema.Resource{
		Description: "Resource for creating Harness Gitops Repositories.",

		CreateContext: resourceGitOpsRepositoryCreate,
		ReadContext:   resourceGitOpsRepositoryRead,
		UpdateContext: resourceGitOpsRepositoryUpdate,
		DeleteContext: resourceGitOpsRepositoryDelete,
		Importer:      helpers.GitopsAgentResourceImporter,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Account Identifier for the Repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project_id": {
				Description: "Project Identifier for the  Repository.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"org_id": {
				Description: "Organization Identifier for the Repository.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"agent_id": {
				Description: "Agent identifier for the Repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"identifier": {
				Description: "Identifier of the Repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"repo": {
				Description: "Repo Details holding application configurations",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"repo": {
							Description: "URL to the remote repository.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"username": {
							Description: "user name used for authenticating at the remote repository.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"password": {
							Description: "password or PAT used for authenticating at the remote repository.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"ssh_private_key": {
							Description: "the PEM data for authenticating at the repo server. Only used with Git repos.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"insecure_ignore_host_key": {
							Description: "InsecureIgnoreHostKey should not be used anymore, Insecure is favoured Used only for Git repos.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"insecure": {
							Description: "specifies whether the connection to the repository ignores any errors when verifying TLS certificates or SSH host keys.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"enable_lfs": {
							Description: " whether git-lfs support should be enabled for this repo. Only valid for Git repositories.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"tls_client_cert_data": {
							Description: "certificate in PEM format for authenticating at the repo server.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"tls_client_cert_key": {
							Description: "private key in PEM format for authenticating at the repo server.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"type_": {
							Description: "Type specifies the type of the repo. Can be either \"git\" or \"helm. \"git\" is assumed if empty or absent.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"name": {
							Description: "name to be used for this repo. Only used with Helm repos.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"inherited_creds": {
							Description: "Whether credentials were inherited from a credential set.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"enable_oci": {
							Description: "whether helm-oci support should be enabled for this repo",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"github_app_private_key": {
							Description: "Github App Private Key PEM data.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"github_app_id": {
							Description: "the ID of the GitHub app used to access the repo.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"github_app_installation_id": {
							Description: " the installation ID of the GitHub App used to access the repo.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"github_app_enterprise_base_url": {
							Description: "the base URL of GitHub Enterprise installation. If empty will default to https://api.github.com",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"proxy": {
							Description: "the HTTP/HTTPS proxy used to access the repo.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"project": {
							Description: "Reference between project and repository that allow you automatically to be added as item inside SourceRepos project entity.",
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"connection_type": {
							Description: "Identifies the authentication method used to connect to the repository",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"upsert": {
				Description: "Whether to create in upsert mode.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"creds_only": {
				Description: "Whether to operate on credential set instead of repository.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"query_repo": {
				Description: "Repo to Query.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"query_project": {
				Description: "Project to Query for Repo.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"query_force_refresh": {
				Description: "Force refresh query for Repo.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"update_mask": {
				Description: "Update mask of the Repository.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"paths": {
							Description: "The set of field mask paths.",
							Optional:    true,
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
	return resource
}

func resourceGitOpsRepositoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)
	ctx = context.WithValue(ctx, nextgen.ContextAccessToken, hh.EnvVars.BearerToken.Get())
	var agentIdentifier, accountIdentifier, orgIdentifier, projectIdentifier, identifier string
	accountIdentifier = c.AccountId
	if attr, ok := d.GetOk("agent_id"); ok {
		agentIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("org_id"); ok {
		orgIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("project_id"); ok {
		projectIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("identifier"); ok {
		identifier = attr.(string)
	}

	createRepoRequest := buildCreateRepoRequest(d)
	resp, httpResp, err := c.RepositoriesApiService.AgentRepositoryServiceCreateRepository(ctx, createRepoRequest, agentIdentifier, &nextgen.RepositoriesApiAgentRepositoryServiceCreateRepositoryOpts{
		AccountIdentifier: optional.NewString(accountIdentifier),
		OrgIdentifier:     optional.NewString(orgIdentifier),
		ProjectIdentifier: optional.NewString(projectIdentifier),
		Identifier:        optional.NewString(identifier),
	})

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}
	// Soft delete lookup error handling
	// https://harness.atlassian.net/browse/PL-23765
	if resp.Repository == nil {
		d.SetId("")
		d.MarkNewResource()
		return nil
	}
	setRepositoryDetails(d, &resp)
	return nil
}

func resourceGitOpsRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)
	var orgIdentifier, projectIdentifier, agentIdentifier, identifier, queryRepo, queryProject string
	var queryForceRefresh bool
	if attr, ok := d.GetOk("org_id"); ok {
		orgIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("project_id"); ok {
		projectIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("agent_id"); ok {
		agentIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("identifier"); ok {
		identifier = attr.(string)
	}
	if attr, ok := d.GetOk("query_repo"); ok {
		queryRepo = attr.(string)
	}
	if attr, ok := d.GetOk("query_project"); ok {
		queryProject = attr.(string)
	}
	if attr, ok := d.GetOk("query_force_refresh"); ok {
		queryForceRefresh = attr.(bool)
	}
	resp, httpResp, err := c.RepositoriesApiService.AgentRepositoryServiceGet(ctx, agentIdentifier, identifier, c.AccountId, &nextgen.RepositoriesApiAgentRepositoryServiceGetOpts{
		OrgIdentifier:     optional.NewString(orgIdentifier),
		ProjectIdentifier: optional.NewString(projectIdentifier),
		QueryRepo:         optional.NewString(queryRepo),
		QueryForceRefresh: optional.NewBool(queryForceRefresh),
		QueryProject:      optional.NewString(queryProject),
	})

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}
	// Soft delete lookup error handling
	// https://harness.atlassian.net/browse/PL-23765
	if resp.Repository == nil {
		d.SetId("")
		d.MarkNewResource()
		return nil
	}
	setRepositoryDetails(d, &resp)
	return nil

}

func resourceGitOpsRepositoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)
	var orgIdentifier, projectIdentifier, agentIdentifier, identifier string
	if attr, ok := d.GetOk("org_id"); ok {
		orgIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("project_id"); ok {
		projectIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("agent_id"); ok {
		agentIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("identifier"); ok {
		identifier = attr.(string)
	}
	updateRepoRequest := buildUpdateRepoRequest(d)
	resp, httpResp, err := c.RepositoriesApiService.AgentRepositoryServiceUpdateRepository(ctx, updateRepoRequest, agentIdentifier, identifier, &nextgen.RepositoriesApiAgentRepositoryServiceUpdateRepositoryOpts{
		AccountIdentifier: optional.NewString(c.AccountId),
		OrgIdentifier:     optional.NewString(orgIdentifier),
		ProjectIdentifier: optional.NewString(projectIdentifier),
	})

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}
	// Soft delete lookup error handling
	// https://harness.atlassian.net/browse/PL-23765
	if resp.Repository == nil {
		d.SetId("")
		d.MarkNewResource()
		return nil
	}
	setRepositoryDetails(d, &resp)
	return nil
}

func resourceGitOpsRepositoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ctx := meta.(*internal.Session).GetPlatformClientWithContext(ctx)
	var orgIdentifier, projectIdentifier, agentIdentifier, identifier, queryRepo, queryProject string
	var queryForceRefresh bool
	if attr, ok := d.GetOk("org_id"); ok {
		orgIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("project_id"); ok {
		projectIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("agent_id"); ok {
		agentIdentifier = attr.(string)
	}
	if attr, ok := d.GetOk("identifier"); ok {
		identifier = attr.(string)
	}
	if attr, ok := d.GetOk("query_repo"); ok {
		queryRepo = attr.(string)
	}
	if attr, ok := d.GetOk("query_project"); ok {
		queryProject = attr.(string)
	}
	if attr, ok := d.GetOk("query_force_refresh"); ok {
		queryForceRefresh = attr.(bool)
	}
	_, httpResp, err := c.RepositoriesApiService.AgentRepositoryServiceDeleteRepository(ctx, agentIdentifier, identifier, &nextgen.RepositoriesApiAgentRepositoryServiceDeleteRepositoryOpts{
		AccountIdentifier: optional.NewString(c.AccountId),
		OrgIdentifier:     optional.NewString(orgIdentifier),
		ProjectIdentifier: optional.NewString(projectIdentifier),
		QueryRepo:         optional.NewString(queryRepo),
		QueryForceRefresh: optional.NewBool(queryForceRefresh),
		QueryProject:      optional.NewString(queryProject),
	})
	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}
	// Soft delete lookup error handling
	// https://harness.atlassian.net/browse/PL-23765
	return nil
}

func buildUpdateRepoRequest(d *schema.ResourceData) nextgen.RepositoriesRepoUpdateRequest {
	var updateMask map[string]interface{}
	if attr, ok := d.GetOk("update_mask"); ok {
		if len(attr.([]interface{})) > 0 {
			updateMask = attr.([]interface{})[0].(map[string]interface{})
		}
	}
	var updateMaskPath []string
	if updateMask != nil && updateMask["paths"] != nil && len(updateMask["paths"].([]interface{})) > 0 {
		for _, v := range updateMask["paths"].([]interface{}) {
			updateMaskPath = append(updateMaskPath, v.(string))
		}
	}
	return nextgen.RepositoriesRepoUpdateRequest{
		Repo: buildRepo(d),
		UpdateMask: &nextgen.ProtobufFieldMask{
			Paths: updateMaskPath,
		},
	}
}

func buildCreateRepoRequest(d *schema.ResourceData) nextgen.RepositoriesRepoCreateRequest {
	var upsert, credsOnly bool
	if attr, ok := d.GetOk("upsert"); ok {
		upsert = attr.(bool)
	}
	if attr, ok := d.GetOk("creds_only"); ok {
		credsOnly = attr.(bool)
	}
	return nextgen.RepositoriesRepoCreateRequest{
		Upsert:    upsert,
		CredsOnly: credsOnly,
		Repo:      buildRepo(d),
	}
}

func buildRepo(d *schema.ResourceData) *nextgen.RepositoriesRepository {
	var repoObj = nextgen.RepositoriesRepository{}
	if attr, ok := d.GetOk("repo"); ok {
		if attr != nil && len(attr.([]interface{})) > 0 {
			var repo = attr.([]interface{})[0].(map[string]interface{})
			if repo["repo"] != nil {
				repoObj.Repo = repo["repo"].(string)
			}
			if repo["username"] != nil {
				repoObj.Username = repo["username"].(string)
			}
			if repo["password"] != nil {
				repoObj.Password = repo["password"].(string)
			}
			if repo["ssh_private_key"] != nil {
				repoObj.SshPrivateKey = repo["ssh_private_key"].(string)
			}
			if repo["insecure_ignore_host_key"] != nil {
				repoObj.InsecureIgnoreHostKey = repo["insecure_ignore_host_key"].(bool)
			}
			if repo["insecure"] != nil {
				repoObj.Insecure = repo["insecure"].(bool)
			}
			if repo["enable_lfs"] != nil {
				repoObj.EnableLfs = repo["enable_lfs"].(bool)
			}
			if repo["tls_client_cert_data"] != nil {
				repoObj.TlsClientCertData = repo["tls_client_cert_data"].(string)
			}
			if repo["tls_client_cert_key"] != nil {
				repoObj.TlsClientCertKey = repo["tls_client_cert_key"].(string)
			}
			if repo["type_"] != nil {
				repoObj.Type_ = repo["type_"].(string)
			}
			if repo["name"] != nil {
				repoObj.Name = repo["name"].(string)
			}
			if repo["inherited_creds"] != nil {
				repoObj.InheritedCreds = repo["inherited_creds"].(bool)
			}
			if repo["enable_oci"] != nil {
				repoObj.EnableOCI = repo["enable_oci"].(bool)
			}
			if repo["github_app_private_key"] != nil {
				repoObj.GithubAppPrivateKey = repo["github_app_private_key"].(string)
			}
			if repo["github_app_id"] != nil {
				repoObj.GithubAppID = repo["github_app_id"].(string)
			}
			if repo["github_app_installation_id"] != nil {
				repoObj.GithubAppInstallationID = repo["github_app_installation_id"].(string)
			}
			if repo["github_app_enterprise_base_url"] != nil {
				repoObj.GithubAppEnterpriseBaseUrl = repo["github_app_enterprise_base_url"].(string)
			}
			if repo["proxy"] != nil {
				repoObj.Proxy = repo["proxy"].(string)
			}
			if repo["project"] != nil {
				repoObj.Project = repo["project"].(string)
			}
			if repo["connection_type"] != nil {
				repoObj.ConnectionType = repo["connection_type"].(string)
			}
		}
	}
	return &repoObj
}

func setRepositoryDetails(d *schema.ResourceData, repo *nextgen.Servicev1Repository) {
	d.SetId(repo.Identifier)
	d.Set("account_id", repo.AccountIdentifier)
	d.Set("org_id", repo.OrgIdentifier)
	d.Set("project_id", repo.ProjectIdentifier)
	d.Set("agent_id", repo.AgentIdentifier)
	d.Set("identifier", repo.Identifier)
	if repo.Repository != nil {
		repoList := []interface{}{}
		repoO := map[string]interface{}{}
		repoO["repo"] = repo.Repository.Repo
		repoO["username"] = repo.Repository.Username
		repoO["password"] = repo.Repository.Password
		repoO["ssh_private_key"] = repo.Repository.SshPrivateKey
		repoO["insecure_ignore_host_key"] = repo.Repository.InsecureIgnoreHostKey
		repoO["insecure"] = repo.Repository.Insecure
		repoO["enable_lfs"] = repo.Repository.EnableLfs
		repoO["tls_client_cert_data"] = repo.Repository.TlsClientCertData
		repoO["tls_client_cert_key"] = repo.Repository.TlsClientCertKey
		repoO["type_"] = repo.Repository.Type_
		repoO["name"] = repo.Repository.Name
		repoO["inherited_creds"] = repo.Repository.InheritedCreds
		repoO["enable_oci"] = repo.Repository.EnableOCI
		repoO["github_app_private_key"] = repo.Repository.GithubAppPrivateKey
		repoO["github_app_id"] = repo.Repository.GithubAppID
		repoO["github_app_installation_id"] = repo.Repository.GithubAppInstallationID
		repoO["github_app_enterprise_base_url"] = repo.Repository.GithubAppEnterpriseBaseUrl
		repoO["proxy"] = repo.Repository.Proxy
		repoO["project"] = repo.Repository.Project
		repoO["connection_type"] = repo.Repository.ConnectionType

		repoList = append(repoList, repoO)
		d.Set("repo", repoList)
	}
}
