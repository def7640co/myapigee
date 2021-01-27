/*
 * Deployments
 *
 * Manage API proxy and shared flow deployments.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// ApiProxyRevisionDeployments API proxy revision deployments.
type ApiProxyRevisionDeployments struct {
	// Name of the API proxy
	APIProxy interface{} `json:"aPIProxy,omitempty"`
	// API proxy deployment details in the environment.
	Environment []ApiProxyRevisionDeploymentsEnvironment `json:"environment,omitempty"`
	// Revision of the API proxy.
	Name string `json:"name,omitempty"`
	// Name of the organization.
	Organization string `json:"organization,omitempty"`
}
