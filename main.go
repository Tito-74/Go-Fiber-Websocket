package main

import (
	"fmt"
	"io"
	"log"

	// "net/http"

	"github.com/Tito-74/fiber-websocket/database"
	"github.com/Tito-74/fiber-websocket/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	// "github.com/valyala/fasthttp"
)


type Server struct{
	conns map[*websocket.Conn]bool
}

// var upgrader = websocket.FastHTTPUpgrader{
// 	ReadBufferSize:  4096,
// 	WriteBufferSize: 4096,
// 	CheckOrigin: func(fctx *fasthttp.RequestCtx) bool {
// 		if cfg.Origins[0] == "*" {
// 				return true
// 		}
// 	},
// }

func NewServer() *Server{
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWs(c *websocket.Conn) {
	log.Printf("new connection from %v", c.RemoteAddr())

	s.conns[c] = true

	s.readLoop(c)

	


		



}

func(s *Server)readLoop(c *websocket.Conn){
	var msg models.Message
	for {
		 err := c.ReadJSON(&msg)
		if  err != nil {
			
			if err == io.EOF {
				break
			}
			log.Println("read:", err)
			continue
		}
		log.Printf("recv: %s", msg)
		username := c.Locals("username")

		fmt.Println("username1:", username)
		// if err := c.WriteMessage(websocket.TextMessage, []byte("thank you !!")); err != nil {
		// 	fmt.Println("write error:", err)
		// }
		// database.Database.Db.Create(&msg)
		s.Broadcast(&msg)

	}

}

func main() {
	database.ConnectDb()

	server := NewServer()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
    AllowOrigins: "*",
    AllowHeaders:  "Origin, Content-Type, Accept",
}))
	app.Use("/ws/:username", func(c *fiber.Ctx) error {

		username := c.Params("username")
		fmt.Println("username", username)
		c.Locals("username", username)
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:username", websocket.New(server.handleWs))

	

	log.Fatal(app.Listen(":3000"))
	// Access the websocket server: ws://localhost:3000/ws/123?v=1.0
	// https://www.websocket.org/echo.html
}

func(s *Server) Broadcast(b *models.Message){
	// var s *Server
	for ws := range s.conns{
		go func(ws *websocket.Conn) {
			if err := ws.WriteMessage(websocket.TextMessage, []byte(b.Content)); err != nil {
				fmt.Println("write err:", err)
			}
			database.Database.Db.Create(&b)
		}(ws)
	}
	// var c *websocket.Conn

	// for {
	// 		if err := c.WriteMessage(websocket.TextMessage, []byte(b.Content)); err != nil {
	// 			fmt.Println("write err:", err)
	// 		}
	// 		database.Database.Db.Create(&b)
	// }
}