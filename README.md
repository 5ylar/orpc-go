# orpc-go
A simple JSON RPC Server with ability to change HTTP Server behind whatever you like!

## Quick start
```go
package main

import (
	"context"
	orpcgo "github.com/5ylar/orpc-go"
)

type WalletDepositRequest struct {
	Amount float64 `json:"amount"`
}

type WalletDepositReply struct {
	Balance float64 `json:"balance"`
}

type GameBetRequest struct {
	Number int16   `json:"number"`
	Amount float64 `json:"amount"`
}

type GameBetReply struct {
	IsWin   bool    `json:"is_win"`
	Balance float64 `json:"balance"`
}

func main() {
	o := orpcgo.NewORPC(
		orpcgo.NewDefaultAdapter(),
	)

	o.Handle("wallet.deposit", func(c orpcgo.Context, i *WalletDepositRequest) (*WalletDepositReply, error) {
		return &WalletDepositReply{
			i.Amount,
		}, nil
	})

	o.Handle("game.bet", func(c orpcgo.Context, i *GameBetRequest) (*GameBetReply, error) {
		return &GameBetReply{
			IsWin:   false,
			Balance: 0,
		}, nil
	})

	_ = o.Start(context.Background())
}
```

```sh
go run example/quick-start/main.go
```

```sh
curl -X POST http://localhost:8080/rpc/game.bet -d '{"number": 9, "amount": 100}' | jq
```
