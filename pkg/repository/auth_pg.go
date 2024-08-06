package repository

import (
	"bank"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (s *AuthPostgres) CreateCustomer(customer bank.Customer) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, surname, age, balance,phone, password_hash ) values($1, $2, $3, $4, $5, $6) RETURNING id", customersTable)
	row := s.db.QueryRow(query, customer.Name, customer.Surname, customer.Age, customer.Balance, customer.Phone, customer.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AuthPostgres) GetCustomer(phone, password string) (bank.Customer, error) {
	var customer bank.Customer
	query := fmt.Sprintf("SELECT id FROM %s WHERE phone = $1 AND password_hash = $2", customersTable)
	err := s.db.Get(&customer, query, phone, password)
	if err != nil {
		logrus.Error("Wrong phone or password")

	}
	return customer, err
}
