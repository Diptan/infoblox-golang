package addressbook

import (
	"address-book/pkg/storage"
	"fmt"
	"time"
)

type Store interface {
	Find(m storage.Matcher) ([]interface{}, error)
	Insert(r interface{}) error
}

// Repository is a type that manipulates addressbook data
type Repository struct {
	db Store
}

// AddressBook is a type that represents addressbook model
type AddressBook struct {
	ID       storage.ID `json:"id"`
	Username string     `json:"username"`
	Address  string     `json:"address"`
	Phone    time.Time  `json:"phone"`
}

// NewRepository creates anew instance of Repository
func NewRepository(db Store) Repository {
	return Repository{
		db: db,
	}
}

// Create creates a new addressbook entity
func (r Repository) Create(d AddressBook) (AddressBook, error) {
	d.ID = storage.NewID()

	if err := r.db.Insert(d); err != nil {
		return AddressBook{}, fmt.Errorf("save command: %w", err)
	}

	return d, nil
}

// Get returns a addressbook entity for the specified ID
func (r Repository) Get(id string) (AddressBook, error) {
	res, err := r.db.Find(func(value interface{}) bool {
		if c, ok := value.(AddressBook); ok {
			return string(c.ID) == id
		}

		return false
	})

	if err != nil {
		return AddressBook{}, fmt.Errorf("get command: %w", err)
	}

	if len(res) == 0 {
		return AddressBook{}, storage.ErrNotFound
	}

	return res[0].(AddressBook), nil
}
