/*
 * OFAC API
 *
 * OFAC (Office of Foreign Assets Control) API is designed to facilitate the enforcement of US government economic sanctions programs required by federal law. This project implements a modern REST HTTP API for companies and organizations to obey federal law and use OFAC data in their applications.
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// Sdn Specially designated national from OFAC list
type Sdn struct {
	EntityID string `json:"entityID,omitempty"`
	SdnName  string `json:"sdnName,omitempty"`
	// SDN's typically represent an individual (customer) or trust/company/organization. OFAC endpoints refer to customers or companies as different entities, but underlying both is the same SDN metadata.
	SdnType string `json:"sdnType,omitempty"`
	Program string `json:"program,omitempty"`
	Title   string `json:"title,omitempty"`
	Remarks string `json:"remarks,omitempty"`
	// Remarks on SDN and often additional information about the SDN
	Match float32 `json:"match,omitempty"`
}
