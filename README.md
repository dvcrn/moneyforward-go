
# MoneyForward Go Client

A Go client library for interacting with the MoneyForward API.

## Installation


go get github.com/dvcrn/moneyforward-go


## Usage


Check out `examples/main.go` for a full example.


```
// Initialize client
cookieString := "YOUR_COOKIE_STRING"
client := moneyforward.NewClient(cookieString)

// Get account summaries
accounts, err := client.GetAccountSummaries()
if err != nil {
    log.Fatal(err)
}

// Get recent transactions
transactions, err := client.GetUserAssetActs(moneyforward.UserAssetActsParams{
    Size:         16,
    IsNew:        true,
    IsContinuous: true,
})
```

## Available Methods

- `GetAccountSummaries()` - Get summary of all accounts
- `GetUserAssetActs(params)` - Get transaction history with pagination
- `GetUserAssetAct(id)` - Get details of a specific transaction
- `GetHomeTimeline(limit)` - Get home timeline data
- `ForceUpdate()` - Force update of account data
- `GetTransactions()` - Get all transactions
- `GetAccount(path)` - Get details of a specific account

## Configuration

The client can be configured with:

- `SetCookie(cookie)` - Set authentication cookie
- `SetBaseURL(url)` - Override default API URL
- `WithHeader(key, value)` - Add custom headers to requests

## Example

See [cmd/run/main.go](cmd/run/main.go) for a complete example implementation.

## License

MIT
