package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TransPostgres struct {
	db *sqlx.DB
}

func NewTransPostgres(db *sqlx.DB) *TransPostgres {
	return &TransPostgres{db: db}
}

func (s *TransPostgres) DepositMoney(id, value int) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	var balance int
	updateCustomerBalanceQuery := fmt.Sprintf("UPDATE %s SET balance = balance + $1 WHERE id = $2 RETURNING balance", customersTable)
	row := tx.QueryRow(updateCustomerBalanceQuery, value, id)
	if err := row.Scan(&balance); err != nil {
		tx.Rollback()
		return 0, err
	}
	updateBankBalanceQuery := fmt.Sprintf("UPDATE %s SET storage = storage + $1", storage)
	_, err = tx.Exec(updateBankBalanceQuery, value)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}
	return balance, tx.Commit()
}

func (s *TransPostgres) WithdrawMoney(id, value int) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	var balance int
	updateCustomerBalanceQuery := fmt.Sprintf("UPDATE %s SET balance = balance - $1 WHERE id = $2 RETURNING balance", customersTable)
	row := tx.QueryRow(updateCustomerBalanceQuery, value, id)
	if err := row.Scan(&balance); err != nil {
		tx.Rollback()
		return 0, err
	}
	updateBankBalanceQuery := fmt.Sprintf("UPDATE %s SET storage = storage - $1", storage)
	_, err = tx.Exec(updateBankBalanceQuery, value)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}
	return balance, tx.Commit()
}
