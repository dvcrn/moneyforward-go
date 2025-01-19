package moneyforward

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://moneyforward.com"
)

// Client represents a MoneyForward API client
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	cookie     string // Changed from authToken to cookie

	// Optional headers
	userAgent      string
	acceptLanguage string
}

// NewClient creates a new MoneyForward API client
func NewClient(cookieString string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		cookie:     cookieString,
	}
}

// SetCookie sets the authentication cookie for requests
func (c *Client) SetCookie(cookie string) {
	c.cookie = cookie
}

// RequestOption allows customizing requests
type RequestOption func(*http.Request)

// WithHeader adds a custom header to the request
func WithHeader(key, value string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}

func (c *Client) newRequest(method, spath string, opts ...RequestOption) (*http.Request, error) {
	u := *c.baseURL
	// Use double slash for sp2 endpoints
	if spath[0] == '/' {
		spath = "/" + spath
	}
	u.Path = spath

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Set cookie instead of Bearer token
	if c.cookie != "" {
		req.Header.Set("Cookie", c.cookie)
	}

	req.Header.Set("User-Agent", "iPhone(iOS:18.2), MoneyFwd-SP(18.1.0) Build:10614")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept", "*/*")

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	// make a new reader with body to allow reading it again
	bodyCloser := io.NopCloser(bytes.NewBuffer(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if v != nil {
		if err := json.NewDecoder(io.NopCloser(bodyCloser)).Decode(v); err != nil {
			// output the raw resp.Body on error
			fmt.Printf("Failed to decode response: %s\n", string(body))

			return err
		}
	}

	return nil
}

// GetHomeTimeline gets the home timeline data
func (c *Client) GetHomeTimeline(limit int) (*HomeTimelineResponse, error) {
	req, err := c.newRequest("GET", "/sp2/home_timeline")
	if err != nil {
		return nil, err
	}

	// Add limit parameter
	params := map[string]string{
		"limit": fmt.Sprintf("%d", limit),
	}
	c.addQueryParams(req, params)

	var resp HomeTimelineResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ForceUpdate forces an update of the data
func (c *Client) ForceUpdate() error {
	req, err := c.newRequest("GET", "/sp2/force_update")
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// GetAccountSummaries gets account summaries
func (c *Client) GetAccountSummaries() (*AccountSummariesResponse, error) {
	req, err := c.newRequest("GET", "/sp2/account_summaries")
	if err != nil {
		return nil, err
	}

	var resp AccountSummariesResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetTransactions gets transaction data
func (c *Client) GetTransactions() (*TransactionsResponse, error) {
	req, err := c.newRequest("GET", "/sp2/transactions")
	if err != nil {
		return nil, err
	}

	var resp TransactionsResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetUserAssetActivities gets user asset activities with pagination and filters
func (c *Client) GetUserAssetActivities(params UserAssetActsParams) (*UserAssetActsResponse, error) {
	req, err := c.newRequest("GET", "/sp2/user_asset_acts")
	if err != nil {
		return nil, err
	}

	// Convert boolean params to integers
	queryParams := map[string]string{
		"is_old":        "0",
		"is_new":        "1",
		"is_continuous": "1",
		"offset":        "0",
		"size":          "16",
	}

	if params.Size > 0 {
		queryParams["size"] = fmt.Sprintf("%d", params.Size)
	}
	if params.Offset > 0 {
		queryParams["offset"] = fmt.Sprintf("%d", params.Offset)
	}
	if params.IsOld {
		queryParams["is_old"] = "1"
	}
	if !params.IsNew {
		queryParams["is_new"] = "0"
	}
	if !params.IsContinuous {
		queryParams["is_continuous"] = "0"
	}

	c.addQueryParams(req, queryParams)

	var resp UserAssetActsResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// UserAssetActsParams represents the parameters for GetUserAssetActivities
type UserAssetActsParams struct {
	IsOld        bool
	IsNew        bool
	IsContinuous bool
	Offset       int
	Size         int
}

// GetUserAssetActivity gets a specific user asset activity by ID
func (c *Client) GetUserAssetActivity(activityID string) (*UserAssetActResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/sp2/user_asset_acts/%s", activityID))
	if err != nil {
		return nil, err
	}

	var resp UserAssetActResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetAccount gets details for a specific account
func (c *Client) GetAccount(mfPath MFShowPath) (*AccountResponse, error) {
	req, err := c.newRequest("GET", string(mfPath))
	if err != nil {
		return nil, err
	}

	var resp AccountResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// SetBaseURL allows overriding the default API URL
func (c *Client) SetBaseURL(urlStr string) error {
	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	c.baseURL = baseURL
	return nil
}

// AddQueryParams adds query parameters to the request URL
func (c *Client) addQueryParams(req *http.Request, params map[string]string) {
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
}

// GetAccountCashFlowTermData gets cash flow data for a specific sub-account within a date range
func (c *Client) GetAccountCashFlowTermData(accountIDHash string, from, to string) (*CashFlowTermDataResponse, error) {
	req, err := c.newRequest("GET", "/sp/cf_term_data_by_account")
	if err != nil {
		return nil, err
	}

	// Add query parameters
	params := map[string]string{
		"account_id_hash": accountIDHash,
		"from":            from,
		"to":              to,
	}
	c.addQueryParams(req, params)

	var resp CashFlowTermDataResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetSubAccountCashFlowTermData gets cash flow data for a specific sub-account within a date range
func (c *Client) GetSubAccountCashFlowTermData(subAccountIDHash string, from, to string) (*CashFlowTermDataResponse, error) {
	req, err := c.newRequest("GET", "/sp/cf_term_data_by_sub_account")
	if err != nil {
		return nil, err
	}

	// Add query parameters
	params := map[string]string{
		"sub_account_id_hash": subAccountIDHash,
		"from":                from,
		"to":                  to,
	}
	c.addQueryParams(req, params)

	var resp CashFlowTermDataResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetSubAccountDetail gets detailed information for a specific sub-account
func (c *Client) GetAccountDetail(accountIDHash string) (*AccountDetailResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/sp/service_detail/%s", accountIDHash))
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"range": "0",
	}
	c.addQueryParams(req, params)

	var resp AccountDetailResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetSubAccountDetail gets detailed information for a specific sub-account
func (c *Client) GetSubAccountDetail(accountIDHash, subAccountIDHash string) (*AccountDetailResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/sp/service_detail/%s", accountIDHash))
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"range":               "0",
		"sub_account_id_hash": subAccountIDHash,
	}
	c.addQueryParams(req, params)

	var resp AccountDetailResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// TriggerAccountAggregation triggers a data aggregation for a specific account
func (c *Client) TriggerAccountAggregation(accountIDHash string) error {
	path := fmt.Sprintf("/sp2/accounts/%s/aggregation_queue", accountIDHash)
	req, err := c.newRequest("POST", path)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
