package domain

type User struct {
	ID        uint
	Username  string
	Password  string
	CreatedAt int64
}

type UserRepository interface {
	Create(user *User) error
	Update(id uint, user *User) error
	Delete(id uint) error
	GetByID(id uint) (*User, error)
	GetByUsername(username string) (*User, error)
}
