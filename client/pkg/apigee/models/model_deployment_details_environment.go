/*
 * Deployments
 *
 * Manage API proxy and shared flow deployments.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// DeploymentDetailsEnvironment struct for DeploymentDetailsEnvironment
type DeploymentDetailsEnvironment struct {
	// Name of the environment.
	Name string `json:"name,omitempty"`
	// Revision details.
	Revision []DeploymentDetailsRevision `json:"revision,omitempty"`
}
