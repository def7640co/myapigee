/*
 * Deployments
 *
 * Manage API proxy and shared flow deployments.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// ApiProxyDeploymentsOrg API proxy deployment details for an organization.
type ApiProxyDeploymentsOrg struct {
	// Name of the environment.
	Environment []ApiProxyDeploymentsOrgEnvironment `json:"environment,omitempty"`
	// Name of the organization.
	Name string `json:"name,omitempty"`
}
