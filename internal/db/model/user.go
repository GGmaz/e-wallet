package model

type User struct {
	ID        int64  `json:"user_id" gorm:"primaryKey;index:idx_user_id;type=varchar(50)"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Address   string `json:"address" gorm:"not null"`
	//Balance   float64   `json:"balance" gorm:"not null"`
	Email    string    `json:"email" gorm:"unique; not null"`
	Accounts []Account `json:"accounts" gorm:"foreignKey:UserId"`
}

func (tx User) GetName() string {
	return "User"
}
