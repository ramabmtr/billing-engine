package service

import (
	"errors"
	"testing"
	"time"

	"github.com/ramabmtr/billing-engine/internal/model"
	"github.com/ramabmtr/billing-engine/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockBorrowerRepo is a mock implementation of repository.BorrowerRepo
type MockBorrowerRepo struct {
	mock.Mock
}

func (m *MockBorrowerRepo) WithTx(tx *gorm.DB) repository.BorrowerRepo {
	args := m.Called(tx)
	return args.Get(0).(repository.BorrowerRepo)
}

func (m *MockBorrowerRepo) Create(b *model.Borrower) error {
	args := m.Called(b)
	// Set ID if it's empty to simulate the BeforeCreate hook
	if b.ID == "" {
		b.ID = uuid.Must(uuid.NewV7()).String()
	}
	return args.Error(0)
}

func (m *MockBorrowerRepo) List() ([]*model.BorrowerWithDelinquentStatus, error) {
	args := m.Called()
	return args.Get(0).([]*model.BorrowerWithDelinquentStatus), args.Error(1)
}

func TestBorrowerService_Create(t *testing.T) {
	tests := []struct {
		name          string
		borrowerName  string
		mockSetup     func(mockRepo *MockBorrowerRepo)
		expectedError bool
	}{
		{
			name:         "Success",
			borrowerName: "John Doe",
			mockSetup: func(mockRepo *MockBorrowerRepo) {
				mockRepo.On("Create", mock.MatchedBy(func(b *model.Borrower) bool {
					return b.Name == "John Doe"
				})).Return(nil)
			},
			expectedError: false,
		},
		{
			name:         "Repository Error",
			borrowerName: "Jane Doe",
			mockSetup: func(mockRepo *MockBorrowerRepo) {
				mockRepo.On("Create", mock.MatchedBy(func(b *model.Borrower) bool {
					return b.Name == "Jane Doe"
				})).Return(errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBorrowerRepo)
			tt.mockSetup(mockRepo)

			service := NewBorrowerService(mockRepo)
			borrower, err := service.Create(tt.borrowerName)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, borrower)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, borrower)
				assert.Equal(t, tt.borrowerName, borrower.Name)
				assert.NotEmpty(t, borrower.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBorrowerService_List(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(mockRepo *MockBorrowerRepo)
		expectedError bool
		expectedCount int
	}{
		{
			name: "Success with borrowers",
			mockSetup: func(mockRepo *MockBorrowerRepo) {
				borrowers := []*model.BorrowerWithDelinquentStatus{
					{
						Borrower: model.Borrower{
							ID:        uuid.Must(uuid.NewV7()).String(),
							Name:      "John Doe",
							CreatedAt: time.Now(),
						},
						IsDelinquent: false,
					},
					{
						Borrower: model.Borrower{
							ID:        uuid.Must(uuid.NewV7()).String(),
							Name:      "Jane Doe",
							CreatedAt: time.Now(),
						},
						IsDelinquent: true,
					},
				}
				mockRepo.On("List").Return(borrowers, nil)
			},
			expectedError: false,
			expectedCount: 2,
		},
		{
			name: "Success with empty list",
			mockSetup: func(mockRepo *MockBorrowerRepo) {
				borrowers := []*model.BorrowerWithDelinquentStatus{}
				mockRepo.On("List").Return(borrowers, nil)
			},
			expectedError: false,
			expectedCount: 0,
		},
		{
			name: "Repository Error",
			mockSetup: func(mockRepo *MockBorrowerRepo) {
				mockRepo.On("List").Return([]*model.BorrowerWithDelinquentStatus{}, errors.New("database error"))
			},
			expectedError: true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBorrowerRepo)
			tt.mockSetup(mockRepo)

			service := NewBorrowerService(mockRepo)
			borrowers, err := service.List()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, borrowers, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
