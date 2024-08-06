package repository

import (
	"bank"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type CreditPostgres struct {
	db *sqlx.DB
}

func NewCreditPostgres(db *sqlx.DB) *CreditPostgres {
	return &CreditPostgres{db: db}
}

func (s *CreditPostgres) TakeCredit(credit bank.Credit) (int, int, error) {
	var availableMoney int
	checkStorage := fmt.Sprintf("SELECT storage FROM %s ", storage)
	err := s.db.Get(&availableMoney, checkStorage)
	if (availableMoney / 2) < credit.Value {
		return 0, 0, errors.New("credit denied")
	}
	tx, err := s.db.Begin()
	var balance int
	var id int
	newCreditQuery := fmt.Sprintf("INSERT INTO %s (customer_id,value,percentage,loan_period,current_debt) VALUES($1, $2, $3, $4, $5) RETURNING id", creditTable)
	newCreditRow := tx.QueryRow(newCreditQuery, credit.CustomerId, credit.Value, credit.Percentage, credit.LoanPeriod, credit.CurrentDebt)
	if err := newCreditRow.Scan(&id); err != nil {
		tx.Rollback()
		return 0, 0, err
	}
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}
	updateCustomerBalanceQuery := fmt.Sprintf("UPDATE %s SET balance = balance + $1 WHERE id = $2 RETURNING balance", customersTable)
	updateCustomerBalanceRow := tx.QueryRow(updateCustomerBalanceQuery, credit.Value, credit.CustomerId)
	if err := updateCustomerBalanceRow.Scan(&balance); err != nil {
		tx.Rollback()
		return 0, 0, err
	}
	updateBankBalanceQuery := fmt.Sprintf("UPDATE %s SET storage = storage - $1", storage)
	_, err = tx.Exec(updateBankBalanceQuery, credit.Value)
	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	return id, balance, tx.Commit()
}

func (s *CreditPostgres) CloseCredit(creditId, id, value int) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	var currentDebt int
	updateCredits := fmt.Sprintf("UPDATE %s SET current_debt = current_debt - $1 WHERE id = $2 AND customer_id = $3 RETURNING current_debt", creditTable)
	row := tx.QueryRow(updateCredits, value, creditId, id)
	if err := row.Scan(&currentDebt); err != nil {
		logrus.Println("2")
		tx.Rollback()
		return 0, err
	}

	updateCustomerBalance := fmt.Sprintf("UPDATE %s SET balance = balance - $1 WHERE ID = $2", customersTable)
	_, err = tx.Exec(updateCustomerBalance, value, id)
	if err != nil {
		logrus.Println("3")
		tx.Rollback()
		return 0, err
	}

	updateBankBalance := fmt.Sprintf("UPDATE %s SET storage = storage + $1", storage)

	_, err = tx.Exec(updateBankBalance, value)
	if err != nil {
		logrus.Println("4")
		tx.Rollback()
		return 0, nil
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	if currentDebt == 0 {
		deleteCredits := fmt.Sprintf("DELETE FROM %s WHERE  id = $1 AND customer_id = $2", creditTable)
		_, err := s.db.Exec(deleteCredits, creditId, id)
		if err != nil {
			logrus.Println("5")
			return 0, err
		}
	}
	return currentDebt, nil

}
