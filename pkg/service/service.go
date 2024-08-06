package service

import (
	"bank"
	"bank/pkg/repository"
)

type Authorization interface {
	CreateCustomer(customer bank.Customer) (int, error)
	GenerateToken(phone, password string) (string, error)
	ParseToken(token string) (int, error)
}
type Transaction interface {
	DepositMoney(id, value int) (int, error)
	WithdrawMoney(id, value int) (int, error)
}
type Credit interface {
	TakeCredit(id, value, variable int) (int, bank.Credit, error)
	CloseCredit(creditId, id, value int) (int, error)
}

type Service struct {
	Authorization
	Transaction
	Credit
}

func NewService(repository *repository.Repository) *Service {

	return &Service{
		Authorization: NewAuthService(repository.Authorization),
		Transaction:   NewTransService(repository.Transaction),
		Credit:        NewCreditService(repository.Credit),
	}

}
