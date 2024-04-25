package main

import (
	"context"
	orpcgo "github.com/5ylar/orpc-go"
)

type TestRequest struct {
	Text string `json:"text"`
}

type TestReply struct {
	CharNum int `json:"char_num"`
}

func main() {
	o := orpcgo.NewORPC(
		orpcgo.NewDefaultAdapter(),
	)

	o.Handle("oms.test", func(c orpcgo.Context, i *TestRequest) (*TestReply, error) {
		return &TestReply{len(i.Text)}, nil
	})

	o.Handle("oms.test2", func(c orpcgo.Context, i *TestRequest) (*TestReply, error) {
		return &TestReply{1000}, nil
	})

	_ = o.Start(context.Background())
}
