package service

import "bank/pkg/repository"

type TransService struct {
	repo repository.Transaction
}

func NewTransService(repo repository.Transaction) *TransService {
	return &TransService{repo: repo}
}

func (s *TransService) DepositMoney(id, value int) (int, error) {
	return s.repo.DepositMoney(id, value)
}

func (s *TransService) WithdrawMoney(id, value int) (int, error) {
	return s.repo.WithdrawMoney(id, value)

}
