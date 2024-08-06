package repository

import (
	"bank"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateCustomer(customer bank.Customer) (int, error)
	GetCustomer(phone, password string) (bank.Customer, error)
}

type Transaction interface {
	DepositMoney(id, value int) (int, error)
	WithdrawMoney(id, value int) (int, error)
}

type Credit interface {
	TakeCredit(credit bank.Credit) (int, int, error)
	CloseCredit(creditId, id, value int) (int, error)
}

type Repository struct {
	Authorization
	Transaction
	Credit
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Transaction:   NewTransPostgres(db),
		Credit:        NewCreditPostgres(db),
	}
}
