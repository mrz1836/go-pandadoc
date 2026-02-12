// Package models provides data models for the PandaDoc API.
package models

import "time"

// Document represents a PandaDoc document.
type Document struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Status         string                 `json:"status"`
	DateCreated    time.Time              `json:"date_created"`
	DateModified   time.Time              `json:"date_modified"`
	DateCompleted  time.Time              `json:"date_completed,omitempty"`
	ExpirationDate *time.Time             `json:"expiration_date,omitempty"`
	UUID           string                 `json:"uuid,omitempty"`
	Version        string                 `json:"version,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// DocumentStatus represents the status of a document.
type DocumentStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// DocumentField represents an embedded field in a document.
type DocumentField struct {
	Name     string      `json:"name"`
	Value    interface{} `json:"value"`
	Type     string      `json:"type,omitempty"`
	Required bool        `json:"required,omitempty"`
}

// DocumentDetails represents detailed document information including fields.
type DocumentDetails struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Status        string                 `json:"status"`
	DateCreated   time.Time              `json:"date_created"`
	DateModified  time.Time              `json:"date_modified"`
	Fields        map[string]interface{} `json:"fields,omitempty"`
	Tokens        []DocumentToken        `json:"tokens,omitempty"`
	PricingTables []PricingTable         `json:"pricing_tables,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// DocumentToken represents a token/variable in a document.
type DocumentToken struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// PricingTable represents a pricing table in a document.
type PricingTable struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Sections []TableSection `json:"sections,omitempty"`
}

// TableSection represents a section in a pricing table.
type TableSection struct {
	Title string     `json:"title,omitempty"`
	Rows  []TableRow `json:"rows,omitempty"`
}

// TableRow represents a row in a pricing table section.
type TableRow struct {
	ID           string                 `json:"id"`
	Data         map[string]interface{} `json:"data,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

// DocumentListResponse represents a paginated list of documents.
type DocumentListResponse struct {
	Results  []Document `json:"results"`
	Count    int        `json:"count"`
	Next     *string    `json:"next,omitempty"`
	Previous *string    `json:"previous,omitempty"`
}
