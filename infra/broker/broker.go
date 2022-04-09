package broker

import "go-micro.dev/v4/broker"

// Client
var Client broker.Broker

// NewLocalBroker -
func NewLocalBroker() error {
	b := broker.NewMemoryBroker()

	// connect (use local memory, so never have error)
	if err := b.Connect(); err != nil {
		return err
	}

	Client = b

	return nil
}