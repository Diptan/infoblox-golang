package addressbook

import (
	"fmt"
	"infoblox-golang/internal/platform/storage"

	"gorm.io/gorm"
)

// Repository is a type that manipulates user data
type Repository struct {
	db *gorm.DB
}

// User is a type that represents user model
type User struct {
	gorm.Model
	ID       storage.ID `json:"id"`
	Username string     `json:"username"`
	Address  string     `json:"address"`
	Phone    string     `json:"phone"`
}

// NewRepository creates anew instance of Repository
func NewRepository(db *gorm.DB) Repository {
	return Repository{
		db: db,
	}
}

// Create creates a new user entity
func (r Repository) Create(u User) (User, error) {
	u.ID = storage.NewID()
	if res := r.db.Create(&u); res.Error != nil {
		return u, fmt.Errorf("save user: %w", res.Error)
	}

	return u, nil
}

// Remove removes an existing user entity by id
func (r Repository) Delete(id string) error {
	if res := r.db.Where("id = ?", id).Delete(&User{}); res.Error != nil {
		return fmt.Errorf("delete user: %w", res.Error)
	}

	return nil
}

// Update updates an existing user entity
func (r Repository) Update(id string, u User) (User, error) {
	u.ID = storage.ID(id)
	if r.db.Model(&u).Where("id = ?", id).Updates(&u).RowsAffected == 0 {
		r.db.Create(&u) // create new record from u
	}

	return u, nil
}

// Get returns a user entity for the specified ID
func (r Repository) GetById(id string) (User, error) {
	var user User
	res := r.db.First(&user, id)
	if res.Error != nil {
		return User{}, fmt.Errorf("get user: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return User{}, storage.ErrNotFound
	}

	return user, nil
}

// GetAll returns all user entities
func (r Repository) GetAll() ([]User, error) {
	var users []User
	res := r.db.Find(&users)
	if res.Error != nil {
		return nil, fmt.Errorf("get all users: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return []User{}, storage.ErrNotFound
	}

	return users, nil
}
