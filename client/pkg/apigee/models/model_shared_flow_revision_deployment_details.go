/*
 * Deployments
 *
 * Manage API proxy and shared flow deployments.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// SharedFlowRevisionDeploymentDetails Shared flow revision deployment details in an environment.
type SharedFlowRevisionDeploymentDetails struct {
	// Name of the environment.
	Environment string `json:"environment,omitempty"`
	// Revision of the shared flow.
	Name string `json:"name,omitempty"`
	// Name of the organization.
	Organization string `json:"organization,omitempty"`
	// Revision of the shared flow.
	Revision string `json:"revision,omitempty"`
	// Used by Apigee support to identify servers that support the API proxy or shared flow deployment.
	Server []ApiProxyRevisionDeploymentsServer `json:"server,omitempty"`
	// Name of the shared flow.
	SharedFlow string `json:"sharedFlow,omitempty"`
	// Deployment status, such as `deployed` or `undeployed`.
	State string `json:"state,omitempty"`
}