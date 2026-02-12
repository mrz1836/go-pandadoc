package pandadoc

// ProductCatalogItemType is the catalog item type enum.
type ProductCatalogItemType string

// Product catalog item type constants.
const (
	// ProductCatalogItemTypeRegular represents a regular catalog item.
	ProductCatalogItemTypeRegular ProductCatalogItemType = "regular"
	// ProductCatalogItemTypeBundle represents a bundle catalog item.
	ProductCatalogItemTypeBundle ProductCatalogItemType = "bundle"
)

// ProductCatalogBillingType is the catalog billing type enum.
type ProductCatalogBillingType string

// Product catalog billing type constants.
const (
	// ProductCatalogBillingTypeOneTime represents a one-time billing type.
	ProductCatalogBillingTypeOneTime ProductCatalogBillingType = "one_time"
	// ProductCatalogBillingTypeRecurring represents a recurring billing type.
	ProductCatalogBillingTypeRecurring ProductCatalogBillingType = "recurring"
)

// SearchProductCatalogItemsOptions controls catalog search filters.
type SearchProductCatalogItemsOptions struct {
	Page         int
	PerPage      int
	Query        string
	OrderBy      string
	Types        []ProductCatalogItemType
	BillingTypes []ProductCatalogBillingType
	ExcludeUUIDs []string
	CategoryID   string
	NoCategory   *bool
}

// SearchProductCatalogItemsResponse is returned by catalog search.
type SearchProductCatalogItemsResponse struct {
	Items        []ProductCatalogSearchItem `json:"items"`
	HasMoreItems bool                       `json:"has_more_items"`
	Total        int                        `json:"total"`
}

// ProductCatalogSearchItem is a search result entry.
type ProductCatalogSearchItem struct {
	UUID             string   `json:"uuid,omitempty"`
	WorkspaceID      string   `json:"workspace_id,omitempty"`
	Title            string   `json:"title,omitempty"`
	SKU              string   `json:"sku,omitempty"`
	Description      string   `json:"description,omitempty"`
	Type             string   `json:"type,omitempty"`
	BillingType      string   `json:"billing_type,omitempty"`
	BillingCycle     int      `json:"billing_cycle,omitempty"`
	Currency         string   `json:"currency,omitempty"`
	CategoryID       string   `json:"category_id,omitempty"`
	CategoryName     string   `json:"category_name,omitempty"`
	CreatedBy        string   `json:"created_by,omitempty"`
	ModifiedBy       string   `json:"modified_by,omitempty"`
	DateCreated      string   `json:"date_created,omitempty"`
	DateModified     string   `json:"date_modified,omitempty"`
	PricingMethod    int      `json:"pricing_method,omitempty"`
	BundleItemsCount int      `json:"bundle_items_count,omitempty"`
	ImageSrc         string   `json:"image_src,omitempty"`
	Price            *float64 `json:"price,omitempty"`
	Cost             *float64 `json:"cost,omitempty"`
	MinTierValue     *float64 `json:"min_tier_value,omitempty"`
	MaxTierValue     *float64 `json:"max_tier_value,omitempty"`
	CustomFields     RawJSON  `json:"custom_fields,omitempty"`
	Images           RawJSON  `json:"images,omitempty"`
	Highlights       RawJSON  `json:"highlights,omitempty"`
	Tiers            RawJSON  `json:"tiers,omitempty"`
}

// CreateProductCatalogItemRequest is a flexible create payload.
type CreateProductCatalogItemRequest map[string]any

// UpdateProductCatalogItemRequest is a flexible update payload.
type UpdateProductCatalogItemRequest map[string]any

// ProductCatalogItemResponse is returned by create/get/update endpoints.
type ProductCatalogItemResponse struct {
	UUID                      string  `json:"uuid,omitempty"`
	Title                     string  `json:"title,omitempty"`
	Type                      string  `json:"type,omitempty"`
	CategoryID                string  `json:"category_id,omitempty"`
	CategoryName              string  `json:"category_name,omitempty"`
	CreatedBy                 string  `json:"created_by,omitempty"`
	ModifiedBy                string  `json:"modified_by,omitempty"`
	DateCreated               string  `json:"date_created,omitempty"`
	DateModified              string  `json:"date_modified,omitempty"`
	DefaultPriceConfiguration RawJSON `json:"default_price_configuration,omitempty"`
	Variants                  RawJSON `json:"variants,omitempty"`
	BundleItems               RawJSON `json:"bundle_items,omitempty"`
}
