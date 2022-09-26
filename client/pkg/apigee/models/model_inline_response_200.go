/*
 * Resource files API
 *
 * Manage files containing executable code or definitions used by API policies to enable custom behavior and extensibility.   Resource files are executable code or other types of assets (for example XSLT) that are used by API proxies at runtime. Resource files can be stored at one of three levels: * **API proxy**: Available to any policies in an API proxy. * **Environment**: Available to any policies in any API proxy deployed in the environment. * **Organization**: Available to any API proxy deployed in any environment in an organization.  Resource files are resolved by name. Apigee Edge resolves resource files from most specific (API proxy) to the most general (organization). This enables you to store generic code that provides utility processing at the organization level. This provides for greater maintainability, since generic code is not repeated across multiple API proxies. A good example of code that might be scoped to the organization is a library to do Base64 encoding.  For more information, see <a href=\"https://docs.apigee.com/api-platform/develop/resource-files\">Manage resources</a>.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// InlineResponse200 struct for InlineResponse200
type InlineResponse200 struct {
	ResourceFile []ResourceFile `json:"resourceFile,omitempty"`
}
