/*
 * OFAC API
 *
 * OFAC (Office of Foreign Assets Control) API is designed to facilitate the enforcement of US government economic sanctions programs required by federal law. This project implements a modern REST HTTP API for companies and organizations to obey federal law and use OFAC data in their applications.
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"time"
)

// Download Metadata and stats about downloaded OFAC data
type Download struct {
	SDNs      int32     `json:"SDNs,omitempty"`
	AltNames  int32     `json:"altNames,omitempty"`
	Addresses int32     `json:"addresses,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
