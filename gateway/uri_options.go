package gateway

type UriOptions struct {
	Endpoint            string
	TickerUri           string
	DepthUri            string
	KlineUri            string
	GetOrderUri         string
	GetPendingOrdersUri string
	GetHistoryOrdersUri string
	CancelOrderUri      string
	NewOrderUri         string
}
