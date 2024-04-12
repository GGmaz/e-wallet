package wire

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/internal/repo"
	"github.com/GGmaz/wallet-arringo/internal/service"
)

type Wires struct {
	TransactionService service.TransactionServiceImpl
	UserService        service.UserServiceImpl
	AccountService     service.AccountServiceImpl
}

func Init() *Wires {
	w := Wires{
		TransactionService: service.TransactionServiceImpl{
			TransactionRepo: repo.Repo[model.Transaction]{},
			UserRepo:        repo.Repo[model.User]{},
			AccountRepo:     repo.Repo[model.Account]{},
		},
		UserService: service.UserServiceImpl{
			UserRepo: repo.Repo[model.User]{},
		},
		AccountService: service.AccountServiceImpl{
			AccountRepo: repo.Repo[model.Account]{},
			UserRepo:    repo.Repo[model.User]{},
		},
	}
	return &w
}

var Svc = Init()
