package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	portEnv := os.Getenv("PORT")
	// if len(portEnv) != 0 {
	// 	portEnv = ":8081"
	// }
	r := GetHandler("/ws")
	r.Run(portEnv)
}

func GetHandler(url string) *gin.Engine {
	r := gin.Default()
	r.GET(url, func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "{msg:'its running'}")
	})
	return r
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	wsupgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsupgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		fmt.Print(err)
	}

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}
