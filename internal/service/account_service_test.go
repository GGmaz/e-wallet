package service

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/mocks"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http/httptest"
	"testing"
)

func TestVerifyAccount(t *testing.T) {
	// Mock the AccountRepo
	mockAccountRepo := &mocks.Repo[model.Account]{}

	// Create a new AccountServiceImpl with the mockAccountRepo
	s := &AccountServiceImpl{
		AccountRepo: mockAccountRepo,
	}

	// Create a new gin.Context
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	db := &gorm.DB{}
	c.Set("transaction", db)

	// Create some sample accounts
	accounts := []*model.Account{
		{ID: 1, AccNumber: "1234567890"},
		{ID: 2, AccNumber: "0987654321"},
	}

	// Fill the mockAccountRepo with sample data
	mockAccountRepo.FillRepoWithAccountData(db, accounts)

	// Test with a valid account number
	t.Run("ValidAccount", func(t *testing.T) {
		err := s.VerifyAccount(c, "1234567890")
		if err != nil {
			t.Errorf("Expected error to be nil for valid account, got %v", err)
		}
	})

	// Test with an invalid account number
	t.Run("InvalidAccount", func(t *testing.T) {
		err := s.VerifyAccount(c, "111122223333") // Example invalid account number
		if err == nil {
			t.Error("Expected error for invalid account, but got nil")
		}

		// You can further customize this test case based on how you handle invalid account numbers
	})
}
