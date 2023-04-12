package addressbook

import (
	"address-book/pkg/storage"
	"fmt"
)

type Store interface {
	Find(m storage.Matcher) ([]interface{}, error)
	Insert(r interface{}) error
	Update(r interface{}) error
	Remove(id storage.ID) error
}

// Repository is a type that manipulates user data
type Repository struct {
	db Store
}

// User is a type that represents user model
type User struct {
	ID       storage.ID `json:"id"`
	Username string     `json:"username"`
	Address  string     `json:"address"`
	Phone    string     `json:"phone"`
}

// NewRepository creates anew instance of Repository
func NewRepository(db Store) Repository {
	return Repository{
		db: db,
	}
}

// Create creates a new user entity
func (r Repository) Create(d User) (User, error) {
	d.ID = storage.NewID()

	if err := r.db.Insert(d); err != nil {
		return User{}, fmt.Errorf("save user: %w", err)
	}

	return d, nil
}

// Remove removes an existing user entity by id
func (r Repository) Delete(id string) error {
	if err := r.db.Remove(storage.ID(id)); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

// Update updates an existing user entity
func (r Repository) Update(id string, u User) (User, error) {
	u.ID = storage.ID(id)
	err := r.db.Update(u)
	if err != nil {
		return User{}, fmt.Errorf("update user: %w", err)
	}

	return u, nil
}

// Get returns a user entity for the specified ID
func (r Repository) GetById(id string) (User, error) {
	res, err := r.db.Find(func(value interface{}) bool {
		if c, ok := value.(User); ok {
			return string(c.ID) == id
		}

		return false
	})

	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}

	if len(res) == 0 {
		return User{}, storage.ErrNotFound
	}

	return res[0].(User), nil
}

// GetAll returns all user entities
func (r Repository) GetAll() ([]User, error) {
	res, err := r.db.Find(func(value interface{}) bool {
		_, ok := value.(User)
		return ok
	})

	if err != nil {
		return []User{}, fmt.Errorf("get all users: %w", err)
	}

	if len(res) == 0 {
		return []User{}, storage.ErrNotFound
	}

	users := make([]User, len(res))

	for i, withID := range res {
		users[i] = withID.(User)
	}

	return users, nil
}
