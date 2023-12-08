package usecase

import (
	"bank-api/internal/domain"
	"context"
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
	Password []byte
}

func (useCase CreateUserUseCase) CreateUserHandler(ctx context.Context, command *CreateUserCommand) (*domain.User, error) {
	password, err := GenerateFromPassword(command.Password)
	if err != nil {
		return nil, err
	}
	user := domain.NewUser(command.Login, password)
	err = useCase.userRepository.Save(ctx, user)
	return user, err
}
