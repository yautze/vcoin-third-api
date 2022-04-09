package record

// Record -
type Record struct {
	// Stream
	Stream string `json:"stream"`

	Data *Data `json:"data"`
}

// Data -
type Data struct {
	// 事件類型
	EventType string `json:"e"`

	// 事件時間
	EventTime int64 `json:"E"`

	// 交易對
	S string `json:"s"`

	// 交易編號
	T int64 `json:"t"`

	// 成交價格
	P string `json:"p"`

	// 成交數量
	Q string `json:"q"`

	// 買方訂單編號
	B int64 `json:"b"`

	// 賣方訂單編號
	A int64 `json:"a"`

	// 成交時間
	Timestamp int64 `json:"T"`

	// 買方是否是做市方。如true，則此次成交是一個主動賣出單，否則是一個主動買入單。
	M bool `json:"m"`

	MM bool `json:"M"`
}
