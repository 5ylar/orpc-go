package main

import (
	"context"
	"errors"
	"log"
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
	ctxtype := reflect.TypeOf(h).In(0)
	preqtype := reflect.TypeOf(h).In(1)
	prepltype := reflect.TypeOf(h).Out(0)
	errtype := reflect.TypeOf(h).Out(1)

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

	app.Post("/rpc/:name", func(c fiber.Ctx) error {
		name := c.Params("name", "")

		if len(strings.TrimSpace(name)) == 0 {
			return errors.New("invalid method name")
		}

		h, ok := o.handlers[name]

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

	log.Fatal(app.Listen(":8080"))

	return nil
}

type TestRequest struct {
	Text string `json:"text"`
}

type TestReply struct {
	CharNum int `json:"char_num"`
}

func main() {
	o := NewORPC()

	o.Register("oms.test", func(ctx Context, i *TestRequest) (*TestReply, error) {
		return &TestReply{len(i.Text)}, nil
	})

	o.Register("oms.test2", func(ctx Context, i *TestRequest) (*TestReply, error) {
		return &TestReply{1000}, nil
	})

	_ = o.Start(context.Background())
}
