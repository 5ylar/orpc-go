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
	SendJSON   func(status int, resp interface{}) error
}

type Adapter interface {
	Start(h func(c AdapterCtx) error) error
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

func (a *DefaultAdapter) Start(h func(c AdapterCtx) error) error {
	a.app.Post("/rpc/:method_name", func(c fiber.Ctx) error {
		methodName := c.Params("method_name", "")

		err := h(AdapterCtx{
			Ctx:        c.Context(),
			MethodName: methodName,
			Headers:    c.GetReqHeaders(),
			Bind:       c.Bind().Body,
			SendJSON: func(status int, resp interface{}) error {
				return c.Status(status).JSON(resp)
			},
		})

		if err != nil {
			return err
		}

		return nil
	})

	return a.app.Listen(":8080")
}
