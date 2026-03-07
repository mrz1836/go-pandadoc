package pandadoc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var errOnlyOneBodyType = fmt.Errorf("only one request body type can be set")

type request struct {
	method      string
	path        string
	query       url.Values
	headers     http.Header
	requireAuth bool
	accept      string

	jsonBody  any
	formBody  url.Values
	multipart *multipartPayload

	expectedStatus []int
}

type multipartPayload struct {
	Fields map[string]string
	Files  []multipartFile
}

type multipartFile struct {
	FieldName   string
	FileName    string
	ContentType string
	Reader      io.Reader
}

// DownloadResponse is a streaming response for binary download endpoints.
type DownloadResponse struct {
	Body               io.ReadCloser
	Headers            http.Header
	StatusCode         int
	ContentType        string
	ContentDisposition string
	ContentLength      int64
}

// Close closes the response body.
func (d *DownloadResponse) Close() error {
	if d == nil || d.Body == nil {
		return nil
	}
	return d.Body.Close()
}

func (c *Client) do(ctx context.Context, req *request) (*http.Response, error) {
	if req == nil {
		return nil, ErrNilRequest
	}

	fullURL, err := c.buildURL(req.path, req.query)
	if err != nil {
		return nil, err
	}

	bodyBytes, contentType, err := encodeRequestBody(req)
	if err != nil {
		return nil, err
	}

	attempt := 0
	for {
		ok, resp, err := c.doAttemptWithHandling(ctx, req, fullURL, bodyBytes, contentType, attempt)
		if err != nil {
			return nil, err
		}
		if ok {
			return resp, nil
		}
		attempt++
	}
}

func (c *Client) doAttemptWithHandling(ctx context.Context, req *request, fullURL string, bodyBytes []byte, contentType string, attempt int) (bool, *http.Response, error) {
	c.logDebug("API Request: %s %s (attempt %d)", req.method, fullURL, attempt+1)

	resp, retryable, err := c.doAttempt(ctx, req, fullURL, bodyBytes, contentType)
	if err != nil {
		c.logError("Request failed: %v", err)
		if !retryable || !c.shouldRetryOnError(attempt, err) {
			return false, nil, err
		}
		c.logInfo("Retrying after error: %v", err)
		if sleepErr := sleepWithContext(ctx, c.backoff(attempt)); sleepErr != nil {
			return false, nil, sleepErr
		}
		return false, nil, nil
	}

	c.logDebug("API Response: %d %s", resp.StatusCode, resp.Status)

	if c.shouldRetryOnStatus(attempt, resp) {
		retryDelay := c.backoff(attempt)
		if retryAfter, ok := parseRetryAfter(resp.Header.Get("Retry-After")); ok {
			retryDelay = retryAfter
		}
		c.logInfo("Retrying on status %d, waiting %v", resp.StatusCode, retryDelay)
		_ = drainAndClose(resp.Body)

		if sleepErr := sleepWithContext(ctx, retryDelay); sleepErr != nil {
			return false, nil, sleepErr
		}
		return false, nil, nil
	}

	if !statusExpected(resp.StatusCode, req.expectedStatus) {
		apiErr := parseAPIError(resp)
		c.logError("API Error: %v", apiErr)
		_ = resp.Body.Close()
		return false, nil, apiErr
	}

	return true, resp, nil
}

func (c *Client) doAttempt(ctx context.Context, req *request, fullURL string, bodyBytes []byte, contentType string) (*http.Response, bool, error) {
	httpReq, buildErr := http.NewRequestWithContext(ctx, req.method, fullURL, bytes.NewReader(bodyBytes))
	if buildErr != nil {
		return nil, false, fmt.Errorf("build request: %w", buildErr)
	}

	if len(bodyBytes) > 0 {
		httpReq.Header.Set("Content-Type", contentType)
	}
	accept := req.accept
	if accept == "" {
		accept = "application/json"
	}
	httpReq.Header.Set("Accept", accept)
	httpReq.Header.Set("User-Agent", c.userAgent)

	for k, vals := range req.headers {
		for _, v := range vals {
			httpReq.Header.Add(k, v)
		}
	}

	if err := c.injectAuth(httpReq, req.requireAuth); err != nil {
		return nil, false, err
	}

	resp, doErr := c.httpClient.Do(httpReq)
	return resp, true, doErr
}

func (c *Client) decodeJSON(ctx context.Context, req *request, out any) error {
	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if out == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		if err == io.EOF {
			return nil
		}
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}

func (c *Client) download(ctx context.Context, req *request) (*DownloadResponse, error) {
	resp, err := c.do(ctx, req)
	if err != nil {
		if resp != nil {
			_ = resp.Body.Close()
		}
		return nil, err
	}

	return &DownloadResponse{
		Body:               resp.Body,
		Headers:            resp.Header.Clone(),
		StatusCode:         resp.StatusCode,
		ContentType:        resp.Header.Get("Content-Type"),
		ContentDisposition: resp.Header.Get("Content-Disposition"),
		ContentLength:      resp.ContentLength,
	}, nil
}

func (c *Client) injectAuth(req *http.Request, required bool) error {
	if c.apiKey != "" {
		req.Header.Set("Authorization", "API-Key "+c.apiKey)
		return nil
	}
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
		return nil
	}
	if required {
		return ErrMissingAuthentication
	}
	return nil
}

func (c *Client) shouldRetryOnError(attempt int, _ error) bool {
	return attempt < c.retryPolicy.MaxRetries
}

func (c *Client) shouldRetryOnStatus(attempt int, resp *http.Response) bool {
	if attempt >= c.retryPolicy.MaxRetries || resp == nil {
		return false
	}

	if resp.StatusCode == http.StatusTooManyRequests && c.retryPolicy.RetryOn429 {
		return true
	}

	if resp.StatusCode >= 500 && resp.StatusCode <= 599 && c.retryPolicy.RetryOn5xx {
		return true
	}

	return false
}

func (c *Client) backoff(attempt int) time.Duration {
	if attempt <= 0 {
		return c.retryPolicy.InitialBackoff
	}
	backoff := c.retryPolicy.InitialBackoff
	for i := 0; i < attempt; i++ {
		backoff *= 2
		if backoff >= c.retryPolicy.MaxBackoff {
			return c.retryPolicy.MaxBackoff
		}
	}
	return backoff
}

var errEndpointPathRequired = fmt.Errorf("endpoint path is required")

func (c *Client) buildURL(endpointPath string, query url.Values) (string, error) {
	if strings.TrimSpace(endpointPath) == "" {
		return "", errEndpointPathRequired
	}

	rel, err := url.Parse(endpointPath)
	if err != nil {
		return "", fmt.Errorf("parse endpoint path: %w", err)
	}

	u := *c.baseURL
	u.Path = joinPaths(c.baseURL.Path, rel.Path)

	finalQuery := u.Query()
	for k, vals := range rel.Query() {
		for _, v := range vals {
			finalQuery.Add(k, v)
		}
	}
	for k, vals := range query {
		for _, v := range vals {
			finalQuery.Add(k, v)
		}
	}
	u.RawQuery = finalQuery.Encode()

	return u.String(), nil
}

func joinPaths(basePath, relPath string) string {
	basePath = strings.TrimSpace(basePath)
	relPath = strings.TrimSpace(relPath)

	if basePath == "" {
		basePath = "/"
	}
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}

	relPath = strings.TrimPrefix(relPath, "/")
	if relPath == "" {
		if strings.HasSuffix(basePath, "/") {
			trimmed := strings.TrimSuffix(basePath, "/")
			if trimmed == "" {
				return "/"
			}
			return trimmed
		}
		return basePath
	}

	trimmedBase := strings.TrimSuffix(basePath, "/")
	if trimmedBase == "" {
		trimmedBase = "/"
	}

	if trimmedBase == "/" {
		return "/" + relPath
	}
	return trimmedBase + "/" + relPath
}

func encodeRequestBody(req *request) ([]byte, string, error) {
	bodyKinds := 0
	if req.jsonBody != nil {
		bodyKinds++
	}
	if req.formBody != nil {
		bodyKinds++
	}
	if req.multipart != nil {
		bodyKinds++
	}
	if bodyKinds > 1 {
		return nil, "", errOnlyOneBodyType
	}

	if req.jsonBody != nil {
		payload, err := json.Marshal(req.jsonBody)
		if err != nil {
			return nil, "", fmt.Errorf("encode JSON request body: %w", err)
		}
		return payload, "application/json", nil
	}
	if req.formBody != nil {
		return []byte(req.formBody.Encode()), "application/x-www-form-urlencoded", nil
	}
	if req.multipart != nil {
		payload, contentType, err := encodeMultipart(req.multipart)
		if err != nil {
			return nil, "", err
		}
		return payload, contentType, nil
	}

	return nil, "", nil
}

func encodeMultipart(payload *multipartPayload) ([]byte, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for k, v := range payload.Fields {
		if err := writer.WriteField(k, v); err != nil {
			return nil, "", fmt.Errorf("write multipart field %q: %w", k, err)
		}
	}

	for _, file := range payload.Files {
		if file.Reader == nil {
			return nil, "", ErrNilFileReader
		}
		field := file.FieldName
		if field == "" {
			field = "file"
		}
		fileName := file.FileName
		if fileName == "" {
			fileName = "upload.bin"
		}

		part, err := writer.CreateFormFile(field, fileName)
		if err != nil {
			return nil, "", fmt.Errorf("create multipart file %q: %w", field, err)
		}
		if _, err := io.Copy(part, file.Reader); err != nil {
			return nil, "", fmt.Errorf("copy multipart file %q: %w", field, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("close multipart body: %w", err)
	}

	return body.Bytes(), writer.FormDataContentType(), nil
}

func statusExpected(status int, expected []int) bool {
	if len(expected) == 0 {
		return status >= 200 && status < 300
	}
	for _, want := range expected {
		if status == want {
			return true
		}
	}
	return false
}

func parseAPIError(resp *http.Response) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		RequestID:  headerFirst(resp.Header, "X-Request-Id", "X-Request-ID", "Request-Id"),
		RetryAfter: resp.Header.Get("Retry-After"),
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr == nil {
		apiErr.RawBody = string(body)
		populateAPIErrorFromBody(apiErr, body)
	}

	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(resp.StatusCode)
	}

	return apiErr
}

func populateAPIErrorFromBody(apiErr *APIError, body []byte) {
	if len(body) == 0 {
		return
	}

	var obj map[string]any
	if err := json.Unmarshal(body, &obj); err != nil {
		if apiErr.Message == "" {
			apiErr.Message = strings.TrimSpace(string(body))
		}
		return
	}

	extractErrorCode(apiErr, obj)
	extractErrorMessage(apiErr, obj)
	extractErrorDetails(apiErr, obj)
}

func extractErrorCode(apiErr *APIError, obj map[string]any) {
	for _, key := range []string{"code", "type", "error"} {
		if code, ok := stringValue(obj[key]); ok && code != "" {
			apiErr.Code = code
			return
		}
	}
}

func extractErrorMessage(apiErr *APIError, obj map[string]any) {
	for _, key := range []string{"message", "detail", "error_description", "error"} {
		if message, ok := stringValue(obj[key]); ok && message != "" {
			apiErr.Message = message
			return
		}
	}
}

func extractErrorDetails(apiErr *APIError, obj map[string]any) {
	if details, ok := obj["details"]; ok {
		apiErr.Details = details
		return
	}
	if details, ok := obj["detail"]; ok {
		switch details.(type) {
		case map[string]any, []any:
			apiErr.Details = details
		}
	}
}

func parseRetryAfter(value string) (time.Duration, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	if seconds, err := strconv.Atoi(value); err == nil && seconds >= 0 {
		return time.Duration(seconds) * time.Second, true
	}
	if when, err := http.ParseTime(value); err == nil {
		delta := time.Until(when)
		if delta < 0 {
			return 0, true
		}
		return delta, true
	}
	return 0, false
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func drainAndClose(body io.ReadCloser) error {
	if body == nil {
		return nil
	}
	_, _ = io.Copy(io.Discard, body)
	return body.Close()
}

func headerFirst(h http.Header, keys ...string) string {
	for _, k := range keys {
		if v := strings.TrimSpace(h.Get(k)); v != "" {
			return v
		}
	}
	return ""
}

func stringValue(v any) (string, bool) {
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	return s, true
}

func escapePathParam(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrEmptyPathParameter
	}
	return url.PathEscape(trimmed), nil
}
