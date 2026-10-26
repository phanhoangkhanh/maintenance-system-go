package tradingSocket

import "encoding/json"


type TradingKeyInsideKeyMessage struct {
	DomainName   string `json:"domainName"`
	Prefix 	string `json:"prefix"`
	MessageName string `json:"messageName"`
	Suffix string `json:"suffix"`
	MessageType string `json:"messageType"`
}

type TradingKeyMessage struct {
	Key   TradingKeyInsideKeyMessage `json:"key"`
	Payload string `json:"payload"`
}

type MarketPriceInfo struct {
    Time   int64  `json:"t"`
    Symbol string `json:"s"`

    Open   string `json:"o"`
    High   string `json:"h"`
    Low    string `json:"l"`
    Last   string `json:"lt"`

    Volume string `json:"vol"`
    CVol   string `json:"cVol"`
    Price  string `json:"p"`

    Ask    string `json:"a"`
    Bid    string `json:"b"`
}

//Normal msg
type MessageTrading struct {
	Key json.RawMessage  `json:"key"`
	SessionId json.RawMessage `json:"sessionId"`
	Timestamp json.RawMessage `json:"timestamp"`
	Payload []string `json:"payload"`
}