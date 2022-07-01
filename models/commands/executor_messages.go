package commands

import (
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/config"
)

type GetAccountRequest struct {
	RequestID uint64
	Account   *config.Account
}

type GetAccountResponse struct {
	Account *account.Account
	Err     error
}
