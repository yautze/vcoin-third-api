package middle

import (
	"github.com/kataras/iris/v12"
)

// C -
type C struct {
	iris.Context
}

// response -
type response struct {
	// Data
	Data interface{} `json:"data"`
	// ErrMsg
	ErrMsg interface{} `json:"errMsg"`
}

// HandleFunc -
func HandleFunc(handler func(*C)) func(iris.Context) {
	return func(c iris.Context) {

		customerContext := &C{
			c,
		}

		handler(customerContext)
	}
}

// R -
func (c *C) R(data interface{}) {
	c.StatusCode(iris.StatusOK)
	c.JSON(response{Data: data, ErrMsg: nil})
}

// E -
func (c *C) E(err error) {
	c.StatusCode(iris.StatusOK)
	c.JSON(response{Data: nil, ErrMsg: err.Error()})
}
