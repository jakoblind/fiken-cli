package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// Client is the Fiken API HTTP client with auth, rate limiting, and pagination.
type Client struct {
	token      string
	httpClient *http.Client
	baseURL    string

	// Rate limiting: max 4 req/sec, 1 concurrent
	mu       sync.Mutex
	lastReq  time.Time
	minDelay time.Duration
}

// NewClient creates a new Fiken API client.
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:  BaseURL,
		minDelay: 250 * time.Millisecond, // 4 req/sec
	}
}

// PaginationInfo holds pagination metadata from response headers.
type PaginationInfo struct {
	Page        int
	PageSize    int
	PageCount   int
	ResultCount int
}

// APIError represents an error from the Fiken API.
type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("fiken API error %d: %s", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("fiken API error %d: %s", e.StatusCode, e.Status)
}

// doRequest performs a rate-limited HTTP request.
func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	elapsed := time.Since(c.lastReq)
	if elapsed < c.minDelay {
		time.Sleep(c.minDelay - elapsed)
	}
	c.lastReq = time.Now()
	c.mu.Unlock()

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if req.Method == http.MethodPost || req.Method == http.MethodPut {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(body),
		}
	}

	return resp, nil
}

// Get performs a GET request to the given path and decodes the response.
func (c *Client) Get(path string, result interface{}) (*PaginationInfo, error) {
	return c.GetWithParams(path, nil, result)
}

// GetWithParams performs a GET request with query parameters.
func (c *Client) GetWithParams(path string, params url.Values, result interface{}) (*PaginationInfo, error) {
	u := c.baseURL + path
	if params != nil {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return nil, fmt.Errorf("decoding response: %w (body: %s)", err, truncate(string(body), 200))
		}
	}

	pagination := parsePagination(resp)
	return pagination, nil
}

// Post performs a POST request with a JSON body.
func (c *Client) Post(path string, body interface{}, result interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("encoding request: %w", err)
	}

	u := c.baseURL + path
	req, err := http.NewRequest(http.MethodPost, u, io.NopCloser(
		io.NewSectionReader(newBytesReaderAt(bodyBytes), 0, int64(len(bodyBytes))),
	))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading response: %w", err)
		}
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

// GetAllPages fetches all pages for a paginated endpoint.
// The fetchPage function should perform the actual request and return results + whether there are more pages.
func (c *Client) GetAllPages(path string, pageSize int, fetchPage func(page int) (int, error)) error {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	page := 0
	for {
		pageCount, err := fetchPage(page)
		if err != nil {
			return err
		}
		page++
		if page >= pageCount {
			break
		}
	}
	return nil
}

func parsePagination(resp *http.Response) *PaginationInfo {
	info := &PaginationInfo{}
	if v := resp.Header.Get(HeaderPage); v != "" {
		info.Page, _ = strconv.Atoi(v)
	}
	if v := resp.Header.Get(HeaderPageSize); v != "" {
		info.PageSize, _ = strconv.Atoi(v)
	}
	if v := resp.Header.Get(HeaderPageCount); v != "" {
		info.PageCount, _ = strconv.Atoi(v)
	}
	if v := resp.Header.Get(HeaderResultCount); v != "" {
		info.ResultCount, _ = strconv.Atoi(v)
	}
	return info
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// bytesReaderAt adapts a byte slice to io.ReaderAt.
type bytesReaderAt struct {
	data []byte
}

func newBytesReaderAt(data []byte) *bytesReaderAt {
	return &bytesReaderAt{data: data}
}

func (b *bytesReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if off >= int64(len(b.data)) {
		return 0, io.EOF
	}
	n = copy(p, b.data[off:])
	if n < len(p) {
		err = io.EOF
	}
	return
}
