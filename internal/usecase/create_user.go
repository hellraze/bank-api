package usecase

import (
	"bank-api/internal/domain"
)

type CreateUserUseCase struct {
	userRepository domain.UserRepository
}

func NewCreateUserUseCase(userRepository domain.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
	}
}

type CreateUserCommand struct {
	Login    string
	Password string
}

func (useCase CreateUserUseCase) CreateUserHandler(command *CreateUserCommand) (*domain.User, error) {
	user := domain.NewUser(command.Login, command.Password)
	err := useCase.userRepository.Save(user)
	return user, err
}
