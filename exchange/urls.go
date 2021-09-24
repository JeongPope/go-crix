package exchange

const (
	marketURL   = "https://api.upbit.com/v1/market/all?isDetails=true" // Markets
	tickerURL   = "wss://api.upbit.com/websocket/v1"                   // Tickers, WEBSOCKET API
	snapshotURL = "https://api.upbit.com/v1/ticker?markets="           // Snapshot
)
