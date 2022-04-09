package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	record_controller "vcoin-third-api/controller/record"
	"vcoin-third-api/infra/broker"
	"vcoin-third-api/infra/config"
	"vcoin-third-api/infra/redis"
	"vcoin-third-api/middle"
	"vcoin-third-api/module/record"

	"github.com/kataras/iris/v12"
	"github.com/spf13/cobra"
	b "go-micro.dev/v4/broker"
	"golang.org/x/net/websocket"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		server()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&config.C.Port, "port", "p", "3000", "port")
	serverCmd.Flags().StringVarP(&config.C.WssTarget, "wss", "w", "wss://stream.yshyqxx.com/stream?streams=btcusdt@aggTrade", "wss target")
	serverCmd.Flags().StringVarP(&config.C.HttpsTarget, "https", "s", "https://stream.yshyqxx.com", "https target")
}

// server -
func server() {
	var ws *websocket.Conn
	var err error
	ctx := context.Background()

	// init redis
	if err := redis.New(); err != nil {
		fmt.Printf("Dial Redis failed: %s\n", err.Error())
		os.Exit(1)
	}

	// init memory broker
	if err := broker.NewLocalBroker(); err != nil {
		fmt.Printf("Dial LocalBroker failed: %s\n", err.Error())
		os.Exit(1)
	}

	// wss observer
	go func() {
		ws, err = websocket.Dial(config.C.WssTarget, "", config.C.WssTarget)
		if err != nil {
			fmt.Printf("Dial Wss failed: %s\n", err.Error())
			os.Exit(1)
		}

		msgCh := make(chan []byte)
		go readStream(ws, msgCh)
		for {
			select {
			case message := <-msgCh:
				// record to redis
				record.Create(ctx, config.SetKey, message, 0)

				// publish to broker
				broker.Client.Publish(config.ObserveRecordTopic, &b.Message{
					Body: message,
				})
			}
		}
	}()

	// iris
	irisServer()

}

// irisServer
func irisServer() {
	app := iris.New()
	app.PartyFunc("/api", func(p iris.Party) {
		p.Get("/record/newest", middle.HandleFunc(record_controller.GetNewestRecord))
		p.Get("/record/observe", middle.HandleFunc(record_controller.ObserveRecord))
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%s", config.C.Port)))
}

// readStream -
func readStream(ws *websocket.Conn, msgCh chan []byte) {
	for {
		var msg []byte
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			fmt.Printf("Error::: %s\n", err.Error())
			return
		}
		msgCh <- msg
	}
}
