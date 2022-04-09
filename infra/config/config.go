package config

//
const (
	SetKey             = "streams=btcusdt@aggTrade"
	ObserveRecordTopic = "observe.record.topic"
)

// C -
var C Config

// Config -
type Config struct {
	// server prot
	Port string

	// wsss target
	WssTarget string

	// https target
	HttpsTarget string
}
