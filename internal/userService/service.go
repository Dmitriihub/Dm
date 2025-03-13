package userService

type UserService struct {
	repo userRepository
}

func NewService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user User) (User, error) {
	return s.repo.CreateUser(user)
}

func (s *UserService) GetAllUsers() ([]User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) DeleteUserByID(id uint) error {
	return s.repo.DeleteUserByID(id)
}

// Пример метода в UserService
func (s *UserService) UpdateUserByID(id int, user User) (User, error) {
	// Логика для обновления пользователя по id
	// Пример: обновить пользователя в репозитории
	updatedUser, err := s.repo.UpdateUserByID(uint(id), user)
	if err != nil {
		return User{}, err
	}
	return updatedUser, nil
}
