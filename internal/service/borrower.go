package service

import (
	"github.com/ramabmtr/billing-engine/internal/model"
	"github.com/ramabmtr/billing-engine/internal/repository"
)

type BorrowerService struct {
	borrowerRepo repository.BorrowerRepo
}

func NewBorrowerService(borrowerRepo repository.BorrowerRepo) *BorrowerService {
	return &BorrowerService{
		borrowerRepo: borrowerRepo,
	}
}

func (s *BorrowerService) Create(name string) (*model.Borrower, error) {
	b := &model.Borrower{
		Name: name,
	}
	err := s.borrowerRepo.Create(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *BorrowerService) List() ([]*model.BorrowerWithDelinquentStatus, error) {
	l, err := s.borrowerRepo.List()
	return l, err
}
