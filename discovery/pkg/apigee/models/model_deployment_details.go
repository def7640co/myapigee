/*
 * Deployments
 *
 * Manage API proxy and shared flow deployments.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// DeploymentDetails API proxy or shared flow deployment details.
type DeploymentDetails struct {
	// API proxy or shared flow deployment details in the environment.
	Environment []DeploymentDetailsEnvironment `json:"environment,omitempty"`
	// Name of the API proxy or shared flow.
	Name string `json:"name,omitempty"`
	// Name of the organization.
	Organization string `json:"organization,omitempty"`
}
