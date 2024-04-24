package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

type TestRequest struct {
	Text string `json:"text"`
}

type TestReply struct {
	CharNum int `json:"char_num"`
}

func main() {
	app := fiber.New()

	app.Post("/", func(c fiber.Ctx) error {
		var preq TestRequest

		if err := c.Bind().Body(&preq); err != nil {
			return err
		}

		return c.JSON(&TestReply{CharNum: len(preq.Text)})
	})

	log.Fatal(app.Listen(":8080"))
}
