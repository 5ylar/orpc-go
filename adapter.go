package orpcgo

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

type AdapterCtx struct {
	Ctx        context.Context
	MethodName string
	Headers    map[string][]string
	Bind       func(dest interface{}) error
}

type Adapter interface {
	Start(h func(c AdapterCtx) (interface{}, error)) error
}

type DefaultAdapter struct {
	app *fiber.App
}

func NewDefaultAdapter() *DefaultAdapter {
	app := fiber.New()

	return &DefaultAdapter{
		app,
	}
}

func (a *DefaultAdapter) Start(h func(c AdapterCtx) (interface{}, error)) error {
	a.app.Post("/rpc/:method_name", func(c fiber.Ctx) error {
		methodName := c.Params("method_name", "")

		repl, err := h(AdapterCtx{
			Ctx:        c.Context(),
			MethodName: methodName,
			Headers:    c.GetReqHeaders(),
			Bind:       c.Bind().Body,
		})

		if err != nil {
			c.Status(500)
			return err
		}

		return c.Status(200).JSON(repl)
	})

	return a.app.Listen(":8080")
}
