/*
 * Stats API
 *
 * Access metrics collected by Apigee Edge that measure API consumption and performance that are used to build Analytics reports.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// MetricsEnvironments struct for MetricsEnvironments
type MetricsEnvironments struct {
	// Dimension details.
	Dimensions []MetricsDimensions `json:"dimensions,omitempty"`
	// Environment name.
	Name string `json:"name,omitempty"`
}