/*
 * Stats API
 *
 * Access metrics collected by Apigee Edge that measure API consumption and performance that are used to build Analytics reports.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// MetricsDimensions struct for MetricsDimensions
type MetricsDimensions struct {
	// List of metrics.
	Metrics []MetricsMetrics `json:"metrics,omitempty"`
	// Dimension name.
	Name string `json:"name,omitempty"`
}