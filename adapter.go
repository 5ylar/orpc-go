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
	Handle(methodName string, h func(c AdapterCtx) (interface{}, error))
	Start() error
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

func (a *DefaultAdapter) Handle(methodName string, h func(c AdapterCtx) (interface{}, error)) {
	a.app.Post("/rpc/"+methodName, func(c fiber.Ctx) error {
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
}

func (a *DefaultAdapter) Start() error {
	return a.app.Listen(":8080")
}
