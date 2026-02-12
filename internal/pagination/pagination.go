package pagination

import "github.com/mrz1836/go-pandadoc/models"

// Package pagination provides utilities for handling API pagination.

// NextPage returns the next page options based on the current page.
func NextPage(current *models.ListOptions) *models.ListOptions {
	if current == nil {
		return &models.ListOptions{
			Page:  1,
			Count: 50, // Default count
		}
	}

	return &models.ListOptions{
		Page:  current.Page + 1,
		Count: current.Count,
	}
}

// HasMore determines if there are more pages based on pagination metadata.
func HasMore(meta *models.PaginationMeta) bool {
	return meta != nil && meta.Next != nil
}
