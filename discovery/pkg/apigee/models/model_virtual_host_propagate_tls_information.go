/*
 * Virtual hosts API
 *
 * Manage virtual hosts. Virtual hosts let multiple domain names connect to the same host. A virtual host on Edge defines the domains and ports on which an API proxy is exposed, and, by extension, the URL that apps use to access an API proxy. A virtual host also defines whether the API proxy is accessed by using the HTTP protocol, or by the encrypted HTTPS protocol.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// VirtualHostPropagateTlsInformation **Edge for Public Cloud only (Alpha).**
type VirtualHostPropagateTlsInformation struct {
	// Flag that specifies whether to enable the capture of TLS connection information by Edge. This information will be made available as flow variables in an API proxy. See <a href=\"https://docs.apigee.com/api-platform/system-administration/tls-vars\">Accessing TLS connection information in an API proxy</a> for more.
	ConnectionProperties string `json:"ConnectionProperties,omitempty"`
	// Flag that specifies whether to enable the capture of client cert details captured by Edge in two-way TLS. This information will be made available as flow variables in an API proxy. See <a href=\"https://docs.apigee.com/api-platform/system-administration/tls-vars\">Accessing TLS connection information in an API proxy</a> for more.
	ClientProperties string `json:"ClientProperties,omitempty"`
}
