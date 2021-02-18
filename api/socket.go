package api

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

func (api *API) SocketHandler(ctx *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)

	for {
		if mt, msg, err = ctx.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)

		if err = ctx.WriteMessage(mt, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
