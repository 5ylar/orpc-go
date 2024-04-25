package orpcgo

import (
	"context"
	"errors"
	"reflect"
	"strings"
)

type Context struct {
	Ctx        context.Context
	MethodName string
	Headers    map[string][]string
}

type Middleware func(c Context) error

type Handler struct {
	Middlewares []Middleware
	H           interface{}
}

type ORPC struct {
	adapter           Adapter
	handlers          map[string]Handler
	globalMiddlewares []Middleware
}

func NewORPC(adapter Adapter) *ORPC {
	return &ORPC{
		adapter:  adapter,
		handlers: make(map[string]Handler),
	}
}

func (o *ORPC) SetGlobalMiddlewares(middlewares []Middleware) {
	o.globalMiddlewares = middlewares
}

func (o *ORPC) Handle(name string, h interface{}, middlewares ...Middleware) *ORPC {
	htype := reflect.TypeOf(h)

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

	o.handlers[name] = Handler{middlewares, h}
	return o
}

func (o *ORPC) Start(ctx context.Context) error {
	return o.adapter.Start(func(c AdapterCtx) error {

		if len(strings.TrimSpace(c.MethodName)) == 0 {
			return errors.New("invalid method name")
		}

		h, ok := o.handlers[c.MethodName]

		if !ok {
			return errors.New("not found method name")
		}

		hc := Context{
			Ctx:        c.Ctx,
			MethodName: c.MethodName,
			Headers:    c.Headers,
		}

		for _, m := range o.globalMiddlewares {
			err := m(hc)

			if err != nil {
				return err
			}
		}

		for _, m := range h.Middlewares {
			err := m(hc)

			if err != nil {
				return err
			}
		}

		preqtype := reflect.TypeOf(h.H).In(1)
		preqv := reflect.New(preqtype.Elem())
		preq := preqv.Interface() // pointer

		if err := c.Bind(preq); err != nil {
			return err
		}

		prepl := reflect.ValueOf(h.H).Call(
			[]reflect.Value{
				reflect.ValueOf(hc),
				preqv,
			},
		)

		// an error
		if !prepl[1].IsNil() {
			return prepl[1].Interface().(error)
		}

		return c.SendJSON(200, prepl[0].Interface())
	})
}
