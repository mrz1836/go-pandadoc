package pandadoc

import "context"

// DocumentsService handles document-related PandaDoc API calls.
type DocumentsService interface {
	List(ctx context.Context, opts *ListDocumentsOptions) (*DocumentListResponse, error)
	Create(ctx context.Context, reqBody DocumentCreateRequest) (*DocumentCreateResponse, error)
	CreateFromUpload(ctx context.Context, reqBody *CreateDocumentFromUploadRequest) (*DocumentCreateResponse, error)
	Status(ctx context.Context, id string) (*DocumentStatusResponse, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, reqBody DocumentUpdateRequest) error
	ESignDisclosure(ctx context.Context, documentID string) (*DocumentESignDisclosureResponse, error)
	ChangeStatus(ctx context.Context, id string, reqBody *ChangeDocumentStatusRequest) error
	ChangeStatusWithUpload(ctx context.Context, id string, reqBody *ChangeDocumentStatusWithUploadRequest) error
	RevertToDraft(ctx context.Context, id string) (*DocumentRevertToDraftResponse, error)
	Details(ctx context.Context, id string) (*DocumentDetailsResponse, error)
	Send(ctx context.Context, id string, reqBody DocumentSendRequest) (*DocumentSendResponse, error)
	CreateEditingSession(ctx context.Context, id string, reqBody CreateDocumentEditingSessionRequest) (*CreateDocumentEditingSessionResponse, error)
	CreateSession(ctx context.Context, id string, reqBody CreateDocumentSessionRequest) (*CreateDocumentSessionResponse, error)
	Download(ctx context.Context, id string) (*DownloadResponse, error)
	DownloadProtected(ctx context.Context, id string) (*DownloadResponse, error)
	TransferOwnership(ctx context.Context, id string, reqBody TransferDocumentOwnershipRequest) error
	TransferAllOwnership(ctx context.Context, reqBody TransferAllDocumentsOwnershipRequest) error
	MoveToFolder(ctx context.Context, id, folderID string) error
	AppendContentLibraryItem(ctx context.Context, id string, reqBody AppendContentLibraryItemRequest) (*AppendContentLibraryItemResponse, error)
}

// ProductCatalogService handles product-catalog API operations.
type ProductCatalogService interface {
	Search(ctx context.Context, opts *SearchProductCatalogItemsOptions) (*SearchProductCatalogItemsResponse, error)
	Create(ctx context.Context, reqBody CreateProductCatalogItemRequest) (*ProductCatalogItemResponse, error)
	Get(ctx context.Context, itemUUID string) (*ProductCatalogItemResponse, error)
	Update(ctx context.Context, itemUUID string, reqBody UpdateProductCatalogItemRequest) (*ProductCatalogItemResponse, error)
	Delete(ctx context.Context, itemUUID string) error
}

// OAuthService handles OAuth token operations.
type OAuthService interface {
	Token(ctx context.Context, req *OAuthTokenRequest) (*OAuthTokenResponse, error)
}

// WebhookSubscriptionsService handles webhook-subscription endpoints.
type WebhookSubscriptionsService interface {
	List(ctx context.Context, opts *ListWebhookSubscriptionsOptions) (*WebhookSubscriptionListResponse, error)
	Create(ctx context.Context, reqBody *WebhookSubscriptionRequest) (*WebhookSubscription, error)
	Get(ctx context.Context, id string) (*WebhookSubscription, error)
	Update(ctx context.Context, id string, reqBody *WebhookSubscriptionRequest) (*WebhookSubscription, error)
	Delete(ctx context.Context, id string) error
	RegenerateSharedKey(ctx context.Context, id string) (*UpdateWebhookSubscriptionSharedKeyResponse, error)
}

// WebhookEventsService handles webhook-event endpoints.
type WebhookEventsService interface {
	List(ctx context.Context, opts *ListWebhookEventsOptions) (*WebhookEventListResponse, error)
	Get(ctx context.Context, id string) (*WebhookEventDetailsResponse, error)
}
