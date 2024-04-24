## Example usage
```go
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
	o := orpcgo.NewORPC()

	o.Register("oms.test", func(ctx orpcgo.Context, i *TestRequest) (*TestReply, error) {
		return &TestReply{len(i.Text)}, nil
	})

	o.Register("oms.test2", func(ctx orpcgo.Context, i *TestRequest) (*TestReply, error) {
		return &TestReply{1000}, nil
	})

	_ = o.Start(context.Background())
}
```
