package addressbook

// Service is a type that provides main addressbook operations
type Service struct {
	repository Repository
	//logger     logger.Logger
}

// NewService returns a new user service instance
func NewService(r Repository) *Service {
	return &Service{
		repository: r,
	}
}

func (s *Service) AddUser(user User) (User, error) {
	return s.repository.Create(user)
}

func (s *Service) GetUserById(id string) (User, error) {
	return s.repository.GetById(id)
}

func (s *Service) GetAllUsers() ([]User, error) {
	return s.repository.GetAll()
}

func (s *Service) UpdateUser(id string, u User) (User, error) {
	return s.repository.Update(id, u)
}

func (s *Service) DeleteUser(id string) error {
	return s.repository.Delete(id)
}
