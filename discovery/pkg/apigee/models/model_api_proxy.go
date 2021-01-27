/*
 * API Proxies API
 *
 * Manage API proxies. You expose APIs on Apigee Edge by implementing API proxies.  API proxies decouple the app-facing API from your backend services, shielding those apps from backend code changes. As you make backend changes to your services, apps continue to call the same API without any interruption. For more information, see <a href=\"https://docs.apigee.com/api-platform/fundamentals/understanding-apis-and-api-proxies\">Understanding APIs and API proxies</a>.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

// ApiProxy API proxy metadata and revisions.
type ApiProxy struct {
	MetaData ApiProxyMetaData `json:"metaData,omitempty"`
	// Name of the API proxy.
	Name string `json:"name,omitempty"`
	// Revisions defined for the API proxy.
	Revision []string `json:"revision,omitempty"`
}
