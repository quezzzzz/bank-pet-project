package service

import (
	"bank"
	"bank/pkg/repository"
)

type CreditService struct {
	repo repository.Credit
}

func NewCreditService(repo repository.Credit) *CreditService {
	return &CreditService{repo: repo}
}

func (s *CreditService) TakeCredit(id, value, variable int) (int, bank.Credit, error) {
	var credit bank.Credit
	switch variable {
	case 1:
		credit = bank.Credit{
			CustomerId: id,
			Value:      value,
			Percentage: 10,
			LoanPeriod: 5,
		}
		credit.CurrentDebt = credit.Value + (credit.Value / credit.Percentage)

	case 2:
		credit = bank.Credit{
			CustomerId: id,
			Value:      value,
			Percentage: 15,
			LoanPeriod: 60,
		}
		credit.CurrentDebt = credit.Value + (credit.Value / credit.Percentage)
	case 3:
		credit = bank.Credit{
			CustomerId: id,
			Value:      value,
			Percentage: 20,
			LoanPeriod: 90,
		}
		credit.CurrentDebt = credit.Value + (credit.Value / credit.Percentage)
	}

	var err error
	var balance int

	credit.Id, balance, err = s.repo.TakeCredit(credit)
	if err != nil {
		return 0, credit, err
	}

	return balance, credit, nil
}

func (s *CreditService) CloseCredit(creditId, id, value int) (int, error) {
	return s.repo.CloseCredit(creditId, id, value)
}
