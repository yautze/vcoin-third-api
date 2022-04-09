package ws

import "encoding/json"

// EventHandler -
type EventHandler func(*Event)

// Event -
type Event struct {
	// Name -
	Name string `json:"event"`

	// Data -
	Data interface{} `json:"data"`
}

// NewEvent -
func NewEvent(rawData []byte) (*Event, error) {
	e := new(Event)

	err := json.Unmarshal(rawData, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// Raw -
func (e *Event) Raw() []byte {
	raw, _ := json.Marshal(e)
	return raw
}

// DataRaw -
func (e *Event) DataRaw() []byte {
	raw, _ := json.Marshal(e.Data)
	return raw
}
