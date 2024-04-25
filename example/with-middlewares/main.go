package main

import (
	"context"
	"errors"

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

var (
	ErrUnauthorized = errors.New("Unauthorized")
)

func AuthMiddleware(c orpcgo.Context) error {
	v, ok := c.Headers["Authorization"]

	if len(v) == 0 || !ok {
		return ErrUnauthorized
	}

	if v[0] != "Dummy abc" {
		return ErrUnauthorized
	}

	return nil
}

func IndividualMiddleware(c orpcgo.Context) error {
	return nil
}

func main() {
	o := orpcgo.NewORPC(
		orpcgo.NewDefaultAdapter(),
	)

	o.SetGlobalMiddlewares([]orpcgo.Middleware{
		AuthMiddleware,
	})

	o.Handle(
		"wallet.deposit",
		func(c orpcgo.Context, i *WalletDepositRequest) (*WalletDepositReply, error) {
			return &WalletDepositReply{
				i.Amount,
			}, nil
		},
	)

	o.Handle(
		"game.bet",
		func(c orpcgo.Context, i *GameBetRequest) (*GameBetReply, error) {
			return &GameBetReply{
				IsWin:   false,
				Balance: 0,
			}, nil
		},
		IndividualMiddleware,
	)

	_ = o.Start(context.Background())
}
