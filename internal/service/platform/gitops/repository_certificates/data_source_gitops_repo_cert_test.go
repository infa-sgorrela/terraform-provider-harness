package repository_certificates_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/harness/harness-go-sdk/harness/utils"
	"github.com/harness/terraform-provider-harness/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitOpsRepoCert(t *testing.T) {

	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(6))
	name := id
	agentId := "account.terraformagent1"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")
	clusterName := id
	resourceName := "data.harness_platform_gitops_repo_cert.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepoCert(id, accountId, name, agentId, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "org_id", id),
					resource.TestCheckResourceAttr(resourceName, "project_id", id),
				),
			},
		},
	})
}

func testAccDataSourceRepoCert(id string, accountId string, name string, agentId string, clusterName string) string {
	return fmt.Sprintf(`
		resource "harness_platform_organization" "test" {
			identifier = "%[1]s"
			name = "%[3]s"
		}

		resource "harness_platform_project" "test" {
			identifier = "%[1]s"
			name = "%[3]s"
			org_id = harness_platform_organization.test.id
		}

		resource "harness_platform_gitops_repo_cert" "test" {
			account_id = "%[2]s"
			project_id = harness_platform_project.test.id
			org_id = harness_platform_organization.test.id
			agent_id = "%[4]s"

 			request {
				upsert = true
				certificates {
					metadata {

					}
					items {
						server_name = "rajRepoCert"
						cert_type = "https"
						cert_data = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZsakNDQTM0Q0NRRDlweEZoeFdmb0tUQU5CZ2txaGtpRzl3MEJBUXNGQURDQmpERUxNQWtHQTFVRUJoTUMKU1U0eEZEQVNCZ05WQkFnTUMwMUJTRUZTUVZOSVZGSkJNUTB3Q3dZRFZRUUhEQVJRVlU1Rk1SQXdEZ1lEVlFRSwpEQWRJUVZKT1JWTlRNUTh3RFFZRFZRUUxEQVpFUlZaUFVGTXhFVEFQQmdOVkJBTU1DR0pzWVdoaWJHRm9NU0l3CklBWUpLb1pJaHZjTkFRa0JGaE55WW1GMmFYTnJZWEpBWjIxaGFXd3VZMjl0TUI0WERUSXlNVEF5TVRJeE5ESTAKTjFvWERUSXpNVEF5TVRJeE5ESTBOMW93Z1l3eEN6QUpCZ05WQkFZVEFrbE9NUlF3RWdZRFZRUUlEQXROUVVoQgpVa0ZUU0ZSU1FURU5NQXNHQTFVRUJ3d0VVRlZPUlRFUU1BNEdBMVVFQ2d3SFNFRlNUa1ZUVXpFUE1BMEdBMVVFCkN3d0dSRVZXVDFCVE1SRXdEd1lEVlFRRERBaGliR0ZvWW14aGFERWlNQ0FHQ1NxR1NJYjNEUUVKQVJZVGNtSmgKZG1semEyRnlRR2R0WVdsc0xtTnZiVENDQWlJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dJUEFEQ0NBZ29DZ2dJQgpBT0IrdFZUUmxTcGdNMndaYWtnc2FBN0lTczZ3VmR1YXBpYlV0MXZNZEtjZklVTDJXNUhzK1pRWW9nZHcyK1VmCmVWR2ZZTUZyS3prb2hkWm5Ra2lxTjN3UkVMUGs2bHNqbDd0ZlFrcFpTODcvSklGZjhobytRZWl4am9rdnhrTGkKL0U0UldaR3lxQU1TQ2syNWcvNVpxUzhjdklLaEphVUx5Mm9KaWN3Zm5QUUZlbWpPdG5PV2tTS09ZWk1wY0cxYwpLVzRRQlMwTjdsemJuR05PVjYvcUJ5YkdzQytJa3hORmNBZUNhMHVxMFJSNWNVZThvNm9WS3M3SGwwUVFQTU9qCnpEaEtJc01rWUlETS8zNFEvaG4rQVB2eVNuL3M3UFpVVWNIcE9zUmVTcE1HdVQwcjJ4SmVCRGZXbWVYck1YRDMKRWpwOUlxVGp3TktOcnhZR2RFUytQSDE1TUpIMEk4R2lvV2RyTTE3T3lKbk0ray9odmNyNm1lRmduWXVlb21PVQpsQ0dPUkkvRGM0bExuVGZFdldVUVg1eGNaNEUxdkU4anUyR0srRGVQM3RuT0NJY2ZrZ1hVR0tXRTlCNFMyTHpECnVOSW1hVllNb3ZKM3VYTjQ1MHFiY1ZhMHNhZWZyQjA3TzBwbmo5Mm0vSzJLUGpsRzRYY2hzdU11ZmpyQitOZDgKTHhGL1FjMkZFVHc0bnZGNDRUTGoyTmlrTzdML1JwZVdUMGVTeVZJZUZWdzQ1WmJrYkQraGpUV1lmdHVCOEF5VwpPemhDZ1RhODhDbzVEbzd5UVkzdFNmYm1tWCtBd0hYRDdzS0tzRk03aG9hWEF4Mjl0TGVma04wWGpIcUpXUWgyCkdtTDhCUW1tVVY2dVZkMGRGdUlBSmtKVE5RQXNwVSs0Q2hmS1J2bkdVODVuQWdNQkFBRXdEUVlKS29aSWh2Y04KQVFFTEJRQURnZ0lCQUI2b01MSzFTSFVMbzJ6ditwWk9pbDRlcm95VEcxUUU4Q04wU1haRWRESy9SbkNFaCtNeQpETldOSUxrTG1MZnQydGUzbnFXN3VjVlZhcG9IWFBOS09HRHR2UjlBVEFiSEN3UEF4TUxkT21CZ0R1dFI2dHpBCmNHNklZemU0VVp0NTJ0dTJtRk42REp3aXVoNWRPcUlKK0pYK0tUS2tDYktPZlNOdXhoUlFZVFZ0dHR1NklTQWwKK2ppR0ozYkdSUktkaUxKUWIvVWFubHhValZmVW1wdUIvV2RmZzBES0oxUkpDdW1XbkJFbDF5RXUyak9SUnZqKwpQTEdtNFhOU1A3YVhwSzhBS0p2ZkJaMmxiNGhIa2JrQUxiVys0OTMvWkZRV1orZ0NRQ1MvUEtaWThyYkVBTkZPCnpWMVdzTms0WWdSaE5rYllFNkV2eDBzNXV1MnovZTZIeVRrR25tdXdrSVNIcjJsUGRrS2tFYjJrbTZsUG4zYnoKRkNpNFl4Vy8vYzBJeDdZaHQ5MnNDLzAyUisyd0cvSTlWMW1zLzFhOHAzUGhvOHB6aEdFMGdzcENFMEIyRHd6aQovajdudEVCZ000SGFmcEJGMDRETEYyUjNkQTJUSmNWQWpXdlZLUkhTbkhvNmlMaklvaXppQUxobHlDbHFQWVhnClI5c2JJWWtBamg3SnRTaVlkTGVjeHVNaU5XU0xUd1R5dk56ck5LQ3hvWlhiSE14N1dlY0VTTUZIM3NQdkJLOVAKSUFuVGVoRTFrU1FLZmtBVmt2TWVaQnRaU2JIdStORU9MRmJsaUhwallGMHhmQjdOTnV1Vy9uYk9XeXA0N1NQNgpKTysvZ2llMVhRb2thVloyUG9sR3pmVHpTV1c5V1NROUpEMjRjUXBFczIzVTBsWlV2UWsvcXdyRQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t"
					}
				}
			}
		}

		data "harness_platform_gitops_repo_cert" "test" {
			account_id = "%[2]s"
			project_id = harness_platform_project.test.id
			org_id = harness_platform_organization.test.id
			agent_id = "%[4]s"
		}
`, id, accountId, name, agentId, clusterName)
}
