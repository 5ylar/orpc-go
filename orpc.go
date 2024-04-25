package orpcgo

import (
	"context"
	"reflect"
)

type Context struct {
	Ctx        context.Context
	MethodName string
	Headers    map[string][]string
	SetStatus  func(status int)
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

func (o *ORPC) Run(ctx context.Context) error {
	for methodName, h := range o.handlers {

		preqtype := reflect.TypeOf(h.H).In(1)
		hr := reflect.ValueOf(h.H)

		o.adapter.Handle(methodName, func(c AdapterCtx) (interface{}, error) {
			hc := Context{
				Ctx:        c.Ctx,
				MethodName: methodName,
				Headers:    c.Headers,
				SetStatus:  c.SetStatus,
			}

			for _, m := range o.globalMiddlewares {
				err := m(hc)

				if err != nil {
					return nil, err
				}
			}

			for _, m := range h.Middlewares {
				err := m(hc)

				if err != nil {
					return nil, err
				}
			}

			preqv := reflect.New(preqtype.Elem())
			preq := preqv.Interface() // pointer

			if err := c.Bind(preq); err != nil {
				return nil, err
			}

			prepl := hr.Call(
				[]reflect.Value{
					reflect.ValueOf(hc),
					preqv,
				},
			)

			// an error
			if !prepl[1].IsNil() {
				return nil, prepl[1].Interface().(error)
			}

			return prepl[0].Interface(), nil
		})
	}

	return o.adapter.Start()
}
