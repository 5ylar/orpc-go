package orpcgo

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type Context struct {
	Ctx context.Context
}

type ORPC struct {
	handlers map[string]interface{}
}

func NewORPC() *ORPC {
	return &ORPC{
		handlers: make(map[string]interface{}),
	}
}

func (o *ORPC) Register(name string, h interface{}) *ORPC {
    htype :=reflect.TypeOf(h)

	ctxtype := htype.In(0)
	preqtype := htype.In(1)
	prepltype := htype.Out(0)
	errtype := htype.Out(1)

	if ctxtype != reflect.TypeOf(Context{}) {
		panic("cannot register")
	}

	if preqtype.Kind() != reflect.Ptr {
		panic("cannot register")
	}

	if preqtype.Elem().Kind() != reflect.Struct {
		panic("cannot register")
	}

	if prepltype.Kind() != reflect.Ptr {
		panic("cannot register")
	}

	if prepltype.Elem().Kind() != reflect.Struct {
		panic("cannot register")
	}

	if errtype != reflect.TypeOf((*error)(nil)).Elem() {
		panic("cannot register")
	}

	o.handlers[name] = h
	return o
}

func (o *ORPC) Start(ctx context.Context) error {
	app := fiber.New()

	app.Post("/rpc/:method_name", func(c fiber.Ctx) error {
		methodName := c.Params("method_name", "")

		if len(strings.TrimSpace(methodName)) == 0 {
			return errors.New("invalid method name")
		}

		h, ok := o.handlers[methodName]

		if !ok {
			return errors.New("not found method name")
		}

		preqtype := reflect.TypeOf(h).In(1)
		preqv := reflect.New(preqtype.Elem())
		preq := preqv.Interface()

		if err := c.Bind().Body(preq); err != nil {
			return err
		}

		prepl := reflect.ValueOf(h).Call(
			[]reflect.Value{
				reflect.ValueOf(Context{
					Ctx: c.UserContext(),
				}),
				preqv,
			},
		)

		// an error
		if !prepl[1].IsNil() {
			return prepl[1].Interface().(error)
		}

		return c.JSON(prepl[0].Interface())
	})

	return app.Listen(":8080")
}
