// mocks package
package mocks

import (
	"errors"
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Repo is a mock implementation of the repo.Repo[model.Account] interface
type Repo[T any] struct {
	mock.Mock
}

// FillRepoWithAccountData fills the mock Repo with some data for testing purposes
func (m *Repo[T]) FillRepoWithAccountData(db *gorm.DB, accounts []*model.Account) {
	// Mock the GetByField method to return the provided accounts
	for _, acc := range accounts {
		accCopy := acc // Create a copy of acc to use inside the closure
		// Return a closure that updates the provided account pointer with the account data
		m.On("GetByField", db, mock.AnythingOfType("*model.Account"), "acc_number", acc.AccNumber).Run(func(args mock.Arguments) {
			accPtr := args.Get(1).(*model.Account) // Get the provided account pointer from the arguments
			*accPtr = *accCopy                     // Update the provided account pointer with the account data
		}).Return(db).Once()

		m.On("Update", db, mock.AnythingOfType("*model.Account"), acc.ID).Return(db).Once()
	}

	// Handle the case for the invalid account number
	invalidAccNumber := "111122223333"
	dbCopy := db // Create a copy of acc to use inside the closure
	m.On("GetByField", db, mock.AnythingOfType("*model.Account"), "acc_number", invalidAccNumber).Return(db).Once()
	m.On("Update", db, mock.AnythingOfType("*model.Account"), int64(0)).Run(func(args mock.Arguments) {
		dbPtr := args.Get(0).(*gorm.DB)
		dbPtr.Error = errors.New("error updating account")
		*dbPtr = *dbCopy
	}).Return(db).Once()
}

func (m *Repo[T]) GetByField(db *gorm.DB, t *T, field string, fieldVal string, preload ...string) *gorm.DB {
	var args mock.Arguments

	if len(preload) == 0 {
		args = m.Called(db, t, field, fieldVal)
	} else {
		args = m.Called(db, t, field, fieldVal, preload)
	}

	return args.Get(0).(*gorm.DB)
}

func (m *Repo[T]) Update(db *gorm.DB, t *T, id int64) *gorm.DB {
	args := m.Called(db, t, id)
	return args.Get(0).(*gorm.DB)
}

func (m *Repo[T]) Create(db *gorm.DB, t *T) *gorm.DB {
	//TODO implement me
	panic("implement me")
}

func (m *Repo[T]) GetById(db *gorm.DB, t *T, id int64, preload ...string) *gorm.DB {
	//TODO implement me
	panic("implement me")
}
