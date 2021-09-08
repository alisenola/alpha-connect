package commands

import (
	"gitlab.com/alphaticks/alpha-connect/account"
	"gitlab.com/alphaticks/alpha-connect/models"
)

type GetAccountRequest struct {
	RequestID uint64
	Account   *models.Account
}

type GetAccountResponse struct {
	Account *account.Account
	Err     error
}
