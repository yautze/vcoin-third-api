package record

import (
	"encoding/json"
	"vcoin-third-api/infra/broker"
	"vcoin-third-api/infra/config"
	"vcoin-third-api/infra/ws"
	"vcoin-third-api/middle"
	"vcoin-third-api/module/record"

	b "go-micro.dev/v4/broker"
)

// GetNewestRecord -
func GetNewestRecord(c *middle.C) {
	res, err := record.Get(c.Request().Context(), config.SetKey)
	if err != nil {
		c.E(err)
		return
	}

	data := new(record.Record)
	if err := json.Unmarshal([]byte(res), &data); err != nil {
		c.E(err)
		return
	}

	c.R(data)
}

// ObserveRecord -
func ObserveRecord(c *middle.C) {
	w, err := ws.NewWebSocket(c.ResponseWriter(), c.Request())
	if err != nil {
		c.E(err)
		return
	}

	brokerClient := broker.Client

	brokerClient.Subscribe(config.ObserveRecordTopic, func(e b.Event) error {
		data := new(record.Record)
		json.Unmarshal(e.Message().Body, &data)
		we := &ws.Event{
			Name: config.ObserveRecordTopic,
			Data: data,
		}

		w.Out <- we.Raw()

		return nil
	})

	<-w.Closer

	return
}
